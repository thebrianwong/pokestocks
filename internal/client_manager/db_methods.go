package client_manager

import (
	"context"
	"encoding/json"
	"fmt"
	"pokestocks/internal/helpers"
	"pokestocks/utils"
	"strconv"
	"strings"

	common_pb "pokestocks/proto/common"
	redis_keys "pokestocks/redis"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func (cc *ClientManager) queryDbForPokemonStockPairs(ctx context.Context, pspIds []string) ([]*common_pb.PokemonStockPair, error) {
	db := cc.DB
	redisClient := cc.RedisClient
	redisPipeline := redisClient.Pipeline()

	query := helpers.PspQueryString()

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
	midnightTomorrow := helpers.MidnightTomorrow()

	for rows.Next() {
		queriedData, err := pgx.RowToMap(rows)
		if err != nil {
			return nil, err
		}

		psp := helpers.ConvertDbRowToPokemonStockPair(queriedData)
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

func (cc *ClientManager) QueryPokemonStockPairs(ctx context.Context, pspIds []string) ([]*common_pb.PokemonStockPair, error) {
	redisClient := cc.RedisClient
	redisPipeline := redisClient.Pipeline()

	var psps []*common_pb.PokemonStockPair

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

	if len(cachedIds) == 0 && len(nonCachedIds) == 0 {
		// the above loop on results never occurred so the id arrays are both empty
		utils.LogWarningError("Error querying Redis key for cached PSPs. Falling back to db", err)
		psps, err = cc.queryDbForPokemonStockPairs(ctx, pspIds)
		if err != nil {
			return nil, err
		}
	} else if len(nonCachedIds) > 0 {
		// 1 or more cache misses require querying db
		nonCachedPspsArr, err := cc.queryDbForPokemonStockPairs(ctx, nonCachedIds)
		if err != nil {
			return nil, err
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

	return psps, nil
}

func (cm *ClientManager) QueryPortfolioCash(ctx context.Context, portfolioId int64) (float64, error) {
	db := cm.DB

	query := `
		SELECT cash
		FROM portfolios
		WHERE id = $1
	`
	var cashString string

	err := db.QueryRow(ctx, query, portfolioId).Scan(&cashString)
	if err != nil {
		return 0, err
	}

	cash, err := strconv.ParseFloat(cashString, 64)
	if err != nil {
		return 0, err
	}

	return cash, nil
}
