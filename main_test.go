package main

import (
	"testing"
	"net/http"
)

func BenchmarkMainHandler(b *testing.B) {
	for n := 0; n < b.N; n++ {
		http.Get("http://localhost:8080/searchdb?subject=COMPUTER+SCIENCE&last=ZHANG&city=Newark")
	}
}
