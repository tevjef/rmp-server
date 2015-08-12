package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"fmt"
)

func searchHandler(c *gin.Context) {
	params, err := defaultHandler(c)
	if err == nil {
		p := SearchServers(params, database)
		if (p == nil) {
			params.IsRutgers = false
			p = SearchServers(params, database)
		}
		log.Printf("%#v", p)
		c.JSON(200, p)
	} else {
		c.String(400, err.Error())
	}
}

func reportHandler(c *gin.Context) {
	params, err := defaultHandler(c)
	if err == nil {
		_, count := incrementStaleCount(params, database)
		c.String(200, fmt.Sprintf("Report count: %d", count))
	} else {
		c.String(400, err.Error())
	}
}

func validateParams(param Parameter) error {
	if isEmpty(param.LastName) || isEmpty(param.Department) || isEmpty(param.City) || isEmpty(param.CourseNumber){
		log.Printf("%#v", param)
		return fmt.Errorf("Must supply all parameters. Params = %#v", param)
	}
	return nil
}

func parseQueryBool(queryString string) bool {
	if (l(queryString) == "true" || queryString == "1") {
		return true
	}
	return false
}

func defaultHandler(c *gin.Context) (Parameter, error) {
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

	err := validateParams(params)
	if err != nil {
		return params, err
	}

	return params, nil
}