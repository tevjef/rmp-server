package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

var database *sql.DB

type Mapping struct {
	Parameter
	Fresh bool
}

func init() {
	dbInfo := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		DbUser, DbPassword, DbHost, DbName)

	db, err := sql.Open("postgres", dbInfo)
	database = db
	checkError(err)
}

func SearchServers(params Parameter, db *sql.DB) (professor *Professor) {
	var err error
	var professorId int

	if professor == nil {
		professorId, professor, _ = getProfessorFromRow(queryProfessorMappingByParams(params, db))
		log.Printf("ID: %d SEARCH DIRECT: %#s\n\n",professorId, professor)
	}
	if professor == nil {
		professorId, professor, err = getProfessorFromRow(queryAdjacentMappingsByParams(params, db))
		if err != nil && professorId != -1 {
			insertMapping(params, professorId, db)
			log.Printf("INSERTING ADJACENT")

		}
		log.Printf("ID: %d SEARCH ADJACENT: %#s\n\n",professorId, professor)
	}
	if professor == nil {
		options := Options{
			FilterSearch:  true,
			RutgersSearch: true,
			SortSearch:    true}

		professors := search(params, options)

		if len(professors) > 0 {
			professor = professors[0]
		}

		log.Printf("SEARCH RESULTS: %#s\n\n", professor)


		if professor != nil {
			professorId := insertProfessor(professor, db)
			exclusionIds := insertExclusions(params.Exclusion, db)
			mappingId := insertMapping(params, professorId, db)
			insertMappingExclusions(mappingId, exclusionIds, db)
			_, professor, _ = getProfessorFromRow(queryProfessorMappingById(professorId, db))
			log.Printf("ID: %d RETURNING AFTER INSERT PROFESSOR: %#s\n\n", professorId, professor)
		}
	}
	return professor
}

func insertProfessor(p *Professor, db *sql.DB) (professorId int) {
	var hash string
	db.QueryRow(`SELECT professor_id, hash
		FROM professors
		WHERE professors.hash = $1`, p.hash()).
	Scan(&professorId, &hash)

	if hash != p.hash() {
		err := db.QueryRow(
			`INSERT INTO professors
			(first_name, last_name, department, title, email, phone_number, fax_number, school, state,
			 city, room, address, overall, helpfulness, clarity, easiness, average_grade, hotness,
			 ratings_count, rating_url, hash)
		 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21) RETURNING professor_id`,
			ToNullString(p.FirstName),
			ToNullString(p.LastName),
			ToNullString(p.Department),
			ToNullString(p.Title),
			ToNullString(p.Email),
			ToNullString(p.convertPhoneNumber()),
			ToNullString(p.FaxNumber),

			ToNullString(p.Location.School),
			ToNullString(p.Location.State),
			ToNullString(p.Location.City),
			ToNullString(p.Location.Room),
			ToNullString(p.Location.Address),

			ToNullFloat64(p.Rating.Overall),
			ToNullFloat64(p.Rating.Helpfullness),
			ToNullFloat64(p.Rating.Clarity),
			ToNullFloat64(p.Rating.Easiness),
			ToNullString(p.Rating.AverageGrade),
			ToNullBool(p.Rating.Hotness),
			ToNullFloat64(p.Rating.RatingsCount),
			ToNullString(p.Rating.RatingUrl),
			p.hash()).
		Scan(&professorId)
		checkError(err)
	}
	return
}

func insertMapping(p Parameter, professorId int, db *sql.DB) (mappingId int) {
	var hash string
	db.QueryRow(
		`SELECT mapping_id, hash FROM mapping
		WHERE hash = $1`,
		p.hash(),
	).
	Scan(&mappingId, &hash)

	if hash != p.hash() {
		err := db.QueryRow(
			`INSERT INTO mapping
			(first_name, last_name, subject, course, inclusion, professor_id, hash)
		 VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING mapping_id`,
			ToNullString(p.FirstName),
			ToNullString(p.LastName),
			ToNullString(p.Department),
			ToNullString(p.CourseNumber),
			ToNullString(p.Inclusion),
			professorId,
			p.hash(),
		).
		Scan(&mappingId)
		checkError(err)
	}
	return
}

func insertExclusions(exclusions []string, db *sql.DB) (exclusionIds []int) {
	for _, val := range exclusions {
		var tempId int
		err := db.QueryRow(
			`INSERT INTO exclusions
				(url)
			 VALUES($1) RETURNING exclusion_id`,
			val).
			Scan(&tempId)
		checkError(err)
		exclusionIds = append(exclusionIds, tempId)
	}
	return
}

