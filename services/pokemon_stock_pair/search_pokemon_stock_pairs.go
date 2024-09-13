package pokemon_stock_pair

import (
	"context"
	"fmt"
	common_pb "pokestocks/proto/common"
	psp_pb "pokestocks/proto/pokemon_stock_pair"
	"pokestocks/redis"
	"strings"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SearchPokemonStockPairs(ctx context.Context, in *psp_pb.SearchPokemonStockPairsRequest) (*psp_pb.SearchPokemonStockPairsResponse, error) {
	db := s.DB
	redisClient := s.RedisClient

	searchValue := in.SearchValue

	var ids []string

	elasticKeyInstances, _ := redisClient.Exists(ctx, redis.ElasticCacheKey(searchValue)).Result()
	if elasticKeyInstances == 1 {
		cachedIds, err := redisClient.SMembers(ctx, redis.ElasticCacheKey(searchValue)).Result()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error querying Redis for key %v: %v", redis.ElasticCacheKey(searchValue), err)
		}

		ids = cachedIds
	} else {
		searchResults, err := s.searchElasticIndex(searchValue)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error searching Elastic data: %v", err)
		}

		elasticPsps, err := convertPokemonStockPairElasticDocuments(searchResults)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error formatting Elastic data: %v", err)
		}

		ids = extractPokemonStockPairIds(elasticPsps)

		redisPipeline := redisClient.Pipeline()
		redisPipeline.SAdd(ctx, redis.ElasticCacheKey(searchValue), ids)
		redisPipeline.Expire(ctx, redis.ElasticCacheKey(searchValue), time.Second*10)

		_, err = redisPipeline.Exec(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error adding Redis data for key %v: %v", redis.ElasticCacheKey(searchValue), err)
		}
	}

	if len(ids) == 0 {
		return &psp_pb.SearchPokemonStockPairsResponse{Data: nil}, nil
	}

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

	for i, id := range ids {
		queryArgs = append(queryArgs, id)
		positionalParams = append(positionalParams, fmt.Sprintf("$%d", i+1))
	}

	orderByArgs := strings.Join(ids, ",")
	queryArgs = append(queryArgs, orderByArgs)

	positionalParamsString := strings.Join(positionalParams, ", ")
	query += fmt.Sprintf("WHERE psp.id IN (%s)", positionalParamsString)
	query += fmt.Sprintf("ORDER BY POSITION(psp.id::text IN $%d)", len(positionalParams)+1)

	rows, err := db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while querying for %v: %v", searchValue, err)
	}
	defer rows.Close()

	var psps []*common_pb.PokemonStockPair

	for rows.Next() {
		queriedData, err := pgx.RowToMap(rows)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error converting queried data to map: %v", err)
		}

		psp := convertDbRowToPokemonStockPair(queriedData)
		psps = append(psps, psp)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading queried data: %v", err)
	}

	err = s.enrichWithStockPrices(psps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)

	}
	return &psp_pb.SearchPokemonStockPairsResponse{Data: psps}, nil
}
