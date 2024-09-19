package pokemon_stock_pair

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	common_pb "pokestocks/proto/common"
	psp_pb "pokestocks/proto/pokemon_stock_pair"
	redis_keys "pokestocks/redis"
	"pokestocks/utils"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Clock struct {
	Timestamp time.Time `json:"timestamp"`
	IsOpen    bool      `json:"is_open"`
	NextOpen  time.Time `json:"next_open"`
	NextClose time.Time `json:"next_close"`
}

func (s *Server) isMarketOpen(ctx context.Context) (bool, error) {
	redisClient := s.RedisClient
	redisPipeline := redisClient.Pipeline()
	alpacaTradingClient := s.AlpacaTradingClient

	cachedMarketStatus, err := redisClient.Get(ctx, redis_keys.MarketStatusKey()).Result()
	if err == nil {
		return cachedMarketStatus == "open", nil
	} else {
		clock, err := getAlpacaClock(alpacaTradingClient)
		if err != nil {
			utils.LogWarning(fmt.Sprintf("Error hitting Alpaca clock API: %v", err))
			return false, err
		}

		marketIsOpen := clock.IsOpen
		if marketIsOpen {
			marketCloseTime := clock.NextClose
			redisPipeline.Set(ctx, redis_keys.MarketStatusKey(), "open", 0)
			redisPipeline.ExpireAt(ctx, redis_keys.MarketStatusKey(), marketCloseTime)
		} else {
			marketOpenTime := clock.NextOpen
			redisPipeline.Set(ctx, redis_keys.MarketStatusKey(), "close", 0)
			redisPipeline.ExpireAt(ctx, redis_keys.MarketStatusKey(), marketOpenTime)
		}

		_, err = redisPipeline.Exec(ctx)
		if err != nil {
			utils.LogWarning(fmt.Sprintf("Error caching market status: %v", err))
		}

		return marketIsOpen, nil
	}
}

func (s *Server) queryDbForPokemonStockPairs(ctx context.Context, pspIds []string) ([]*common_pb.PokemonStockPair, error) {
	// preparingDbQuery := time.Now()

	db := s.DB
	redisClient := s.RedisClient
	redisPipeline := redisClient.Pipeline()

	query := `
		SELECT 
			psp.id AS "pspId",
			pokemon.id AS "pokemonId",
			pokemon."name" AS "pokemonName",
			pokemon.pokedex_number AS "pokedexNumber",
			pokemon.created_at AS "pokemonCreatedAt",
			pokemon.updated_at AS "pokemonUpdatedAt",
			pokemon.sprite_url as "pokemonSpriteUrl",
			pokemon_types1.id AS "type1Id",
			pokemon_types1."type" AS "type1Name",
			pokemon_types1.sprite_url AS "type1SpriteUrl",
			pokemon_types2.id AS "type2Id",
			pokemon_types2."type" AS "type2Name",
			pokemon_types2.sprite_url AS "type2SpriteUrl",
			stocks.id AS "stockId",
			stocks.symbol AS "stockSymbol", 
			stocks."name" AS "stockName",
			stocks.created_at AS "stockCreatedAt",
			stocks.updated_at AS "stockUpdatedAt",
			stocks.active AS "stockActive",
			seasons.id AS "seasonId",
			seasons."name" AS "seasonName",
			seasons.active AS "seasonActive"
		FROM pokemon_stock_pairs AS psp
		INNER JOIN pokemon ON pokemon.id = psp.pokemon_id
		INNER JOIN pokemon_types AS pokemon_types1 ON pokemon_types1.id = pokemon.type_1_id
		LEFT JOIN pokemon_types AS pokemon_types2 ON pokemon_types2.id = pokemon.type_2_id
		INNER JOIN stocks ON stocks.id = psp.stock_id
		INNER JOIN seasons ON seasons.id = psp.season_id
	`

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

	// fmt.Println("preparingDbQuery took:", time.Since(preparingDbQuery))

	// queryingDb := time.Now()
	rows, err := db.Query(ctx, query, queryArgs...)
	// fmt.Println("queryingDb took:", time.Since(queryingDb))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// processingDbQuery := time.Now()

	var psps []*common_pb.PokemonStockPair

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
		redisPipeline.Expire(ctx, redis_keys.DbCacheKey(fmt.Sprint(psp.Id)), time.Second*10)
	}

	// savingDbKeys := time.Now()
	_, err = redisPipeline.Exec(ctx)
	// fmt.Println("savingDbKeys:", time.Since(savingDbKeys))

	if err != nil {
		log.Printf("Error inserting PSP JSON into Redis: %v", err)
	}

	// fmt.Println("processingDbQuery took:", time.Since(processingDbQuery))

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return psps, nil
}

