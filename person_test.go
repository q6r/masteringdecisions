package main

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHPersonCreateHPersonInfo(t *testing.T) {
	dbmap = InitDatabase()
	defer dbmap.Db.Close()
	dbmap.TruncateTables()

	data := `{"email":"2@2.2","pw_hash":"2","name_first":"3","name_last":"4"}`
	var p1, p2 Person

	RunSimplePost("/person", data,
		func(c *gin.Context) {
			HPersonCreate(c)
		},
		func(r *httptest.ResponseRecorder) {
			err := json.Unmarshal(r.Body.Bytes(), &p1)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, p1.Email, "2@2.2", "Should be equal")
			assert.Equal(t, p1.PW_hash, "2", "Should be equal")
			assert.Equal(t, p1.Name_First, "3", "Should be equal")
			assert.Equal(t, p1.Name_Last, "4", "Should be equal")
		})

	RunSimpleGet("/person/<something>/info",
		func(c *gin.Context) {
			c.Params = gin.Params{gin.Param{Key: "person_id", Value: strconv.Itoa(p1.Person_ID)}}
			HPersonInfo(c)
		},
		func(r *httptest.ResponseRecorder) {
			err := json.Unmarshal(r.Body.Bytes(), &p2)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, p1, p2, "Should be equal")
		})
}

type ErrorJson struct {
	Result string `json:"result"`
}

func TestHPersonCreateHPersonDelete(t *testing.T) {
	dbmap = InitDatabase()
	defer dbmap.Db.Close()
	dbmap.TruncateTables()

	data := `{"email":"2@2.2","pw_hash":"2","name_first":"3","name_last":"4"}`
	var p1 Person

	RunSimplePost("/person", data,
		func(c *gin.Context) {
			HPersonCreate(c)
		},
		func(r *httptest.ResponseRecorder) {
			err := json.Unmarshal(r.Body.Bytes(), &p1)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, p1.Email, "2@2.2", "Should be equal")
			assert.Equal(t, p1.PW_hash, "2", "Should be equal")
			assert.Equal(t, p1.Name_First, "3", "Should be equal")
			assert.Equal(t, p1.Name_Last, "4", "Should be equal")
		})

	RunSimpleDelete("/person/<something>",
		func(c *gin.Context) {
			c.Params = gin.Params{gin.Param{Key: "person_id", Value: strconv.Itoa(p1.Person_ID)}}
			HPersonDelete(c)
		},
		func(r *httptest.ResponseRecorder) {
			fmt.Printf(string(r.Body.Bytes()))
			var jerr ErrorJson
			err := json.Unmarshal(r.Body.Bytes(), &jerr)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, jerr.Result, "deleted", "Should be equal")
		})
}
