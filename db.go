package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
	"errors"
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

func RefreshDatabase(db *sql.DB) {
	params := constructParametersFromRows(queryStaleMappingsForUpdate(db), db)

	for _, param := range params{
		professor := SearchRMP(param)

		if professor != nil {
			professorId := insertProfessor(professor, db)
			insertOrUpdateMapping(param, professorId, db)
		}

		fmt.Printf("REFRESHED: %#s", professor)
	}
}

func Search(params Parameter, db *sql.DB) (professor *Professor) {
	fmt.Printf("Search() %#v Hash: %s \n", params, params.hash())
	//First search the database for the professor
	professor, _ = SearchDatabase(params, db)

	if mappingExists := checkMappingExists(params, db); professor == nil && !mappingExists {
		//If they're not in the database, scrape RMP
		professor = SearchRMP(params)
		var professorId int64
		if professor != nil {
			professorId = insertProfessor(professor, db)
		}
		mappingId := insertOrUpdateMapping(params, professorId, db)
		exclusionIds := insertExclusions(params.Exclusion, db)
		insertMappingExclusions(mappingId, exclusionIds, db)
		_, professor, _ = getProfessorFromRow(queryProfessorMappingById(professorId, db))
	}
	return
}

func SearchDatabase(params Parameter, db *sql.DB) (professor *Professor, err error) {
	var professorId int64

	if professor == nil {
		professorId, professor, err = getProfessorFromRow(queryProfessorMappingByParams(params, db))
		fmt.Printf("ID: %dRESULT SEARCH DIRECT: %#v\n\n",professorId, professor)
	}

	if professor == nil {
		professorId, professor, err = getProfessorFromRow(queryAdjacentMappingsByParams(params, db))
		if err != nil && professorId != -1 {
			insertOrUpdateMapping(params, professorId, db)
			fmt.Printf("ID: %d INSERTING ADJACENT: %#v\n\n",professorId, professor)

		}
		fmt.Printf("ID: %d SEARCH ADJACENT: %#s\n\n",professorId, professor)
	}
	if err != nil && professor == nil {
		err = errors.New("No professor found")
	}
	return professor, nil
}

