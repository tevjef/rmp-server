package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

/*func TestProfessorsToJson(t *testing.T) {
	sql.NullString{}
	p := makeProfessors()
	out, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(out))
	}

	return
}*/

func TestSortByCity(t *testing.T) {
	expected := "Newark"
	p := makeProfessors()
	byCity := ProfessorsByCity{city: expected, professors: p}
	sort.Stable(byCity)
	t.Logf("\n\nSort by City\n")
	t.Logf(printTestProf(p))

	assert.True(t, p[0].Location.City == expected)
	assert.True(t, p[1].Location.City == expected)
	assert.True(t, p[2].Location.City == expected)
	assert.False(t, p[3].Location.City == expected)
	assert.False(t, p[4].Location.City == expected)

	return
}

func TestSortByDepartment(t *testing.T) {
	expected := "Biology"
	p := makeProfessors()
	byDepartment := ProfessorsByDepartment{department: expected, professors: p}
	sort.Stable(byDepartment)
	t.Logf("\n\nSort by Department\n")
	t.Logf(printTestProf(p))

	assert.True(t, p[0].Department == expected)
	assert.False(t, p[1].Department == expected)
	assert.False(t, p[2].Department == expected)
	assert.False(t, p[3].Department == expected)
	assert.False(t, p[4].Department == expected)

	return
}

func TestSortByName(t *testing.T) {
	expected := "D"
	p := makeProfessors()
	byName := ProfessorsByName{FirstName: expected, professors: p}
	sort.Stable(byName)
	t.Logf("\n\nSort by Name\n")
	t.Logf(printTestProf(p))

	assert.True(t, string(p[0].FirstName[0]) == expected)
	assert.True(t, string(p[1].FirstName[0]) == expected)
	assert.False(t, string(p[2].FirstName[0]) == expected)
	assert.False(t, string(p[3].FirstName[0]) == expected)
	assert.False(t, string(p[4].FirstName[0]) == expected)

	return
}

func TestSortByAll(t *testing.T) {
	expected := "Biology"
	p := makeProfessors()

	byCity := ProfessorsByCity{city: "Newark", professors: p}
	sort.Stable(byCity)
	t.Logf("\n\nSort by City\n")
	t.Logf(printTestProf(p))

	byDepartment := ProfessorsByDepartment{department: "Biology", professors: p}
	sort.Stable(byDepartment)
	t.Logf("\n\nSort by Department\n")
	t.Logf(printTestProf(p))

	byName := ProfessorsByName{FirstName: "", professors: p}
	sort.Stable(byName)
	t.Logf("\n\nSort by Name\n")
	t.Logf(printTestProf(p))

	assert.True(t, string(p[0].Department) == expected)

	return
}

func makeProfessors() (professors Professors) {
	p1 := &Professor{FirstName: "Douglas", LastName: "Morrison", Location: Location{City: "Newark"}, Department: "Science"}
	p3 := &Professor{FirstName: "Karl", LastName: "Morrison", Location: Location{City: "New Brunswick"}, Department: "History"}
	p2 := &Professor{FirstName: "Debrorah", LastName: "Morrison", Location: Location{City: "Newark"}, Department: "History"}
	p4 := &Professor{FirstName: "Brittany", LastName: "Morrison", Location: Location{City: "Camden"}, Department: "Biology"}
	p5 := &Professor{FirstName: "Victoria", LastName: "Morrison", Location: Location{City: "Newark"}, Department: "Law"}

	professors = append(professors, p1, p2, p3, p4, p5)
	return
}

func makeFullProfessor() (p *Professor) {
	return &Professor{
		FirstName:   "Victoria",
		LastName:    "Morrison",
		Department:  "Science",
		Title:       "ASSOCIATE PROFESSOR",
		Email:       "randall@rutgers.edu",
		PhoneNumber: []string{"(347) 320 4109", "(347) 320 7422", "(347) 123 4567"},
		FaxNumber:   "(347) 320 4109",
		Location: Location{
			School:  "Rutgers State University of New Jersey",
			City:    "Newark",
			State:   "NJ",
			Room:    "Room 205",
			Address: "101 Warren Streent",
		},
		Rating: Rating{
			Overall:      4.5,
			Helpfullness: 4.1,
			Clarity:      4.9,
			Easiness:     4.9,
			AverageGrade: "F",
			Hotness:      true,
			RatingUrl:    "SOMeurl.com",
			RatingsCount: 21,
		},
	}
}
func printTestProf(p Professors) string {
	s := make([]string, 1)
	s = append(s, "\n")
	for i, val := range p {
		s = append(s, fmt.Sprintf("%d. %20s | %8s | %s\n", i, val.FullName(), val.Location.City, val.Department))
	}
	return fmt.Sprintf("%s", s)
}
