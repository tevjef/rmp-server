package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/search", searchHandler)
	r.GET("/report", reportHandler)
	r.Run(":8080")

	defer database.Close()
}
