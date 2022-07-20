package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var store = NewUrlStore()

func main() {
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
		c.Error(err)
	}
	c.Redirect(http.StatusFound, url)
}

func add(c *gin.Context) {
	var urlStruct struct {
		Url string `json:"url" binding:"required"`
	}
	c.BindJSON(&urlStruct)
	key, err := store.Put(urlStruct.Url)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, gin.H{"key": key})
}
