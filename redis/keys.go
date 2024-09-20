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

func MarketStatusKey() string {
	return keyPrefix + "#marketStatus"
}

func StockSymbolKey(symbol string) string {
	return keyPrefix + ":stockPrice#" + symbol
}

func GetSymbolFromStockSymbolKey(key string) string {
	return strings.Split(key, "#")[1]
}

func NextMarketOpenKey() string {
	return keyPrefix + "#nextMarketOpen"
}
