package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPersonCreate(t *testing.T) {
	data := []byte(`{"email":"1@1.1","pw_hash":"2","name_first":"3","name_last":"4"}`)

	var p1 Person
	var p2 Person

	SimplePostRequest("/person", data, func(resp string, err error) {
		CheckErr(t, err)

		err = json.Unmarshal([]byte(resp), &p1)
		CheckErr(t, err)

		assert.Equal(t, p1.Email, "1@1.1", "they should equal")
		assert.Equal(t, p1.PW_hash, "2", "they should equal")
		assert.Equal(t, p1.Name_First, "3", "they should equal")
		assert.Equal(t, p1.Name_Last, "4", "they should equal")

	})

	SimpleGetRequest(fmt.Sprintf("/person/%d/info", p1.Person_ID),
		func(resp string, err error) {
			CheckErr(t, err)

			err = json.Unmarshal([]byte(resp), &p2)
			CheckErr(t, err)

			assert.Equal(t, p1, p2, "they should equal")
		})
}
