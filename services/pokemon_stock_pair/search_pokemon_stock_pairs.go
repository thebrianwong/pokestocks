package pokemon_stock_pair

import (
	"context"
	"encoding/json"
	"fmt"
	common_pb "pokestocks/proto/common"
	psp_pb "pokestocks/proto/pokemon_stock_pair"
	redis_keys "pokestocks/redis"
	"pokestocks/utils"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SearchPokemonStockPairs(ctx context.Context, in *psp_pb.SearchPokemonStockPairsRequest) (*psp_pb.SearchPokemonStockPairsResponse, error) {
	// startSearchPokemonStockPairs := time.Now()

	searchValue := in.SearchValue

	redisClient := s.RedisClient
	redisPipeline := redisClient.Pipeline()

	pspIds, err := s.searchPokemonStockPairIds(ctx, searchValue)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error searching for PSP ids: %v", err)
	}
	if pspIds == nil {
		return &psp_pb.SearchPokemonStockPairsResponse{Data: nil}, nil
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
