package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"log"
	"crypto/sha1"
)

const (
	SearchLimit       = 200
	DefaultSearchSize = 20
	BaseRMPUrl        = "http://www.ratemyprofessors.com"
)

type (
	Parameter struct {
		FirstName    string
		LastName     string
		Department   string
		City         string
		CourseNumber string
		Exclusion    []string
		Inclusion    string
		IsRutgers    bool
	}

	Options struct {
		SortSearch    bool
		FilterSearch  bool
		RutgersSearch bool
	}
)

func search(params Parameter, options Options) (professors Professors) {
	var wg sync.WaitGroup
	//Execute normal search if there is no explicit professor to search for
	if !params.hasInclusion() {
		//Initial search to get professors
		doc := execSearch(params.LastName, params.IsRutgers, 0)
		numOfProfessors := getNumberOfProfessors(doc)
		professors = appendProfessors(professors, extractListings(doc))

		//If more searches are necessary, make them
		if numOfProfessors > DefaultSearchSize {
			for offset := DefaultSearchSize; offset <= numOfProfessors; offset += DefaultSearchSize {
				wg.Add(1)
				go func(offset int) {
					doc := execSearch(params.LastName, params.IsRutgers, offset)
					professors = appendProfessors(professors, extractListings(doc))
					defer wg.Done()
				}(offset)
			}
		}
		wg.Wait()
	} else {
		//Since this was explicitly added it should not be filtered out if it doesn't match the filter's conditions
		options.FilterSearch = false
		professors = appendProfessors(professors, []*Professor{
			&Professor{
				Rating: Rating{
					RatingUrl: params.Inclusion,
				},
			}})
	}

	//Once finished discovering all possible professors remove listings that match our exclusion list
	professors = filterListings(params, professors)

	//Deeper search for professors info.
	for _, val := range professors {
		wg.Add(1)
		go func(prof *Professor) {
			execLookup(prof)
			wg.Done()
		}(val)
	}
	wg.Wait()

	//Filters out professors that don't match certain conditions.
	if options.FilterSearch {
		professors = filterProfessors(professors, params)
	}

	//Pull addition professor information from rutgers directory
	if options.RutgersSearch {
		searchRutgersProfessors(professors)
	}

	if options.SortSearch {
		sortProfessors(professors, params)
	}

	return
}

func filterListings(params Parameter, professors Professors) (filtered Professors) {
	for _, prof := range professors {
		appendFlag := true
		for _, excl := range params.Exclusion {
			if prof.Rating.RatingUrl == excl {
				appendFlag = false
			}
		}
		if appendFlag {
			filtered = append(filtered, prof)
		}
	}
	return
}

func sortProfessors(professors Professors, params Parameter) {
	byCity := ProfessorsByCity{city: params.City, professors: professors}
	sort.Stable(byCity)

	byDepartment := ProfessorsByDepartment{department: params.Department, professors: professors}
	sort.Stable(byDepartment)

	byName := ProfessorsByName{FirstName: params.FirstName, LastName: params.LastName, professors: professors}
	sort.Stable(byName)
}

func filterProfessors(professors Professors, params Parameter) (filtered Professors) {
	for _, val := range professors {
		if l(val.Location.City) == l(params.City) && l(val.LastName) == l(params.LastName) {
			filtered = append(filtered, val)
		}
	}
	return
}

func searchRutgersProfessors(professors Professors) {
	var wg sync.WaitGroup
	for _, val := range professors {
		wg.Add(1)
		go func(prof *Professor) {
			prof = execPeopleSearch(prof)
			wg.Done()
		}(val)
	}
	wg.Wait()
}

func printProfs(p Professors) {
	for i, val := range p {
		fmt.Printf("%d. %10s | %8s | %s | Title: %s\n", i, val.FullName(), val.Location.City, val.Department, val.Title)
	}
}

func printProf(p *Professor) string {
	return fmt.Sprintf("%20s | %8s | %s | Title: %s\n", p.FullName(), p.Location.City, p.Department, p.Title)
}

func appendProfessors(profs Professors, toAppend Professors) Professors {
	for _, val := range toAppend {
		profs = append(profs, val)
	}
	return profs
}

