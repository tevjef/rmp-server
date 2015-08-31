package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"log"
	"strings"
)

const (
	SearchListing    = "testRes/search_dummydata.txt"
	ProfessorListing = "testRes/professor_dummydata.txt"
	RutgersSearch    = "testRes/people_search_dummydata.txt"
)

func TestSearch(t *testing.T) {
	params := Parameter{
		LastName:   "watrous-",
		Department: "Computer Science",
		City:       "Newark",
		IsRutgers:  false}

	result := search(params)

	t.Logf("Result: %#s", result)
	assert.True(t, result != nil)
}

func TestNJITSearch(t *testing.T) {
	params := Parameter{
		FirstName:"L",
		LastName:   "Lay",
		Department: "Computer Science",
		City:       "Newark",
		IsRutgers:  false}

	result := search(params)

	t.Logf("Result: %#s", result)
	assert.True(t, result != nil)
	return
}

func TestSearchImpossible(t *testing.T) {
	params := Parameter{
		LastName:   "Asami-Sato",
		Department: "Computer Science",
		City:       "Newark",
		IsRutgers:  false}

	result := search(params)

	t.Logf("Result: %#s", result)
	assert.True(t, result != nil)
}


func TestSearchWithInclusion(t *testing.T) {
	expectedFirstName := "David"
	expectedCity := "New Brunswick"

	params := Parameter{
		Inclusion: "/ShowRatings.jsp?tid=1284875"}

	result := search(params)

	t.Logf("Result: %#s", result)
	assert.True(t, result != nil)
	assert.Equal(t, expectedFirstName, result.FirstName)
	assert.Equal(t, expectedCity,result.Location.City)

	return
}

func TestSearchWithExclusion(t *testing.T) {
	expectedFirstName := "Douglas"

	params := Parameter{
		LastName:   "Morrison",
		Department: "Biology",
		City:       "Newark",
		Exclusion:  []string{"/ShowRatings.jsp?tid=1834574", "/ShowRatings.jsp?tid=1208552"},
		IsRutgers:  true}

	result := search(params)

	t.Logf("Result: %#s", result)
	assert.Equal(t, expectedFirstName, result.FirstName)
	return
}

func TestFilterListings(t *testing.T) {
	params := Parameter{
		Exclusion: []string{"/ShowRatings.jsp?tid=126031"},
	}
	listings := extractListings(getDummyDoc(SearchListing))
	expected := len(listings) - 1
	listings = filterListings(params, listings)
	assert.True(t, expected == len(listings))
}

func TestSearchMultiPage(t *testing.T) {
	params := Parameter{
		LastName:  "John",
		City:      "Newark",
		IsRutgers: true}

	result := search(params)
	t.Logf("Result: %#s", result)
	return
}

func TestSortProfessors(t *testing.T) {
	expected := "Biology"
	p := makeProfessors()
	params := Parameter{City: "Newark", FirstName: "", LastName: "Morrison", Department: expected}
	sortProfessors(p, params)
	assert.True(t, string(p[0].Department) == expected)
	return
}

func TestFilterProfessors(t *testing.T) {
	expected := "Newark"
	p := makeProfessors()
	params := Parameter{City: "Newark", FirstName: "", LastName: "Morrison", Department: expected}
	p = filterProfessors(p, params)
	for _, val := range p {
		assert.Equal(t, val.Location.City, expected)
	}

	return
}

func TestExtractListings(t *testing.T) {
	expected := 13
	listings := extractListings(getDummyDoc(SearchListing))
	result := len(listings)
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestGetNumberOfProfessors(t *testing.T) {
	expected := 13
	result := getNumberOfProfessors(getDummyDoc(SearchListing))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractDepartment(t *testing.T) {
	expected := "Science"
	result := extractDepartment(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractUniversity(t *testing.T) {
	expected := "Rutgers - State University of New Jersey"
	result := extractUniversity(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractCity(t *testing.T) {
	expected := "Newark"
	result := extractCity(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractState(t *testing.T) {
	expected := "NJ"
	result := extractState(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractFirstName(t *testing.T) {
	expected := "Douglas"
	result := extractFirstName(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractLastName(t *testing.T) {
	expected := "Morrison"
	result := extractLastName(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractOverall(t *testing.T) {
	result := extractOverall(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.True(t, result > 1)
}

func TestExtractHelpfulness(t *testing.T) {
	result := extractHelpfulness(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.True(t, result > 1)
}

func TestExtractClarity(t *testing.T) {
	result := extractClarity(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.True(t, result > 1)
}

func TestExtractEasiness(t *testing.T) {
	result := extractEasiness(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.True(t, result > 1)
}

func TestExtractAverageGrade(t *testing.T) {
	expected := "C"
	result := extractAverageGrade(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractHotness(t *testing.T) {
	result := extractHotness(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.False(t, result)
}

func TestExtractRatingsCount(t *testing.T) {
	var expected float64
	expected = 142
	result := extractRatingsCount(getDummyDoc(ProfessorListing))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExecutePeopleSearch(t *testing.T) {
	expected := 1
	result := execPeopleSearch(makeProfessors()[0])
	t.Log("Result:", result)
	assert.True(t, len(result.Email) > expected)
	assert.True(t, len(result.Title) > expected)
	assert.True(t, len(result.Location.Address) > expected)
	assert.True(t, len(result.Location.Room) > expected)

}

func TestExtractTitle(t *testing.T) {
	expected := strings.Title(l("ASSOCIATE PROFESSOR"))
	result := extractTitle(getDummyDoc(RutgersSearch))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractPhone1(t *testing.T) {
	expected := "(973) 353-1268"
	result := extractPhone1(getDummyDoc(RutgersSearch))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractPhone2(t *testing.T) {
	expected := "(973) 353-5347"
	result := extractPhone2(getDummyDoc(RutgersSearch))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractFax(t *testing.T) {
	expected := "(973) 353-5518"
	result := extractFax(getDummyDoc(RutgersSearch))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func TestExtractAddress(t *testing.T) {
	result := extractAddress(getDummyDoc(RutgersSearch))
	log.Printf("Result: %#v", result)

	//t.Log("Result:", result)
	assert.True(t, len(result) > 1)
}

func TestExtractRoomLocation(t *testing.T) {
	result := extractRoomLocation(getDummyDoc(RutgersSearch))
	t.Log("Result:", result)
	assert.True(t, len(result) > 1)
}

func TestExtractRoomNumber(t *testing.T) {
	result := extractRoomNumber(getDummyDoc(RutgersSearch))
	t.Log("Result:", result)
	assert.True(t, len(result) > 1)
}

func TestExtractEmail(t *testing.T) {
	expected := "dougmorr@andromeda.rutgers.edu"
	result := extractEmail(getDummyDoc(RutgersSearch))
	t.Log("Result:", result)
	assert.Equal(t, expected, result)
}

func getDummyDoc(filename string) *goquery.Document {
	file, _ := ioutil.ReadFile(filename)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(file))
	return doc
}
