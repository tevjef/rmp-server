package main

import (
	"github.com/gin-gonic/gin"
)

func searchHandler(c *gin.Context) {
	first := l(c.Query("first"))           // ""
	last := l(c.Query("last"))             // Randall
	department := l(c.Query("department")) //Biology
	city := l(c.Query("city"))             //newark, new brunswick, camden

	if isEmpty(last) || isEmpty(department) || isEmpty(city) {
		c.String(400, "Must supply all paramters")
		return
	}
	//semester := c.Query("semester") //
	//year := c.Query("year") // 2015
	//courseNumber := c.Query("course") //198

	params := Parameter{
		FirstName:  first,
		LastName:   last,
		City:       city,
		Department: department,
		IsRutgers:  true}

	options := Options{
		FilterSearch:  true,
		RutgersSearch: true,
		SortSearch:    true}

	p := search(params, options)
	printProfs(p)
	c.JSON(200, p)
}

func searchDbHandler(c *gin.Context) {
	first := l(c.Query("first"))        // ""
	last := l(c.Query("last"))          // Randall
	department := l(c.Query("subject")) //Biology
	courseNumber := c.Query("course")   //198
	city := l(c.Query("city"))          //newark, new brunswick, camden

	if isEmpty(last) || isEmpty(department) || isEmpty(city) {
		c.String(400, "Must supply all paramters")
		return
	}

	params := Parameter{
		FirstName:    first,
		LastName:     last,
		City:         city,
		CourseNumber: courseNumber,
		Department:   department,
		IsRutgers:    true}

	options := Options{
		FilterSearch:  true,
		RutgersSearch: true,
		SortSearch:    true}

	p := search(params, options)
	printProfs(p)
	c.JSON(200, p)
}
