package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"fmt"
)

func searchHandler(c *gin.Context) {
	params := defaultHandler(c)
	p := SearchServers(params, database)
	if (p == nil) {
		params.IsRutgers = false
		p = SearchServers(params, database)
	}
	log.Printf("%#v",p)
	c.JSON(200, p)
}

func reportHandler(c *gin.Context) {
	params := defaultHandler(c)
	_, count := incrementStaleCount(params,database)
	c.String(200, fmt.Sprintf("Report count: %d", count))
}

func validateParams(c *gin.Context, param Parameter) bool {
	if isEmpty(param.LastName) || isEmpty(param.Department) || isEmpty(param.City) || isEmpty(param.CourseNumber){
		log.Printf("%#v", param)
		c.String(400, fmt.Sprintf("Must supply all parameters. Params = %#v", param))
		return true
	}
	return false
}

func parseQueryBool(queryString string) bool {
	if (l(queryString) == "true" || queryString == "1") {
		return true
	}
	return false
}

func defaultHandler(c *gin.Context) Parameter {
	last := l(format(c.Query("last")))          // Randall
	department := l(format(c.Query("subject"))) //Biology
	city := l(format(c.Query("city"))  )        //newark, new brunswick, camden
	first := l(format(c.Query("first")) )       // ""
	courseNumber := format(c.Query("course"))   //198
	isRutgers := parseQueryBool(format(c.DefaultQuery("rutgers", "1")))

	params := Parameter{
		FirstName:    first,
		LastName:     last,
		City:         city,
		CourseNumber: courseNumber,
		Department:   department,
		IsRutgers:    isRutgers}
	if validateParams(c, params) {
		return params
	}
	return params
}