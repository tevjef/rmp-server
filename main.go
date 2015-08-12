package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/search", searchHandler)
	r.POST("/search", reportHandler)
	r.Run(":8080")

	defer database.Close()
}