func insertMappingExclusions(mappingId int, exclusionIds []int, db *sql.DB) {
	for _, val := range exclusionIds {
		smt, err := db.Prepare(
			`INSERT INTO mapping_exclusions
				(exclusion_id, mapping_id)
			 VALUES($1, $2)`)
		res, err := smt.Exec(val, mappingId)
		checkError(err)

		affect, err := res.RowsAffected()
		checkError(err)

		fmt.Println(affect, "rows changed")

		checkError(err)
	}
}

func queryAdjacentMappingsByParams(params Parameter, db *sql.DB) *sql.Row {
	row := db.QueryRow(
		`SELECT professors.professor_id, professors.first_name, professors.last_name, professors.email, professors.department, professors.title,
		professors.phone_number, professors.fax_number,professors.school, professors.state, professors.city,
		professors.room, professors.address, professors.overall, professors.helpfulness, professors.clarity,
		professors.easiness, professors.average_grade,professors.hotness, professors.ratings_count, professors.rating_url
		FROM mapping
		LEFT JOIN mapping_exclusions
			ON mapping_exclusions.mapping_id = mapping.mapping_id
		LEFT JOIN professors
			ON mapping.professor_id = professors.professor_id
		WHERE mapping.professor_id IS NOT NULL
		AND mapping.first_name is NOT NULL
		AND mapping_exclusions.mapping_id IS NULL
		AND mapping.last_name = $1
		AND mapping.subject = $2
		LIMIT 1;`,
		params.LastName,
		params.Department)
	return row
}

func queryProfessorMappingById(professorId int, db *sql.DB) *sql.Row {
	row := db.QueryRow(
		`SELECT professor_id, first_name, last_name, email, department, title, phone_number, fax_number,
		school, state, city, room, address, overall, helpfulness, clarity, easiness, average_grade,
		hotness, ratings_count, rating_url
		FROM professors WHERE professor_id = $1`,
		professorId)
	return row
}

func queryProfessorMappingByParams(params Parameter, db *sql.DB) *sql.Row {
	row := db.QueryRow(
		`SELECT professors.professor_id, professors.first_name, professors.last_name, professors.email, professors.department,
		professors.title, professors.phone_number, professors.fax_number, professors.school, professors.state,
		professors.city, professors.room, professors.address, professors.overall, professors.helpfulness, professors.clarity,
		professors.easiness, professors.average_grade,professors.hotness, professors.ratings_count, professors.rating_url
		FROM mapping
		LEFT JOIN professors ON mapping.professor_id = professors.professor_id
		WHERE mapping.professor_id IS NOT NULL
		AND mapping.hash = $1
		LIMIT 1;`,
		params.hash())
	return row
}

func getProfessorFromRow(row *sql.Row) (professorId int, professor *Professor, err error) {
	//Professor

	var FirstName sql.NullString
	var LastName sql.NullString
	var Email sql.NullString
	var Department sql.NullString
	var Title sql.NullString
	var TempPhoneNumber sql.NullString
	var PhoneNumber []string
	var FaxNumber sql.NullString

	//Location
	var School sql.NullString
	var City sql.NullString
	var State sql.NullString
	var Room sql.NullString
	var Address sql.NullString

	//Rating
	var Overall sql.NullFloat64
	var Helpfulness sql.NullFloat64
	var Easiness sql.NullFloat64
	var Clarity sql.NullFloat64
	var AverageGrade sql.NullString
	var Hotness sql.NullBool
	var RatingsCount sql.NullFloat64
	var RatingUrl sql.NullString

	err = row.Scan(&professorId, &FirstName, &LastName, &Email, &Department, &Title, &TempPhoneNumber, &FaxNumber, &School,
		&State, &City, &Room, &Address, &Overall, &Helpfulness, &Clarity, &Easiness, &AverageGrade, &Hotness,
		&RatingsCount, &RatingUrl)
	PhoneNumber = strings.Split(TempPhoneNumber.String, ",")

	if err == nil {
		result := &Professor{
			FirstName:   FirstName.String,
			LastName:    LastName.String,
			Email:       Email.String,
			Department:  Department.String,
			Title:       Title.String,
			PhoneNumber: PhoneNumber,
			FaxNumber:   FaxNumber.String,
			Location: Location{
				School:  School.String,
				City:    City.String,
				State:   State.String,
				Room:    Room.String,
				Address: Address.String,
			},
			Rating: Rating{
				Overall:      Overall.Float64,
				Helpfullness: Helpfulness.Float64,
				Easiness:     Easiness.Float64,
				Clarity:      Clarity.Float64,
				AverageGrade: AverageGrade.String,
				Hotness:      Hotness.Bool,
				RatingsCount: RatingsCount.Float64,
				RatingUrl:    RatingUrl.String,
			},
		}

		//log.Printf("%#s", result)

		return professorId, result, err
	}
	return -1, nil, err
}
