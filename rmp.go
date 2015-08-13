package main

import (
	"github.com/renstrom/fuzzysearch/fuzzy"
	"crypto/sha1"
	"fmt"
)

type (
	Professors []*Professor

	ProfessorsByName struct {
		FirstName  string
		LastName   string
		professors Professors
	}

	ProfessorsByDepartment struct {
		department string
		professors Professors
	}

	ProfessorsByCity struct {
		city       string
		professors Professors
	}

	Professor struct {
		FirstName   string   `json:"firstName,omitempty"`
		LastName    string   `json:"lastName,omitempty"`
		Email       string   `json:"email,omitempty"`
		Department  string   `json:"department,omitempty"`
		Title       string   `json:"title,omitempty"`
		PhoneNumber []string `json:"phoneNumber,omitempty"`
		FaxNumber   string   `json:"faxNumber,omitempty"`
		Location    Location `json:"location,omitempty"`
		Rating      Rating   `json:"rating,omitempty"`
	}

	Location struct {
		School  string `json:"university,omitempty"`
		City    string `json:"city,omitempty"`
		State   string `json:"state,omitempty"`
		Room    string `json:"room,omitempty"`
		Address string `json:"address,omitempty"`
	}

	Rating struct {
		Overall      float64 `json:"overall"`
		Helpfulness float64 `json:"helpfulness"`
		Easiness     float64 `json:"easiness"`
		Clarity      float64 `json:"clarity"`
		AverageGrade string  `json:"averageGrade"`
		Hotness      bool    `json:"hotness"`
		RatingsCount float64 `json:"ratingsCount"`
		RatingUrl    string  `json:"ratingUrl,omitempty"`
	}
)

func (p *Professors) Remove1(index int) {
	s := *p
	s = append(s[:index], s[index+1:]...)
	*p = s
}

func (p ProfessorsByName) Len() int {
	return len(p.professors)
}

func (p ProfessorsByName) Less(i, j int) bool {
	rank := p.compareProfessorName(p.professors[i], p.professors[j])
	if rank == 1 {
		return false
	} else if rank == -1 {
		return true
	}
	return false
}

func (p ProfessorsByName) Swap(i, j int) {
	p.professors[i], p.professors[j] = p.professors[j], p.professors[i]
}

func (p ProfessorsByName) compareProfessorName(prof1, prof2 *Professor) int {
	if len(p.FirstName) != 0 {
		p1 := l(string(prof1.FirstName[0]))
		p2 := l(string(prof2.FirstName[0]))
		param := l(string(p.FirstName[0]))

		if p1 == param && p2 != param {
			return -1
		} else if p1 != param &&
			p2 == param {
			return 1
		} else {
			return compareLength(p1, p2)
		}
	}
	return 0
}

func (p ProfessorsByDepartment) Len() int {
	return len(p.professors)
}

func (p ProfessorsByDepartment) Less(i, j int) bool {
	rank := p.compareProfessorDepartment(p.professors[i], p.professors[j])
	if rank == 1 {
		return false
	} else if rank == -1 {
		return true
	}
	return false
}

func (p ProfessorsByDepartment) Swap(i, j int) {
	p.professors[i], p.professors[j] = p.professors[j], p.professors[i]
}

func (p ProfessorsByDepartment) compareProfessorDepartment(prof1, prof2 *Professor) int {
	param := l(p.department)

	p1 := fuzzy.LevenshteinDistance(param, prof1.Department)
	p2 := fuzzy.LevenshteinDistance(param, prof2.Department)

	if p1 < p2 {
		return -1
	} else if p1 > p2 {
		return 1
	} else {
		return 1
	}
}

func (p ProfessorsByCity) Len() int {
	return len(p.professors)
}

func (p ProfessorsByCity) Less(i, j int) bool {
	rank := p.compareProfessorCity(p.professors[i], p.professors[j])
	if rank == 1 {
		return false
	} else if rank == -1 {
		return true
	}
	return false
}

func (p ProfessorsByCity) Swap(i, j int) {
	p.professors[i], p.professors[j] = p.professors[j], p.professors[i]
}

func (p ProfessorsByCity) compareProfessorCity(prof1, prof2 *Professor) int {
	p1 := l(prof1.Location.City)
	p2 := l(prof2.Location.City)
	param := l(p.city)

	if p1 == param && p2 != param {
		return -1
	} else if p1 != param &&
		p2 == param {
		return 1
	} else {
		return compareLength(p1, p2)
	}
}

func compareLength(s1, s2 string) int {
	if len(s1) < len(s2) {
		return 1
	} else {
		return -1
	}
}

func (n *Professor) FullName() string {
	return n.FirstName + " " + n.LastName
}

func (n *Professor) convertPhoneNumber() string {
	var tempStr string
	for _, val := range n.PhoneNumber {
		if val != "" {
			tempStr = tempStr + val + ","
		}
	}
	if len(tempStr) > 0 {
		return tempStr[:len(tempStr)-1]
	}
	return Empty
}

func (p *Professor) hash() string {
	hash := p.FirstName+p.LastName+p.Department+p.Location.City+p.Location.State
	sum := sha1.Sum([]byte(hash))
	//%x	base 16, lower-case a-f, two characters per byte
	return fmt.Sprintf("%x", sum[:8])
}

func (p *Professor) equals(other *Professor) bool {
	return p.hash() == other.hash()
}