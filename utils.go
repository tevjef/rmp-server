package main

import (
	"log"
	"strings"
	"database/sql"
)

const (
	Empty         = ""
	IndexNotFound = -1
)

var l = strings.ToLower
var format = formatString

func formatString(s string) string {
	return strings.Title(l(strings.TrimSpace(s)))
}
func substringAfter(str, separator string) string {
	if isEmpty(str) {
		return str
	}
	pos := strings.Index(str, separator)
	if pos == IndexNotFound {
		return Empty
	}
	return str[pos+len(separator):]
}

func substringAfterLast(str, separator string) string {
	if isEmpty(str) {
		return str
	}
	if isEmpty(separator) {
		return Empty
	}
	pos := strings.LastIndex(str, separator)
	if pos == IndexNotFound || pos == len(str)-len(separator) {
		return Empty
	}
	return str[pos+len(separator):]
}

func substringBefore(str, separator string) string {
	if isEmpty(str) {
		return str
	}
	pos := strings.Index(str, separator)
	if pos == IndexNotFound {
		return str
	}
	return str[:pos]
}

func substringBeforeLast(str, separator string) string {
	if isEmpty(str) || isEmpty(separator) {
		return str
	}
	pos := strings.LastIndex(str, separator)
	if pos == IndexNotFound {
		return Empty
	}
	return str[:pos]
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func isEmpty(str string) bool {
	return len(str) == 0
}

//ToNullString invalidates a sql.NullString if empty, validates if not empty
func ToNullString(s string) sql.NullString {
	return sql.NullString{String : s, Valid : s != ""}
}

func ToNullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64 : f, Valid : f != 0}
}

func ToNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64 : i, Valid : i != 0}
}

func ToNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool : b, Valid : true}
}


