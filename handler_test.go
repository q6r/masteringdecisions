// +build testrunmain

package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/verdverm/frisby"
)

func cleanDatabase(t *testing.T) {
	frisby.Create("Cleaning database").
		Get("http://localhost:9999/clean").
		Send().
		ExpectStatus(200).
		ExpectJson("result", "cleaned")
}

func TestRunMain(t *testing.T) {
	go main()
	time.Sleep(4 * time.Second)
}

/* PERSON
 */

func TestHPersonCreate(t *testing.T) {

	cleanDatabase(t)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200)

	frisby.Create("Test HPersonCreate (fails)").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(403)

	frisby.Global.PrintReport()
}

func TestHPersonsList(t *testing.T) {

	cleanDatabase(t)

	for i := 0; i < 10; i++ {
		frisby.Create("Test HPersonCreate").
			Post("http://localhost:9999/person").
			SetHeader("Content-Type", "application/json").
			SetJson(Person{NameFirst: "a", NameLast: "b", Email: fmt.Sprintf("abcd%d@abcd.com", i), PWHash: "c"}).
			Send().
			ExpectStatus(200)

		frisby.Create("Test HPersonsList").
			Get("http://localhost:9999/persons").
			SetHeader("Content-Type", "application/json").
			Send().
			ExpectStatus(200).
			ExpectJsonLength("persons", i+1)
	}

	frisby.Global.PrintReport()
}

func TestHPersonUpdate(t *testing.T) {
	cleanDatabase(t)

	frisby.Create("Create a person").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		frisby.Create("Update this person").
			Put(fmt.Sprintf("http://localhost:9999/person/%d", pid)).
			SetHeader("Content-Type", "application/json").
			SetJson(Person{NameFirst: "newa", NameLast: "newb", Email: "newe@abcd.com", PWHash: "pwd"}).
			Send().
			ExpectStatus(200).
			ExpectJson("person.email", "newe@abcd.com").
			ExpectJson("person.name_first", "newa").
			ExpectJson("person.name_last", "newb")
	})

	frisby.Global.PrintReport()

}

func TestHPersonInfo(t *testing.T) {
	cleanDatabase(t)

	frisby.Create("Create a person").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}
		frisby.Create("Check the person information").
			Get(fmt.Sprintf("http://localhost:9999/person/%d/info", pid)).
			SetHeader("Content-Type", "application/json").
			Send().
			ExpectStatus(200)
	})

	frisby.Global.PrintReport()
}

func TestHPersonDecisions(t *testing.T) {
	cleanDatabase(t)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		frisby.Create("Test HPersonDecisions").
			Get(fmt.Sprintf("http://localhost:9999/person/%d/decisions", pid)).
			SetHeader("Content-Type", "application/json").
			Send().
			ExpectJsonLength("decisions", 0)

		frisby.Create("Test HDecisionCreate").
			Post("http://localhost:9999/decision").
			SetHeader("Content-Type", "application/json").
			SetJson(Decision{PersonID: pid, Name: "t1", Description: "desc",
			Stage: 1, AlternativeVoteStyle: "alt", CriterionVoteStyle: "crit", ClientSettings: "clis"}).
			Send().
			ExpectStatus(200)

		frisby.Create("Test HPersonDecisions").
			Get(fmt.Sprintf("http://localhost:9999/person/%d/decisions", pid)).
			SetHeader("Content-Type", "application/json").
			Send().
			ExpectJsonLength("decisions", 1)
	})

	frisby.Global.PrintReport()
}

func TestHPersonDelete(t *testing.T) {

	cleanDatabase(t)

	frisby.Create("Test HPersonDelete (fails)").
		Delete("http://localhost:9999/person/1").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(403)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		frisby.Create("Test HPersonDelete").
			Delete(fmt.Sprintf("http://localhost:9999/person/%d", pid)).
			SetHeader("Content-Type", "application/json").
			SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
			Send().
			ExpectStatus(200)

	})

	frisby.Global.PrintReport()
}

/* DECISION
 */

