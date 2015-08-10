package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"testing"
)

var testDatabase *sql.DB

func setup() {
	dbInfo := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		DbUser, DbPassword, DbHost, DbName)

	db, err := sql.Open("postgres", dbInfo)
	testDatabase = db
	checkError(err)
}

func tearDown() {
	testDatabase.Close()
}

func TestInsertProfessor(t *testing.T) {
	setup()

	insertProfessor(makeFullProfessor(), testDatabase)

	tearDown()
}

func TestQueryProfessorById(t *testing.T) {
	setup()

	id, result, err := getProfessorFromRow(queryProfessorMappingById(3, testDatabase))
	checkError(err)
	log.Printf("ID: %d %#v ", id, result)

	tearDown()
}

func BenchmarkInsertProfessor(b *testing.B) {
	setup()
	for n := 0; n < b.N; n++ {
		insertProfessor(makeFullProfessor(), testDatabase)
	}
	tearDown()
}

func TestInsertExclusions(t *testing.T) {
	setup()

	insertExclusions([]string{"/ShowRatings.jsp?tid=373482", "/ShowRatings.jsp?tid=1537234"}, testDatabase)

	tearDown()
}

func TestServerSearch(t *testing.T) {
	setup()

	params := Parameter{
		LastName:   "Kapadia",
		Department: "Math",
		City:       "Newark",
		CourseNumber:"101",

		IsRutgers:  true}

	result := SearchServers(params, testDatabase)

	log.Printf("%#s", result)

	tearDown()
}