func (s *Server) SearchPokemonStockPairs(ctx context.Context, in *psp_pb.SearchPokemonStockPairsRequest) (*psp_pb.SearchPokemonStockPairsResponse, error) {
	// startSearchPokemonStockPairs := time.Now()

	searchValue := in.SearchValue

	redisClient := s.RedisClient
	redisPipeline := redisClient.Pipeline()

	var pspIds []string

	// checkingForCachedElastic := time.Now()
	cachedElasticIds, err := redisClient.ZRange(ctx, redis_keys.ElasticCacheKey(searchValue), 0, -1).Result()
	// fmt.Println("checkingForCachedElastic:", time.Since(checkingForCachedElastic))

	if err == nil && len(cachedElasticIds) != 0 {
		fmt.Println("from cache instead of elastic")
		pspIds = cachedElasticIds
	} else {
		if err != nil {
			// if there is something wrong with Redis and it can't answer our request,
			// we can always just fallback to searching Elastic
			utils.LogWarning(fmt.Sprintf("Error querying Redis key %v for cached PSP ids. Falling back to Elastic: %v", redis_keys.ElasticCacheKey(searchValue), err))
		} else {
			log.Println("cache miss, going to elastic")
		}
		// searchingElastic := time.Now()
		searchResults, err := s.searchElasticIndex(searchValue)
		// fmt.Println("searchingElastic:", time.Since(searchingElastic))
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
		redisPipeline.ZAdd(ctx, redis_keys.ElasticCacheKey(searchValue), sortedSet...)
		redisPipeline.Expire(ctx, redis_keys.ElasticCacheKey(searchValue), time.Second*10)

		// savingElasticKeys := time.Now()
		_, err = redisPipeline.Exec(ctx)
		// fmt.Println("savingElasticKeys:", time.Since(savingElasticKeys))
		if err != nil {
			// don't return a gRPC response with an error
			// a response with data can still be generated even if we can't cache Elasticsearch results
			utils.LogWarning(fmt.Sprintf("Error saving data to Redis for key %v. Skipping: %v", redis_keys.ElasticCacheKey(searchValue), err))
		}
	}

	cachedIds := []string{}
	cachedPspsMap := map[string]*common_pb.PokemonStockPair{}
	nonCachedIds := []string{}

	for _, id := range pspIds {
		redisPipeline.JSONGet(ctx, redis_keys.DbCacheKey(id)).Result()
	}
	// gettingJson := time.Now()
	// if this returns an error, the following loop will not run,
	// resulting in the cachedIds and nonCachedIds to be empty slices
	// and the first condition of the if-else being true
	// therefore, there is no need to do a direct error check here
	results, _ := redisPipeline.Exec(ctx)
	// fmt.Println("gettingJson:", time.Since(gettingJson))

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

	var psps []*common_pb.PokemonStockPair

	if len(cachedIds) == 0 && len(nonCachedIds) == 0 {
		// something wrong with Redis and need to query db for all ids
		utils.LogWarning(fmt.Sprintf("Error querying Redis key %v for cached JSON. Falling back to Elastic: %v", redis_keys.ElasticCacheKey(searchValue), err))
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

	// gettingStockPrices := time.Now()

	err = s.enrichWithStockPrices(psps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)
	}

	// fmt.Println("gettingStockPrices took:", time.Since(gettingStockPrices))
	// fmt.Println("startSearchPokemonStockPairs took:", time.Since(startSearchPokemonStockPairs))
	// fmt.Println("======")

	return &psp_pb.SearchPokemonStockPairsResponse{Data: psps}, nil
}
