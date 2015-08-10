package main

import (
	"testing"
)

func TestSubstringAfter(t *testing.T) {
	t.Log("TestSubstringAfter")

	result := substringAfter("abcba", "b")
	assertEquals(t, "cba", result)

	result = substringAfter("abc", "a")
	assertEquals(t, "bc", result)

	result = substringAfter("abc", "c")
	assertEquals(t, "", result)
}

func TestSubstringAfterLast(t *testing.T) {
	t.Log("TestSubstringAfterLast")

	result := substringAfterLast("abcba", "b")
	assertEquals(t, "a", result)

	result = substringAfterLast("abc", "c")
	assertEquals(t, "", result)

	result = substringAfterLast("a", "a")
	assertEquals(t, "", result)
}

func TestSubstringBefore(t *testing.T) {
	t.Log("TestSubstringBefore")

	result := substringBefore("abcba", "b")
	assertEquals(t, "a", result)

	result = substringBefore("abc", "a")
	assertEquals(t, "", result)

	result = substringBefore("abc", "c")
	assertEquals(t, "ab", result)
}

func TestSubstringBeforeLast(t *testing.T) {
	t.Log("TestSubstringBeforeLast")

	result := substringBeforeLast("abcba", "b")
	assertEquals(t, "abc", result)

	result = substringBeforeLast("abc", "a")
	assertEquals(t, "", result)

	result = substringBeforeLast("abc", "")
	assertEquals(t, "abc", result)
}

func assertEquals(t *testing.T, expected, actual string) {
	t.Logf("%s (expected)\n"+"%s (actual)", expected, actual)
	if expected != actual {
		t.Errorf("Not equal: %#v (expected)\n"+
			"        != %#v (actual)", expected, actual)
	}

}
