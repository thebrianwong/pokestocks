package pokemon_stock_pair

import (
	"context"
	"encoding/json"
	"fmt"
	common_pb "pokestocks/proto/common"
	psp_pb "pokestocks/proto/pokemon_stock_pair"
	redis_keys "pokestocks/redis"
	"pokestocks/utils"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetPokemonStockPair(ctx context.Context, in *psp_pb.GetPokemonStockPairRequest) (*psp_pb.GetPokemonStockPairResponse, error) {
	db := s.DB
	redisClient := s.RedisClient
	redisPipeline := redisClient.Pipeline()

	pspId := in.Id

	var psps []*common_pb.PokemonStockPair

	pspJson, err := redisClient.JSONGet(ctx, redis_keys.DbCacheKey(fmt.Sprint(pspId))).Result()
	if pspJson != "" && err == nil {
		var psp common_pb.PokemonStockPair
		json.Unmarshal([]byte(pspJson), &psp)
		psps = append(psps, &psp)
	} else {
		if err != nil {
			utils.LogWarningError("Error querying Redis key "+redis_keys.ElasticCacheKey(fmt.Sprint(pspId))+" for cached JSON. Falling back to db", err)
		}

		query := pspQueryString() + "WHERE psp.id = $1"
		rows, err := db.Query(ctx, query, pspId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error while querying for %v: %v", pspId, err)
		}
		defer rows.Close()

		for rows.Next() {
			queriedData, err := pgx.RowToMap(rows)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "error converting queried data to map: %v", err)
			}

			psp := convertDbRowToPokemonStockPair(queriedData)
			psps = append(psps, psp)

			jsonBytes, err := json.Marshal(psp)
			if err != nil {
				return nil, err
			}

			midnightTomorrow := midnightTomorrow()

			redisPipeline.JSONSet(ctx, redis_keys.DbCacheKey(fmt.Sprint(pspId)), "$", string(jsonBytes))
			redisPipeline.ExpireAt(ctx, redis_keys.DbCacheKey(fmt.Sprint(pspId)), midnightTomorrow)
		}

		if err = rows.Err(); err != nil {
			return nil, status.Errorf(codes.Internal, "error reading queried data: %v", err)
		}

		_, err = redisPipeline.Exec(ctx)
		if err != nil {
			utils.LogWarningError("Error caching PSP JSON to Redis", err)
		}
	}

	err = s.enrichWithStockPrices(ctx, psps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)
	}

	return &psp_pb.GetPokemonStockPairResponse{Data: psps}, nil
}
