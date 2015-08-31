package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"testing"
	"github.com/stretchr/testify/assert"
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

func InsertProfessor(b *testing.B) {
	setup()
	for n := 0; n < b.N; n++ {
		insertProfessor(makeFullProfessor(), testDatabase)
	}
	tearDown()
}

func TestInsertNullProfessor(t *testing.T) {
	setup()
	params := Parameter{
		LastName:   "Test-Asami-Sato81",
		Department: "Badassery",
		City:       "Newark",
		CourseNumber:"103",

		IsRutgers:  true}

	insertOrUpdateMapping(params, 0, testDatabase)

	tearDown()
}

func TestGetMappingWithNullProfessor(t *testing.T) {
	setup()

	params := Parameter{
		LastName:   "Test-Asami-Sato81",
		Department: "Badassery",
		City:       "Newark",
		CourseNumber:"103",

		IsRutgers:  true}

	result ,err := SearchDatabase(params, testDatabase)

	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.True(t, result != nil)

	log.Printf("Result: %#s", result)

	tearDown()
}

func TestSearchWithNullProfessor(t *testing.T) {
	setup()

	params := Parameter{
		LastName:   "Test-Asami-Sato81",
		Department: "Badassery",
		City:       "Newark",
		CourseNumber:"103",

		IsRutgers:  true}

	result := Search(params, testDatabase)

	assert.True(t, result != nil)

	log.Printf("Result: %#s", result)

	tearDown()
}

func TestCheckNullProfessor(t *testing.T) {
	setup()

	params := Parameter{
		LastName:   "Test-Asami-Sato81",
		Department: "Badassery",
		City:       "Newark",
		CourseNumber:"103",

		IsRutgers:  true}

	assert.True(t, 	checkMappingExists(params, testDatabase))

	tearDown()
}

func TestInsertExclusions(t *testing.T) {
	setup()

	insertExclusions([]string{"/ShowRatings.jsp?tid=373482", "/ShowRatings.jsp?tid=1537234"}, testDatabase)

	tearDown()
}

func TestQueryExclusionsForMapping(t *testing.T) {
	setup()
	var url string
	rows := queryExclusionsForMapping(2, testDatabase)
	for rows.Next() {
		rows.Scan(&url)
	}
	log.Println("Result", url)
	assert.True(t, len(url) > 0)
	tearDown()
}

func TestRefreshDatabase(t *testing.T) {
	setup()

	RefreshDatabase(testDatabase)

	tearDown()
}

func TestIncrementStaleCount(t *testing.T) {
	setup()
	expected := int64(0)
	params := Parameter{
		LastName:   "Friedman",
		Department: "History",
		City:       "Newark",
		CourseNumber:"103",

		IsRutgers:  true}

	result, count := incrementStaleCount(params, testDatabase)
	assert.True(t, result > expected)
	assert.True(t, count > expected, )

	tearDown()
}

func TestConstructParametersFromRows(t *testing.T) {
	setup()

	parameters := constructParametersFromRows(queryStaleMappingsForUpdate(testDatabase), testDatabase)
	log.Printf("%#v", parameters)
	assert.True(t, len(parameters) > 0)

	tearDown()
}

func TestUpdateMappings(t *testing.T) {
	setup()

	rowId := updateMapping(1, 3, testDatabase)
	assert.True(t, rowId == 1)

	tearDown()
}

func TestQueryStaleMappingsForUpdate(t *testing.T) {
	setup()

	var lastName sql.NullString
	var dumb sql.NullString
	rows := queryStaleMappingsForUpdate(testDatabase)
	for rows.Next() {
		err := rows.Scan(&dumb, &dumb, &lastName, &dumb, &dumb, &dumb, &dumb)
		checkError(err)
		log.Println("Result", lastName)
	}
	assert.True(t, len(lastName.String) > 1)
	tearDown()
}

func TestDatabaseSearch(t *testing.T) {
	setup()

	params := Parameter{
		LastName:   "Friedman",
		Department: "History",
		City:       "Newark",
		CourseNumber:"103",

		IsRutgers:  true}

	result ,err := SearchDatabase(params, testDatabase)

	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.True(t, result != nil)

	log.Printf("Result: %#s", result)

	tearDown()
}
