package caches

import (
	"encoding/json"
	// "fmt"
	"github.com/gomodule/redigo/redis"
	"ligulfzhou.com/coincalc/models"
)

var REDIS redis.Conn

const (
	EXPIRE_TIME = 5 * 60
)

func getCoinsKey() string {
	return "coins"
}

func getCoinKey(symbolAndName string) string {
	return "coin_" + symbolAndName
}

func getUsersKey(user string) string {
	return "" + user
}

func Init() (redis.Conn, error) {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic("connect redis error")
	}

	REDIS = c
	return c, err
}

func GetCoins(page int, pagesize int) ([]models.Coin, error) {
	start, stop := (page-1)*pagesize, page*pagesize-1
	res, err := redis.Strings(REDIS.Do("lrange", getCoinsKey(), start, stop))
	if err != nil {
		return nil, err
	}

	coinKeys := []interface{}{}
	for _, coinSymbolAndName := range res {
		coinKeys = append(coinKeys, getCoinKey(coinSymbolAndName))
	}

	coinsr, err := redis.Strings(REDIS.Do("mget", coinKeys...))

	coins := []models.Coin{}
	for _, coin := range coinsr {
		tmp := models.Coin{}
		_ = json.Unmarshal([]byte(coin), &tmp)
		coins = append(coins, tmp)
	}

	return coins, nil
}

func GetUserCoins(user string) ([]models.UserCoinPair, error) {
	// todo: 是否应该将用户的币列表也放到redis里去？
	userCoins, _ := models.GetUserCoins(user)

	coinKeys := []interface{}{}
	for _, userCoin := range userCoins {
		coinKeys = append(coinKeys, getCoinKey(userCoin.Symbol+"_"+userCoin.Name))
	}

	coinsr, _ := redis.Strings(REDIS.Do("mget", coinKeys...))

	coins := []models.UserCoinPair{}
	for idx, coin := range coinsr {
		tmp := models.Coin{}
		_ = json.Unmarshal([]byte(coin), &tmp)
		coins = append(coins, models.UserCoinPair{&tmp, userCoins[idx].User, userCoins[idx].Symbol, userCoins[idx].Name, userCoins[idx].Cnt})
	}
	return coins, nil
}

func PostUserCoin(user string, symbol string, name string, count int) (models.UserCoinPair, error) {
	uc := models.SetUserCoin(user, symbol, name, count)
	coin := models.Coin{}
	coinRaw, _ := redis.String(REDIS.Do("get", getCoinKey(uc.Symbol+"_"+uc.Name)))
	_ = json.Unmarshal([]byte(coinRaw), &coin)

	return models.UserCoinPair{&coin, uc.User, uc.Symbol, uc.Name, uc.Cnt}, nil
}