func TestHDecisionCreate(t *testing.T) {

	cleanDatabase(t)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		frisby.Create("Test HDecisionCreate").
			Post("http://localhost:9999/decision").
			SetHeader("Content-Type", "application/json").
			SetJson(Decision{PersonID: pid, Name: "t1", Description: "desc",
			Stage: 1, AlternativeVoteStyle: "alt", CriterionVoteStyle: "crit", ClientSettings: "clis"}).
			Send().
			ExpectStatus(200)
	})

	frisby.Global.PrintReport()
}

func TestHDecisionInfo(t *testing.T) {

	cleanDatabase(t)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		frisby.Create("Test HDecisionCreate").
			Post("http://localhost:9999/decision").
			SetHeader("Content-Type", "application/json").
			SetJson(Decision{PersonID: pid, Name: "t1", Description: "desc",
			Stage: 1, AlternativeVoteStyle: "alt", CriterionVoteStyle: "crit", ClientSettings: "clis"}).
			Send().
			ExpectStatus(200).
			AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

			did, err := json.Get("decision").Get("decision_id").Int()
			if err != nil {
				t.Error(err)
			}

			frisby.Create("Test HDecisionInfo").
				Get(fmt.Sprintf("http://localhost:9999/decision/%d/info", did)).
				SetHeader("Content-Type", "application/json").
				Send().
				ExpectStatus(200).
				ExpectJson("decision.name", "t1").
				ExpectJson("decision.description", "desc")

		})
	})

	frisby.Global.PrintReport()
}

func TestHDecisionCriterionsList(t *testing.T) {

	cleanDatabase(t)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		frisby.Create("Test HDecisionCreate").
			Post("http://localhost:9999/decision").
			SetHeader("Content-Type", "application/json").
			SetJson(Decision{PersonID: pid, Name: "t1", Description: "desc",
			Stage: 1, AlternativeVoteStyle: "alt", CriterionVoteStyle: "crit", ClientSettings: "clis"}).
			Send().
			ExpectStatus(200).
			AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

			did, err := json.Get("decision").Get("decision_id").Int()
			if err != nil {
				t.Error(err)
			}

			frisby.Create("Test HBallotCreate").
				Post(fmt.Sprintf("http://localhost:9999/decision/%d/ballot", did)).
				SetHeader("Content-Type", "application/json").
				SetJson(Ballot{Name: "b1", Email: "email"}).
				Send().
				ExpectStatus(200).
				ExpectJson("ballot.name", "b1").
				ExpectJson("ballot.email", "email").
				AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

				bid, err := json.Get("ballot").Get("ballot_id").Int()
				if err != nil {
					t.Error(err)
				}

				frisby.Create("Test HDecisionBallotsList").
					Get(fmt.Sprintf("http://localhost:9999/decision/%d/ballots", did)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(200).
					ExpectJsonLength("ballots", 1)

				frisby.Create("Test HBallotUpdate").
					Put(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d", did, bid)).
					SetHeader("Content-Type", "application/json").
					SetJson(Ballot{Name: "b2", Email: "email2"}).
					Send().
					ExpectStatus(200).
					ExpectJson("ballot.name", "b2").
					ExpectJson("ballot.email", "email2")

				frisby.Create("Test HBallotInfo").
					Get(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d/info", did, bid)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(200).
					ExpectJson("ballot.name", "b2").
					ExpectJson("ballot.email", "email2")

				frisby.Create("Test HBallotAllInfo").
					Get(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d", did, bid)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(200).
					ExpectJson("ballot.name", "b2").
					ExpectJson("ballot.email", "email2")

				frisby.Create("Test HBallotDelete").
					Delete(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d", did, bid)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(200).
					ExpectJson("result", "deleted")

				frisby.Create("Test HBallotInfo").
					Get(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d/info", did, bid)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(403)

			})

		})
	})

	frisby.Global.PrintReport()
}

