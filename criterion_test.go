package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCriterionCreateDelete(t *testing.T) {
	dbmap = InitDatabase()
	defer dbmap.Db.Close()
	dbmap.TruncateTables()

	// need to have person and decision first
	p1, err := TCreatePerson(`{"email":"2@2.2","pw_hash":"2","name_first":"3","name_last":"4"}`)
	if err != nil {
		t.Error(err)
	}
	data := fmt.Sprintf(`{"person_id":%d, "name":"hello", "description":"abc", "stage":123, "criterion_vote_style":"def", "alternative_vote_style":"def", "client_settings":"def"}`, p1.Person_ID)
	d1, err := TCreateDecision(data)
	if err != nil {
		t.Error(err)
	}

	c1, err := TCreateCriterion(d1.Decision_ID,
		`{"name": "crittst", "weight":100}`)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, c1.Decision_ID, d1.Decision_ID, "should be equal")

	// Delete the criterion created
	jres, err := TDeleteCriterion(d1.Decision_ID, c1.Criterion_ID)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, jres.Result, "deleted", "should be deleted")
}
