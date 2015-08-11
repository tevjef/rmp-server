package main

import (
	"github.com/gin-gonic/gin"
	"log"
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
	last := l(format(c.Query("last")))          // Randall
	department := l(format(c.Query("subject"))) //Biology
	city := l(format(c.Query("city"))  )        //newark, new brunswick, camden
	first := l(format(c.Query("first")) )       // ""
	courseNumber := format(c.Query("course"))   //198

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

	log.Printf("%#v", params)
	p := SearchServers(params, database)

	if (p == nil) {
		params.IsRutgers = false
		p = SearchServers(params, database)
	}

	log.Printf("%#v",p)
	c.JSON(200, p)
}