func execSearch(name string, isRutgers bool, offset int) *goquery.Document {
	uni := ""
	if isRutgers {
		uni = "+rutgers"
	}
	url := fmt.Sprintf(BaseRMPUrl+"/search.jsp?query=%s%s&stateselect=nj&offset=%d", name, uni, offset)
	doc, _ := goquery.NewDocument(url)
	return doc
}

func extractListings(doc *goquery.Document) (listings Professors) {
	doc.Find(".listing").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find("a").Attr("href")

		professor := new(Professor)
		professor.Rating.RatingUrl = url

		listings = append(listings, professor)
	})
	return
}

func getNumberOfProfessors(doc *goquery.Document) int {
	resultText := doc.Find(".result-count").Text()
	start := strings.LastIndex(resultText, "f") + 2
	end := strings.LastIndex(resultText, "r") - 1
	val, _ := strconv.Atoi(resultText[start:end])
	return val
}

func execLookup(professor *Professor) {
	url := BaseRMPUrl + professor.Rating.RatingUrl
	log.Println("Lookup:", url)
	doc, _ := goquery.NewDocument(url)
	extractProfessor(professor, doc)
}

func extractProfessor(professor *Professor, doc *goquery.Document) {
	professor.FirstName = extractFirstName(doc)
	professor.LastName = extractLastName(doc)

	professor.Department = extractDepartment(doc)

	professor.Rating.AverageGrade = extractAverageGrade(doc)
	professor.Rating.Clarity = extractClarity(doc)
	professor.Rating.Easiness = extractEasiness(doc)
	professor.Rating.Helpfullness = extractHelpfulness(doc)
	professor.Rating.Hotness = extractHotness(doc)
	professor.Rating.Overall = extractOverall(doc)

	professor.Location.School = extractUniversity(doc)
	professor.Location.City = extractCity(doc)
	professor.Location.State = extractState(doc)
}

func extractOverall(doc *goquery.Document) float64 {
	result := doc.Find(".breakdown-header").First().Find(".grade").Text()
	val, _ := strconv.ParseFloat(result, 64)
	return val
}

func extractHelpfulness(doc *goquery.Document) float64 {
	result := doc.Find(".rating-slider").First().Find(".rating").Text()
	val, _ := strconv.ParseFloat(result, 64)
	return val
}

func extractClarity(doc *goquery.Document) float64 {
	result := doc.Find(".rating-slider").First().Next().Find(".rating").Text()
	val, _ := strconv.ParseFloat(result, 64)
	return val
}

func extractEasiness(doc *goquery.Document) float64 {
	result := doc.Find(".rating-slider").First().Next().Next().Find(".rating").Text()
	val, _ := strconv.ParseFloat(result, 64)
	return val
}

func extractAverageGrade(doc *goquery.Document) string {
	result := doc.Find(".breakdown-header").Next().Find(".grade").Text()
	return format(result)
}

func extractHotness(doc *goquery.Document) bool {
	result, _ := doc.Find(".breakdown-header").Next().Next().Html()
	return !strings.Contains(result, "cold")
}

func extractRatingsCount(doc *goquery.Document) int {
	result := format(doc.Find(".rating-count").Text())
	result = substringBefore(result, " ")
	resultInt, err := strconv.Atoi(result)
	checkError(err)
	return resultInt
}

func extractCity(doc *goquery.Document) string {
	result := doc.Find(".result-title").Text()
	result = substringAfter(result, ", ")
	result = substringBefore(result, ",")
	return strings.TrimSpace(result)
}

func extractState(doc *goquery.Document) string {
	result := doc.Find(".result-title").Text()
	dirty := substringAfterLast(result, ", ")
	return strings.TrimSpace(dirty)
}

func extractFirstName(doc *goquery.Document) string {
	result := format(doc.Find(".pfname").First().Text())
	return format(result)
}

func extractLastName(doc *goquery.Document) string {
	result := format(doc.Find(".plname").First().Text())
	return format(result)
}

func extractDepartment(doc *goquery.Document) string {
	result := doc.Find(".result-title").Text()
	dirty := substringAfter(result, "Professor in the ")
	dirty = substringBefore(dirty, " department")
	return format(dirty)
}

func extractUniversity(doc *goquery.Document) string {
	result := doc.Find(".result-title").Find(".school").Text()
	return strings.TrimSpace(result)
}

