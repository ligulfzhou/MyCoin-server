package main

import (
	"github.com/gin-gonic/gin"
	"ligulfzhou.com/coincalc/caches"
	"ligulfzhou.com/coincalc/controllers"
	"ligulfzhou.com/coincalc/models"
)

func main() {
	db, _ := models.Init()
	defer db.Close()

	rs, _ := caches.Init()
	defer rs.Close()

	r := gin.Default()
	r.GET("/coins", controllers.GetCoins)
	r.GET("/user/coins", controllers.GetUserCoin)
	r.POST("/user/coins", controllers.PostUserCoin)
	r.GET("/coins/search", controllers.SearchCoins)

	r.Run("0.0.0.0:8050")
}
