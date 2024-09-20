package pokemon_stock_pair

import (
	"context"
	"encoding/json"
	"fmt"
	common_pb "pokestocks/proto/common"
	psp_pb "pokestocks/proto/pokemon_stock_pair"
	redis_keys "pokestocks/redis"
	"pokestocks/utils"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) queryDbForPokemonStockPairs(ctx context.Context, pspIds []string) ([]*common_pb.PokemonStockPair, error) {
	db := s.DB
	redisClient := s.RedisClient
	redisPipeline := redisClient.Pipeline()

	query := pspQueryString()

	queryArgs := []any{}
	positionalParams := []string{}

	for i, id := range pspIds {
		queryArgs = append(queryArgs, id)
		positionalParams = append(positionalParams, fmt.Sprintf("$%d", i+1))
	}

	orderByArgs := strings.Join(pspIds, ",")
	queryArgs = append(queryArgs, orderByArgs)

	positionalParamsString := strings.Join(positionalParams, ", ")
	query += fmt.Sprintf("WHERE psp.id IN (%s)", positionalParamsString)
	query += fmt.Sprintf("ORDER BY POSITION(psp.id::text IN $%d)", len(positionalParams)+1)

	rows, err := db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var psps []*common_pb.PokemonStockPair
	midnightTomorrow := midnightTomorrow()

	for rows.Next() {
		queriedData, err := pgx.RowToMap(rows)
		if err != nil {
			return nil, err
		}

		psp := convertDbRowToPokemonStockPair(queriedData)
		psps = append(psps, psp)

		jsonBytes, err := json.Marshal(psp)
		if err != nil {
			return nil, err
		}
		redisPipeline.JSONSet(ctx, redis_keys.DbCacheKey(fmt.Sprint(psp.Id)), "$", string(jsonBytes))
		redisPipeline.ExpireAt(ctx, redis_keys.DbCacheKey(fmt.Sprint(psp.Id)), midnightTomorrow)
	}

	if err = rows.Err(); err != nil {
		utils.LogWarningError("Unable to attempt to cache PSP JSON to Redis due to error reading db rows", err)
		return nil, err
	}

	_, err = redisPipeline.Exec(ctx)
	if err != nil {
		utils.LogWarningError("Error caching PSP JSON to Redis", err)
	}

	return psps, nil
}

func (s *Server) SearchPokemonStockPairs(ctx context.Context, in *psp_pb.SearchPokemonStockPairsRequest) (*psp_pb.SearchPokemonStockPairsResponse, error) {
	// startSearchPokemonStockPairs := time.Now()

	searchValue := in.SearchValue

	redisClient := s.RedisClient
	redisPipeline := redisClient.Pipeline()

	var pspIds []string

	cachedElasticIds, err := redisClient.ZRange(ctx, redis_keys.ElasticCacheKey(searchValue), 0, -1).Result()
	if err == nil && len(cachedElasticIds) != 0 {
		pspIds = cachedElasticIds
	} else {
		if err != nil {
			// if there is something wrong with Redis and it can't answer our request,
			// we can always just fallback to searching Elastic
			utils.LogWarningError("Error querying Redis key "+redis_keys.ElasticCacheKey(searchValue)+" for cached PSP ids. Falling back to Elastic", err)
		}
		searchResults, err := s.searchElasticIndex(searchValue)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error searching Elastic data: %v", err)
		}

		elasticPsps, err := convertPokemonStockPairElasticDocuments(searchResults)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error formatting Elastic data: %v", err)
		}
		if len(elasticPsps) == 0 {
			return &psp_pb.SearchPokemonStockPairsResponse{Data: nil}, nil
		}

		pspIds = extractPokemonStockPairIds(elasticPsps)

		redisPipeline := redisClient.Pipeline()
		sortedSet := []redis.Z{}
		for i, id := range pspIds {
			sortedSetMember := redis.Z{
				Score:  float64(i),
				Member: id,
			}
			sortedSet = append(sortedSet, sortedSetMember)
		}
		midnightTomorrow := midnightTomorrow()
		redisPipeline.ZAdd(ctx, redis_keys.ElasticCacheKey(searchValue), sortedSet...)
		redisPipeline.ExpireAt(ctx, redis_keys.ElasticCacheKey(searchValue), midnightTomorrow)

		_, err = redisPipeline.Exec(ctx)
		if err != nil {
			// don't return a gRPC response with an error
			// a response with data can still be generated even if we can't cache Elasticsearch results
			utils.LogWarningError("Error caching data to Redis for key "+redis_keys.ElasticCacheKey(searchValue)+". Skipping", err)
		}
	}

	cachedIds := []string{}
	cachedPspsMap := map[string]*common_pb.PokemonStockPair{}
	nonCachedIds := []string{}

	for _, id := range pspIds {
		redisPipeline.JSONGet(ctx, redis_keys.DbCacheKey(id)).Result()
	}

	results, err := redisPipeline.Exec(ctx)
	if err == nil {
		for _, result := range results {
			jsonString := result.(*redis.JSONCmd).Val()
			key := result.(*redis.JSONCmd).Args()[1]
			id := redis_keys.GetIdFromDbCacheKey(key.(string))
			if jsonString == "" {
				nonCachedIds = append(nonCachedIds, id)
			} else {
				cachedIds = append(cachedIds, id)

				var psp common_pb.PokemonStockPair
				json.Unmarshal([]byte(jsonString), &psp)
				cachedPspsMap[id] = &psp
			}
		}
	}

	var psps []*common_pb.PokemonStockPair

	if len(cachedIds) == 0 && len(nonCachedIds) == 0 {
		// the above loop on results never occurred so the id arrays are both empty
		utils.LogWarningError("Error querying Redis key "+redis_keys.ElasticCacheKey(searchValue)+" for cached JSON. Falling back to db", err)
		psps, err = s.queryDbForPokemonStockPairs(ctx, pspIds)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error querying data from db: %v", err)
		}
	} else if len(nonCachedIds) > 0 {
		// 1 or more cache misses require querying db
		nonCachedPspsArr, err := s.queryDbForPokemonStockPairs(ctx, nonCachedIds)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error querying data from db: %v", err)
		}

		nonCachedPspsMap := map[string]*common_pb.PokemonStockPair{}
		for _, nonCachedPsp := range nonCachedPspsArr {
			nonCachedPspsMap[fmt.Sprint(nonCachedPsp.Id)] = nonCachedPsp
		}

		for _, id := range pspIds {
			val, ok := cachedPspsMap[id]
			if ok {
				psps = append(psps, val)
			} else {
				psps = append(psps, nonCachedPspsMap[id])
			}
		}
	} else {
		// 0 cache misses
		for _, id := range cachedIds {
			psps = append(psps, cachedPspsMap[id])
		}
	}

	err = s.enrichWithStockPrices(ctx, psps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)
	}

	// fmt.Println("startSearchPokemonStockPairs took:", time.Since(startSearchPokemonStockPairs))
	// fmt.Println("======")

	return &psp_pb.SearchPokemonStockPairsResponse{Data: psps}, nil
}
