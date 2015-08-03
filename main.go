package main

import (
	"github.com/gin-gonic/gin"
)

var DB = make(map[string]string)

func main() {
	DB["Tevin"] = "Tevin Jeffrey"
	DB["Shariis"] = "Shariis Jeffrey"
	DB["Tyrece"] = "Tyrece Jeffrey"

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		if value, ok := DB[c.Query("name")]; ok {
			c.String(200, "Congrats " + value)
		} else {
			c.String(200, "User not found")
		}
	})
	r.Run(":8080")
}