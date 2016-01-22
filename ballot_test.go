package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBallotCreateDestroy(t *testing.T) {
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

	// Test Creating 10 ballots for decision
	var bs []Ballot
	for i := 0; i < 10; i++ {
		data := fmt.Sprintf(`{"secret":123, "name": "ballot", "email":"123@gmail.com"}`)
		b, err := TCreateBallot(d1.Decision_ID, data)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, b.Decision_ID, d1.Decision_ID, "should be equal")
		bs = append(bs, b)
	}

	// Delete those ballots
	for _, b := range bs {
		res, err := TDeleteBallot(d1.Decision_ID, b.Ballot_ID)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, res.Result, "deleted", "should be equal")
	}
}
