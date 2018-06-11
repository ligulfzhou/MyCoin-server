package models

import (
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

// not stored in db, but in redis
type Coin struct {
	Id                   int    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Symbol               string `json:"symbol"`
	Name                 string `json:"name"`
	PriceUsd             string `json:"price_usd"`
	PriceBtc             string `json:"price_btc"`
	PercentChangeOneHour string `json:"percent_change_one_hour"`
	PercentChangeOneDay  string `json:"percent_change_one_day"`
	PercentChangeOneWeek string `json:"percent_change_one_week"`
	MarketCapUsd         string `json:"market_cap_usd"`
	OneDayVolumnUsd      string `json:"one_day_volumn_usd"`
	AvailableSupply      string `json:"available_supply"`
	TotalSupply          string `json:"total_supply"`
	ImgUrl               string `json:"img_url"`
}

type TCoin struct {
	Id     int    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	ImgUrl string `json:"img_url"`
}

func (c *TCoin) TableName() string {
	return "coin"
}

type UserCoin struct {
	Id     int    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	User   string `json:"user"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Cnt    int    `json:"cnt"`
}

type UserCoinPair struct {
	*Coin

	User   string `json:"user"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Cnt    int    `json:"cnt"`
}

func (uc *UserCoin) TableName() string {
	return "user_coin"
}

func Init() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:MYSQLzhouligang153@/xcoin?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("connect mysql error")
	}

	DB = db
	return db, err
}

func GetUserCoins(user string) ([]UserCoin, error) {
	userCoins := []UserCoin{}
	DB.Where("user = ?", user).Find(&userCoins)
	return userCoins, nil
}

func SetUserCoin(user string, symbol string, name string, count int) UserCoin {
	userCoin := UserCoin{}

	DB.Where("user = ? AND symbol = ? AND name = ?", user, symbol, name).Find(&userCoin)
	if userCoin.Symbol != "" {
		if userCoin.Cnt != count {
			userCoin.Cnt = count
			// DB.Model(&userCoin).Update("count", count)
			DB.Save(&userCoin)
		}
	} else {
		userCoin.User = user
		userCoin.Symbol = symbol
		userCoin.Name = name
		userCoin.Cnt = count
		DB.Create(&userCoin)
	}

	return userCoin
}

func SearchCoins(search string) []TCoin {
	coins := []TCoin{}
	searchTxt := "%" + search + "%"
	DB.Where("symbol like ? OR name like ?", searchTxt, searchTxt).Find(&coins)
	return coins
}