func insertProfessor(p *Professor, db *sql.DB) (professorId int64) {
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
			ToNullFloat64(p.Rating.Helpfulness),
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

func insertOrUpdateMapping(p Parameter, professorId int64, db *sql.DB) (mappingId int64) {
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
			(first_name, last_name, subject, course, inclusion, professor_id, hash, is_rutgers, city)
		 VALUES($1,$2,$3,$4,$5,$6,$7,$8, $9) RETURNING mapping_id`,
			ToNullString(p.FirstName),
			ToNullString(p.LastName),
			ToNullString(p.Department),
			ToNullString(p.CourseNumber),
			ToNullString(p.Inclusion),
			ToNullInt64(professorId),
			p.hash(),
			ToNullBool(p.IsRutgers),
			ToNullString(p.City),
		).
		Scan(&mappingId)
		checkError(err)
		fmt.Printf("Inserting mapping: %#s", p)

	} else {
		updateMapping(mappingId, professorId, db)
	}
	return
}

func insertExclusions(exclusions []string, db *sql.DB) (exclusionIds []int64) {
	for _, val := range exclusions {
		var tempId int64
		var tempUrl sql.NullString
		db.QueryRow(
		`SELECT url FROM exclusions
				WHERE url = $1 LIMIT 1 RETURNING exclusion_id, url`,
		val).
		Scan(&tempId, &tempUrl)

		if tempUrl.String != val {
			err := db.QueryRow(
				`INSERT INTO exclusions
				(url)
			 VALUES($1) RETURNING exclusion_id`,
				val).
			Scan(&tempId)
			checkError(err)
		}
		exclusionIds = append(exclusionIds, tempId)
	}
	return
}

func insertMappingExclusions(mappingId int64, exclusionIds []int64, db *sql.DB) {
	for _, val := range exclusionIds {

		db.QueryRow(`INSERT INTO mapping_exclusions
				(exclusion_id, mapping_id)
			 VALUES($1, $2)`, val, mappingId)

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

func constructParametersFromRows(rows *sql.Rows, db *sql.DB) (params []Parameter) {
	for rows.Next() {
		var mappingId int64
		var firstName sql.NullString
		var lastName sql.NullString
		var city sql.NullString
		var subject sql.NullString
		var course sql.NullString
		var inclusion sql.NullString
		var exclusions []string
		var isRutgers bool

		rows.Scan(&mappingId, &firstName, &lastName, &subject, &course, &inclusion, &isRutgers, &city)

		exclRows := queryExclusionsForMapping(mappingId, db)
		for exclRows.Next() {
			var tempExclusions sql.NullString
			exclRows.Scan(&tempExclusions)
			exclusions = append(exclusions, tempExclusions.String)
		}

		param :=  Parameter{
			FirstName:firstName.String,
			LastName:lastName.String,
			Department:subject.String,
			City:city.String,
			CourseNumber:course.String,
			Inclusion:inclusion.String,
			Exclusion:exclusions,
			IsRutgers:isRutgers,
		}
		params = append(params, param)
	}
	return
}

func queryStaleMappingsForUpdate(db *sql.DB) *sql.Rows {
	rows, err := db.Query(
		`SELECT mapping.mapping_id, mapping.first_name, mapping.last_name, mapping.subject, mapping.course, mapping.inclusion, is_rutgers, city
		FROM
			mapping
		WHERE
			mapping.is_stale`)
	checkError(err)
	return rows
}

func updateMapping(mappingId, professorId int64, db *sql.DB) (id int64) {
	row := db.QueryRow(
		`UPDATE mapping
		SET professor_id = $1, is_stale = FALSE
		WHERE mapping_id = $2 RETURNING mapping_id`,
		ToNullInt64(professorId),
		ToNullInt64(mappingId))
	err := row.Scan(&id)
	checkError(err)
	fmt.Printf("Updating mapping: %d", mappingId)
	return
}

//Gets all exclusions for a particular mapping
func queryExclusionsForMapping(mappingId int64, db *sql.DB) *sql.Rows {
	rows, err := db.Query(
		`SELECT exclusions.url FROM mapping_exclusions
		INNER JOIN
			exclusions ON mapping_exclusions.exclusion_id = exclusions.exclusion_id
		INNER JOIN
			mapping ON mapping.mapping_id = mapping_exclusions.mapping_id
		WHERE
			mapping.mapping_id = $1`, ToNullInt64(mappingId))
	checkError(err)
	return rows
}

func queryAdjacentMappingsByParams(params Parameter, db *sql.DB) *sql.Row {
	fmt.Printf("queryAdjacentMappingsByParams() Hash: %s \n", params.hash())
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
		AND mapping_exclusions.mapping_id IS NULL
		AND mapping.first_name = $1
		AND mapping.last_name = $2
		AND mapping.subject = $3
		LIMIT 1;`,
		params.FirstName,
		params.LastName,
		params.Department)
	return row
}

func queryProfessorMappingById(professorId int64, db *sql.DB) *sql.Row {
	fmt.Printf("ID: %d queryProfessorMappingById()\n", professorId)
	row := db.QueryRow(
		`SELECT professor_id, first_name, last_name, email, department, title, phone_number, fax_number,
		school, state, city, room, address, overall, helpfulness, clarity, easiness, average_grade,
		hotness, ratings_count, rating_url
		FROM professors WHERE professor_id = $1`,
		ToNullInt64(professorId))
	return row
}

func queryProfessorMappingByParams(params Parameter, db *sql.DB) *sql.Row {
	fmt.Printf("queryProfessorMappingByParams() Hash: %s \n", params.hash())
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

func checkMappingExists(params Parameter, db *sql.DB) (exists bool) {
	db.QueryRow(
		`SELECT EXISTS(SELECT * FROM mapping WHERE hash = $1) AS bool`,
		params.hash()).Scan(&exists)
	fmt.Printf("checkMappingExists() %#v Hash: %s Result: %s \n", params, params.hash(), exists)
	return
}


func getProfessorFromRow(row *sql.Row) (professorId int64, professor *Professor, err error) {
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
		professor := &Professor{
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
				Helpfulness: Helpfulness.Float64,
				Easiness:     Easiness.Float64,
				Clarity:      Clarity.Float64,
				AverageGrade: AverageGrade.String,
				Hotness:      Hotness.Bool,
				RatingsCount: RatingsCount.Float64,
				RatingUrl:    RatingUrl.String,
			},
		}

		fmt.Printf("getProfessorFromRow() %#v\n", professor)

		return professorId, professor, err
	}
	return -1, nil, err
}

func incrementStaleCount(param Parameter, db *sql.DB) (mappingId, count int64) {
	row := db.QueryRow(
		`UPDATE mapping
		SET stale_count = stale_count + 1
		WHERE hash = $1 RETURNING mapping_id, stale_count`, param.hash())
	row.Scan(&mappingId, &count)
	return
}