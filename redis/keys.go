package redis_keys

import "strings"

const (
	keyPrefix = "POKESTOCKS"
)

func ElasticCacheKey(searchValue string) string {
	return keyPrefix + ":elastic#" + searchValue
}

func DbCacheKey(id string) string {
	return keyPrefix + ":db#" + id
}

func GetIdFromDbCacheKey(key string) string {
	return strings.Split(key, "#")[1]
}
