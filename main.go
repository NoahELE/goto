package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var store, err = NewUrlStore()

func main() {
	if err != nil {
		log.Fatal(err)
	}
	r := setupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/")
	})
	r.Static("/web", "./frontend/dist")

	r.GET("/:key", redirect)
	r.POST("/add", add)

	return r
}

func redirect(c *gin.Context) {
	key := c.Params.ByName("key")
	url, err := store.Get(key)
	if err != nil {
		err := c.Error(err)
		if err != nil {
			return
		}
	}
	c.Redirect(http.StatusFound, url)
}

func add(c *gin.Context) {
	var urlStruct struct {
		Url string `json:"url" binding:"required"`
	}
	err := c.BindJSON(&urlStruct)
	if err != nil {
		log.Fatal(err)
	}
	key, err := store.Put(urlStruct.Url)
	if err != nil {
		err := c.Error(err)
		if err != nil {
			log.Fatal(err)
		}
	}
	c.JSON(http.StatusOK, gin.H{"key": key})
}
