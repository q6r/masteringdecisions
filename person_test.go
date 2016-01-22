package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHPersonCreateHPersonInfo
// a sample test that tests HPersonCreate,
// and HpersonInfo Handlers
func TestHPersonCreateHPersonInfo(t *testing.T) {

	dbmap = InitDatabase()
	defer dbmap.Db.Close()
	dbmap.TruncateTables()

	data := `{"email":"2@2.2","pw_hash":"2","name_first":"3","name_last":"4"}`

	p1, err := TCreatePerson(data)
	if err != nil {
		t.Error(err)
	}

	p2, err := TInfoPerson(p1.Person_ID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, p1, p2, "should be equal")
}

func TestHPersonCreateHPersonDelete(t *testing.T) {

	dbmap = InitDatabase()
	defer dbmap.Db.Close()
	dbmap.TruncateTables()

	data := `{"email":"2@2.2","pw_hash":"2","name_first":"3","name_last":"4"}`
	p1, err := TCreatePerson(data)
	if err != nil {
		t.Error(err)
	}

	res, err := TDeletePerson(p1.Person_ID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, res.Result, "deleted")
}
