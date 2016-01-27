package main

import (
	"testing"

	"github.com/verdverm/frisby"
)

func cleanDatabase(t *testing.T) {
	frisby.Create("Cleaning database").
		Get("http://localhost:9999/clean").
		Send().
		ExpectStatus(200).
		ExpectJson("result", "cleaned")
}
