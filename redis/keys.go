package redis

const (
	keyPrefix = "POKESTOCKS"
)

func ElasticCacheKey(searchValue string) string {
	return keyPrefix + ":elastic#" + searchValue
}