func TestHDecisionList(t *testing.T) {

	cleanDatabase(t)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		for i := 0; i < 5; i++ {
			frisby.Create("Test HDecisionCreate").
				Post("http://localhost:9999/decision").
				SetHeader("Content-Type", "application/json").
				SetJson(Decision{PersonID: pid, Name: "t1", Description: "desc",
				Stage: 1, AlternativeVoteStyle: "alt", CriterionVoteStyle: "crit", ClientSettings: "clis"}).
				Send().
				ExpectStatus(200)
		}

		frisby.Create("Test HDecisionsList").
			Get(fmt.Sprintf("http://localhost:9999/decisions")).
			SetHeader("Content-Type", "application/json").
			Send().
			ExpectStatus(200)

	})

	frisby.Global.PrintReport()
}

func TestHDecisionBallotsList(t *testing.T) {

	cleanDatabase(t)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		frisby.Create("Test HDecisionCreate").
			Post("http://localhost:9999/decision").
			SetHeader("Content-Type", "application/json").
			SetJson(Decision{PersonID: pid, Name: "t1", Description: "desc",
			Stage: 1, AlternativeVoteStyle: "alt", CriterionVoteStyle: "crit", ClientSettings: "clis"}).
			Send().
			ExpectStatus(200).
			AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

			did, err := json.Get("decision").Get("decision_id").Int()
			if err != nil {
				t.Error(err)
			}

			frisby.Create("Test HCriterionCreate").
				Post(fmt.Sprintf("http://localhost:9999/decision/%d/criterion", did)).
				SetHeader("Content-Type", "application/json").
				SetJson(Criterion{Name: "c1"}).
				Send().
				ExpectStatus(200).
				ExpectJson("criterion.name", "c1").
				AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
				cid, err := json.Get("criterion").Get("criterion_id").Int()
				if err != nil {
					t.Error(err)
				}

				frisby.Create("Test HDecisionCriterionsList").
					Get(fmt.Sprintf("http://localhost:9999/decision/%d/criterions", did)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(200).
					ExpectJsonLength("criterions", 1)

				frisby.Create("Test HCriterionInfo").
					Get(fmt.Sprintf("http://localhost:9999/decision/%d/criterion/%d/info", did, cid)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(200).
					ExpectJson("criterion.name", "c1")

				frisby.Create("Test HCriterionUpdate").
					Put(fmt.Sprintf("http://localhost:9999/decision/%d/criterion/%d", did, cid)).
					SetHeader("Content-Type", "application/json").
					SetJson(Criterion{Name: "c2"}).
					Send().
					ExpectStatus(200).
					ExpectJson("criterion.name", "c2")

				frisby.Create("Test HCriterionDelete").
					Delete(fmt.Sprintf("http://localhost:9999/decision/%d/criterion/%d", did, cid)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(200).
					ExpectJson("result", "deleted")

				frisby.Create("Test HCriterionInfo").
					Get(fmt.Sprintf("http://localhost:9999/decision/%d/criterion/%d/info", did, cid)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(403)

			})

		})
	})

	frisby.Global.PrintReport()
}

func TestHDecisionDelete(t *testing.T) {

	cleanDatabase(t)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		for i := 0; i < 5; i++ {
			frisby.Create("Test HDecisionCreate").
				Post("http://localhost:9999/decision").
				SetHeader("Content-Type", "application/json").
				SetJson(Decision{PersonID: pid, Name: "t1", Description: "desc",
				Stage: 1, AlternativeVoteStyle: "alt", CriterionVoteStyle: "crit", ClientSettings: "clis"}).
				Send().
				ExpectStatus(200).
				AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

				did, err := json.Get("decision").Get("decision_id").Int()
				if err != nil {
					t.Error(err)
				}

				frisby.Create("Test HBallotCreate").
					Post(fmt.Sprintf("http://localhost:9999/decision/%d/ballot", did)).
					SetHeader("Content-Type", "application/json").
					SetJson(Ballot{Name: "b1", Email: "email"}).
					Send().
					ExpectStatus(200).
					ExpectJson("ballot.name", "b1").
					ExpectJson("ballot.email", "email")

				frisby.Create("Test HDecisionDelete").
					Delete(fmt.Sprintf("http://localhost:9999/decision/%d", did)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(200).
					ExpectJson("result", "deleted")

			})

		}

		frisby.Create("Test HDecisionsList").
			Get(fmt.Sprintf("http://localhost:9999/decisions")).
			SetHeader("Content-Type", "application/json").
			Send().
			ExpectStatus(200)

	})

	frisby.Global.PrintReport()
}

