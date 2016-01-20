package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func CheckErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func SimplePostRequest(url string,
	data []byte, checker func(resp string, err error)) {

	req, err := http.NewRequest("POST", "http://localhost:9999"+url, bytes.NewBuffer(data))
	if err != nil {
		checker("", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		checker("", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		checker("", fmt.Errorf("Error status code %v", resp.StatusCode))
	}

	iob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		checker("", err)
	}

	checker(string(iob), nil)
}

func SimpleGetRequest(url string, checker func(resp string, err error)) {

	req, err := http.NewRequest("GET", "http://localhost:9999"+url, nil)
	if err != nil {
		checker("", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		checker("", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		checker("", fmt.Errorf("Error status code %v", resp.StatusCode))
	}

	iob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		checker("", err)
	}

	checker(string(iob), nil)
}
