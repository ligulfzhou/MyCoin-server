package controllers

import (
	"github.com/gin-gonic/gin"
	"ligulfzhou.com/coincalc/caches"
	"ligulfzhou.com/coincalc/models"
	"net/http"
	"strconv"
)

func GetCoins(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	pagesize, err2 := strconv.Atoi(c.DefaultQuery("pagesize", "100"))

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "params error",
		})
		return
	}

	coins, err := caches.GetCoins(page, pagesize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "db error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"coins": coins,
	})
}

func GetUserCoin(c *gin.Context) {
	user := c.DefaultQuery("user", "")
	if user == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "params error",
		})
		return
	}
	userCoins, err := caches.GetUserCoins(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "db error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"coins": userCoins,
	})
	return
}

func PostUserCoin(c *gin.Context) {
	user := c.DefaultQuery("user", "")
	count, err := strconv.Atoi(c.DefaultQuery("count", "0"))
	symbol := c.DefaultQuery("symbol", "")
	name := c.DefaultQuery("name", "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "params error",
		})
		return
	}

	ucp, _ := caches.PostUserCoin(user, symbol, name, count)
	c.JSON(http.StatusOK, gin.H{
		"user_coin": ucp,
	})
}

func SearchCoins(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "params error",
		})
		return
	}

	coins := models.SearchCoins(search)
	c.JSON(http.StatusOK, gin.H{
		"coins": coins,
	})
}