func TestHDecisionUpdate(t *testing.T) {

	cleanDatabase(t)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		frisby.Create("Test HDecisionCreate").
			Post("http://localhost:9999/decision").
			SetHeader("Content-Type", "application/json").
			SetJson(Decision{PersonID: pid, Name: "t1", Description: "desc",
			Stage: 1, AlternativeVoteStyle: "alt", CriterionVoteStyle: "crit", ClientSettings: "clis"}).
			Send().
			ExpectStatus(200).
			AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

			did, err := json.Get("decision").Get("decision_id").Int()
			if err != nil {
				t.Error(err)
			}

			frisby.Create("Test HDecisionUpdate").
				Put(fmt.Sprintf("http://localhost:9999/decision/%d", did)).
				SetHeader("Content-Type", "application/json").
				SetJson(Decision{PersonID: pid, Name: "t2", Description: "desc2",
				Stage: 1, AlternativeVoteStyle: "alt", CriterionVoteStyle: "crit", ClientSettings: "clis"}).
				Send().
				ExpectStatus(200).
				ExpectJson("decision.name", "t2").
				ExpectJson("decision.description", "desc2").
				ExpectJson("decision.stage", 1)

		})

	})

	frisby.Global.PrintReport()
}

func TestHRatings(t *testing.T) {

	cleanDatabase(t)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		frisby.Create("Test HDecisionCreate").
			Post("http://localhost:9999/decision").
			SetHeader("Content-Type", "application/json").
			SetJson(Decision{PersonID: pid, Name: "t1", Description: "desc",
			Stage: 1, AlternativeVoteStyle: "alt", CriterionVoteStyle: "crit", ClientSettings: "clis"}).
			Send().
			ExpectStatus(200).
			AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

			did, err := json.Get("decision").Get("decision_id").Int()
			if err != nil {
				t.Error(err)
			}

			frisby.Create("Test HCriterionCreate").
				Post(fmt.Sprintf("http://localhost:9999/decision/%d/criterion", did)).
				SetHeader("Content-Type", "application/json").
				SetJson(Criterion{Name: "c1"}).
				Send().
				ExpectStatus(200).
				ExpectJson("criterion.name", "c1").
				AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

				cid, err := json.Get("criterion").Get("criterion_id").Int()
				if err != nil {
					t.Error(err)
				}

				frisby.Create("Test HBallotCreate").
					Post(fmt.Sprintf("http://localhost:9999/decision/%d/ballot", did)).
					SetHeader("Content-Type", "application/json").
					SetJson(Ballot{Name: "b1", Email: "email"}).
					Send().
					ExpectStatus(200).
					ExpectJson("ballot.name", "b1").
					ExpectJson("ballot.email", "email").
					AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

					bid, err := json.Get("ballot").Get("ballot_id").Int()
					if err != nil {
						t.Error(err)
					}

					frisby.Create("Test HRatingCreate").
						Get(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d/criterion/%d/vote/20", did, bid, cid)).
						SetHeader("Content-Type", "application/json").
						Send().
						ExpectStatus(200).
						ExpectJson("rating.ballot_id", bid).
						ExpectJson("rating.criterion_id", cid).
						ExpectJson("rating.rating", 20)

					frisby.Create("Test HRatingUpdate").
						Put(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d/criterion/%d/vote/30", did, bid, cid)).
						SetHeader("Content-Type", "application/json").
						Send().
						ExpectStatus(200).
						ExpectJson("rating.rating", 30).
						ExpectJson("rating.ballot_id", bid).
						ExpectJson("rating.criterion_id", cid)

					frisby.Create("Test HRatingDelete").
						Delete(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d/criterion/%d/vote", did, bid, cid)).
						SetHeader("Content-Type", "application/json").
						Send().
						ExpectStatus(200).
						ExpectJson("result", "deleted")

					frisby.Create("Test HRatingCreate").
						Get(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d/criterion/%d/vote/20", did, bid, cid)).
						SetHeader("Content-Type", "application/json").
						Send().
						ExpectStatus(200).
						ExpectJson("rating.ballot_id", bid).
						ExpectJson("rating.criterion_id", cid).
						ExpectJson("rating.rating", 20)

					frisby.Create("Test HRatings").
						Get(fmt.Sprintf("http://localhost:9999/decision/%d/criterion/%d/votes", did, cid)).
						SetHeader("Content-Type", "application/json").
						Send().
						ExpectStatus(200).
						ExpectJsonLength("ratings", 1)

					// Test voting
					frisby.Create("Test HAlternativeCreate(").
						Post(fmt.Sprintf("http://localhost:9999/decision/%d/alternative", did)).
						SetHeader("Content-Type", "application/json").
						SetJson(Alternative{Name: "a1", Description: "a1", Cost: 123.123}).
						Send().
						ExpectStatus(200).
						ExpectJson("alternative.name", "a1").
						ExpectJson("alternative.description", "a1").
						ExpectJson("alternative.cost", 123.123).
						AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

						aid, err := json.Get("alternative").Get("alternative_id").Int()
						if err != nil {
							t.Error(err)
						}

						// now do voting on pairs
						frisby.Create("Test HVotingCreate").
							Get(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d/alternative/%d/criterion/%d/vote/20", did, bid, aid, cid)).
							SetHeader("Content-Type", "application/json").
							Send().
							ExpectStatus(200).
							ExpectJson("vote.alternative_id", aid).
							ExpectJson("vote.criterion_id", cid).
							ExpectJson("vote.ballot_id", bid).
							ExpectJson("vote.weight", 20)

						frisby.Create("Test HVoteUpdate").
							Put(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d/alternative/%d/criterion/%d/vote/30", did, bid, aid, cid)).
							SetHeader("Content-Type", "application/json").
							Send().
							ExpectStatus(200).
							ExpectJson("vote.weight", 30).
							ExpectJson("vote.ballot_id", bid).
							ExpectJson("vote.criterion_id", cid).
							ExpectJson("vote.alternative_id", aid)

						frisby.Create("Test HVoteDelete").
							Delete(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d/alternative/%d/criterion/%d/vote", did, bid, aid, cid)).
							SetHeader("Content-Type", "application/json").
							Send().
							ExpectStatus(200).
							ExpectJson("result", "deleted")

						frisby.Create("Test HVotingCreate").
							Get(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d/alternative/%d/criterion/%d/vote/20", did, bid, aid, cid)).
							SetHeader("Content-Type", "application/json").
							Send().
							ExpectStatus(200)

						frisby.Create("Test HVotes").
							Get(fmt.Sprintf("http://localhost:9999/decision/%d/ballot/%d/votes", did, bid)).
							SetHeader("Content-Type", "application/json").
							Send().
							ExpectStatus(200).
							ExpectJsonLength("votes", 1)

					})

				})

			})

		})

	})

	frisby.Global.PrintReport()
}