func execPeopleSearch(professor *Professor) *Professor {
	doc := getPeopleSearchDocument(professor)

	professor.Email = extractEmail(doc)
	professor.Title = extractTitle(doc)

	professor.PhoneNumber = []string{extractPhone1(doc), extractPhone2(doc)}
	professor.FaxNumber = extractFax(doc)
	professor.Location.Address = extractAddress(doc)
	professor.Location.Room = extractRoomLocation(doc)

	return professor
}

func getPeopleSearchDocument(professor *Professor) *goquery.Document {
	values := make(map[string][]string)
	values["p_name_last"] = []string{professor.LastName}
	values["p_name_first"] = []string{professor.FirstName}
	resp, _ := http.PostForm("https://www.acs.rutgers.edu/pls/pdb_p/Pdb_Display.search_results", values)

	doc, _ := goquery.NewDocumentFromResponse(resp)
	return doc
}

func extractTitle(doc *goquery.Document) string {
	html, _ := doc.Html()
	doc, _ = goquery.NewDocumentFromReader(bytes.NewBufferString(substringAfter(html, "<dt>Title:</dt>")))
	result := doc.Find("dd").First().Text()
	return format(result)
}

func extractPhone1(doc *goquery.Document) string {
	html, _ := doc.Html()
	dirty := substringAfter(html, "<dt>Phone:</dt>")
	dirty = substringBefore(dirty, "<h4>Email Address</h4>")
	doc, _ = goquery.NewDocumentFromReader(bytes.NewBufferString(dirty))
	result := doc.Find("dd").First().Text()
	if len(result) == 14 {
		return result
	}
	return Empty
}

func extractPhone2(doc *goquery.Document) string {
	html, _ := doc.Html()
	dirty := substringAfter(html, "<dt>Phone:</dt>")
	dirty = substringAfter(dirty, "<dt>Phone:</dt>")
	dirty = substringBefore(dirty, "<h4>Email Address</h4>")

	doc, _ = goquery.NewDocumentFromReader(bytes.NewBufferString(dirty))

	result := doc.Find("dd").First().Text()
	if len(result) == 14 {
		return result
	}
	return Empty
}

func extractFax(doc *goquery.Document) string {
	html, _ := doc.Html()
	dirty := substringAfter(html, "<dt>Phone:</dt>")
	dirty = substringBefore(dirty, "<h4>Email Address</h4>")

	doc, _ = goquery.NewDocumentFromReader(bytes.NewBufferString(dirty))
	result := doc.Find("dd").Last().Text()
	if len(result) == 14 {
		return result
	}
	return Empty
}

func extractAddress(doc *goquery.Document) string {
	html, _ := doc.Html()
	dirty := substringAfter(html, "<h4>Postal Address</h4>")
	dirty = substringBefore(dirty, "<h4>Location</h4>")

	doc, _ = goquery.NewDocumentFromReader(bytes.NewBufferString(dirty))

	result := doc.Find("dd").Last().Text()
	result = strings.Replace(result, "\u00a0", " | ", -1)
	result = strings.Replace(result, "\n", " ", -1)
	return strings.TrimSpace(result)
}

func extractRoomLocation(doc *goquery.Document) string {
	html, _ := doc.Html()
	dirty := substringAfter(html, "<h4>Location</h4>")
	dirty = substringBefore(dirty, "<form")

	doc, _ = goquery.NewDocumentFromReader(bytes.NewBufferString(dirty))

	result := doc.Find("dd").Last().Text()
	//log.Printf("Log %#v", result)
	result = strings.Replace(result, "\u00a0", " | ", -1)
	result = strings.Replace(result, "\n", " ", -1)
	return strings.TrimSpace(result)
}

func extractEmail(doc *goquery.Document) string {
	html, _ := doc.Html()
	dirty := substringAfter(html, "<h4>Email Address</h4>")
	dirty = substringBefore(dirty, "<h4>Postal Address</h4>")
	doc, _ = goquery.NewDocumentFromReader(bytes.NewBufferString(dirty))
	return doc.Find("dd").Find("a").First().Text()
}

func (p *Parameter) hasInclusion() bool {
	return strings.Contains(p.Inclusion, "/ShowRatings.jsp")
}

func (p *Parameter) hash() string {
	hash := p.LastName+p.Department+p.City+p.CourseNumber
	sum := sha1.Sum([]byte(hash))
	//%x	base 16, lower-case a-f, two characters per byte
	return fmt.Sprintf("%x", sum[:8])
}