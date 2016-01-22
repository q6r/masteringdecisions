package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecisionCreate(t *testing.T) {
	dbmap = InitDatabase()
	defer dbmap.Db.Close()
	dbmap.TruncateTables()

	p1, err := TCreatePerson(`{"email":"2@2.2","pw_hash":"2","name_first":"3","name_last":"4"}`)
	if err != nil {
		t.Error(err)
	}

	var ds []Decision
	for i := 0; i < 10; i++ {
		data := fmt.Sprintf(`{"person_id":%d, "name":"hello", "description":"abc", "stage":123, "criterion_vote_style":"def", "alternative_vote_style":"def", "client_settings":"def"}`, p1.Person_ID)
		d, err := TCreateDecision(data)
		if err != nil {
			t.Error(err)
		}
		ds = append(ds, d)
	}

	for _, d := range ds {
		assert.Equal(t, d.Person_ID, p1.Person_ID, "Owner should be same")
	}

	// Delete decisions
	for _, d := range ds {
		res, err := TDeleteDecision(d.Decision_ID)
		assert.Equal(t, err, nil, "Failed while destroying")
		assert.Equal(t, res.Result, "deleted decision, its ballots")
	}
}