func TestHAlternativeStar(t *testing.T) {

	cleanDatabase(t)

	frisby.Create("Test HPersonCreate").
		Post("http://localhost:9999/person").
		SetHeader("Content-Type", "application/json").
		SetJson(Person{NameFirst: "a", NameLast: "b", Email: "abcd@abcd.com", PWHash: "c"}).
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

		pid, err := json.Get("person").Get("person_id").Int()
		if err != nil {
			t.Error(err)
		}

		frisby.Create("Test HDecisionCreate").
			Post("http://localhost:9999/decision").
			SetHeader("Content-Type", "application/json").
			SetJson(Decision{PersonID: pid, Name: "t1", Description: "desc",
			Stage: 1, AlternativeVoteStyle: "alt", CriterionVoteStyle: "crit", ClientSettings: "clis"}).
			Send().
			ExpectStatus(200).
			AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

			did, err := json.Get("decision").Get("decision_id").Int()
			if err != nil {
				t.Error(err)
			}

			frisby.Create("Test HAlternativeCreate").
				Post(fmt.Sprintf("http://localhost:9999/decision/%d/alternative", did)).
				SetHeader("Content-Type", "application/json").
				SetJson(Alternative{Name: "alt", Description: "desc1"}).
				Send().
				ExpectStatus(200).
				ExpectJson("alternative.name", "alt").
				ExpectJson("alternative.description", "desc1")

			frisby.Create("Test HAlternativeCreate").
				Post(fmt.Sprintf("http://localhost:9999/decision/%d/alternative", did)).
				SetHeader("Content-Type", "application/json").
				SetJson(Alternative{Name: "alt", Description: "desc1"}).
				Send().
				ExpectStatus(200).
				ExpectJson("alternative.name", "alt").
				ExpectJson("alternative.description", "desc1")

			frisby.Create("Test HAlternativeCreate").
				Post(fmt.Sprintf("http://localhost:9999/decision/%d/alternative", did)).
				SetHeader("Content-Type", "application/json").
				SetJson(Alternative{Name: "alt", Description: "desc1"}).
				Send().
				ExpectStatus(200).
				ExpectJson("alternative.name", "alt").
				ExpectJson("alternative.description", "desc1").
				AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {

				aid, err := json.Get("alternative").Get("alternative_id").Int()
				if err != nil {
					t.Error(err)
				}

				frisby.Create("Test HAlternativeUpdate").
					Put(fmt.Sprintf("http://localhost:9999/decision/%d/alternative/%d", did, aid)).
					SetHeader("Content-Type", "application/json").
					SetJson(Alternative{Name: "up", Description: "up2"}).
					Send().
					ExpectStatus(200).
					ExpectJson("alternative.name", "up").
					ExpectJson("alternative.description", "up2")

				frisby.Create("Test HAlternativeInfo").
					Get(fmt.Sprintf("http://localhost:9999/decision/%d/alternative/%d/info", did, aid)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(200).
					ExpectJson("alternative.name", "up").
					ExpectJson("alternative.description", "up2")

				frisby.Create("Test HAlternativeDelete").
					Delete(fmt.Sprintf("http://localhost:9999/decision/%d/alternative/%d", did, aid)).
					SetHeader("Content-Type", "application/json").
					Send().
					ExpectStatus(200).
					ExpectJson("result", "deleted")

			})

			frisby.Create("Test HDecisionAlternativesList").
				Get(fmt.Sprintf("http://localhost:9999/decision/%d/alternatives", did)).
				SetHeader("Content-Type", "application/json").
				Send().
				ExpectStatus(200).
				ExpectJsonLength("alternatives", 2)

		})
	})

	frisby.Global.PrintReport()
}

