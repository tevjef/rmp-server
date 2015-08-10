package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/search", searchHandler)
	r.GET("/searchdb", searchDbHandler)
	r.Run(":8080")
}