func TestHRoots(t *testing.T) {
	cleanDatabase(t)

	frisby.Create("Test HRoot").
		Get(fmt.Sprintf("http://localhost:9999/")).
		SetHeader("Accept", "text/html").
		Send().
		ExpectStatus(200)

	frisby.Create("Test HRoot").
		Get(fmt.Sprintf("http://localhost:9999/login.html")).
		SetHeader("Accept", "text/html").
		Send().
		ExpectStatus(200)

	frisby.Create("Test HRoot").
		Get(fmt.Sprintf("http://localhost:9999/logout.html")).
		SetHeader("Accept", "text/html").
		Send().
		ExpectStatus(200)

	frisby.Create("Test HRoot").
		Get(fmt.Sprintf("http://localhost:9999/ballot.html")).
		SetHeader("Accept", "text/html").
		Send().
		ExpectStatus(200)

	frisby.Create("Test HRoot").
		Get(fmt.Sprintf("http://localhost:9999/decision/0")).
		SetHeader("Accept", "text/html").
		Send().
		ExpectStatus(200)

	frisby.Create("Test HRoot").
		Get(fmt.Sprintf("http://localhost:9999/results.html")).
		SetHeader("Accept", "text/html").
		Send().
		ExpectStatus(200)

	frisby.Global.PrintReport()
}
