package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type RequestFunc func(*gin.Context)
type ResponseFunc func(*httptest.ResponseRecorder)

type JResult struct {
	Result string `json:"result"`
}

type RequestConfig struct {
	Method      string
	Path        string
	Body        string
	Headers     map[string]string
	Middlewares []gin.HandlerFunc
	Handler     RequestFunc
	Finaliser   ResponseFunc
}

func RunRequest(rc RequestConfig) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	if rc.Middlewares != nil && len(rc.Middlewares) > 0 {
		for _, mw := range rc.Middlewares {
			r.Use(mw)
		}
	}

	qs := ""
	if strings.Contains(rc.Path, "?") {
		ss := strings.Split(rc.Path, "?")
		rc.Path = ss[0]
		qs = ss[1]
	}

	body := bytes.NewBufferString(rc.Body)

	req, _ := http.NewRequest(rc.Method, rc.Path, body)

	if len(qs) > 0 {
		req.URL.RawQuery = qs
	}

	if len(rc.Headers) > 0 {
		for k, v := range rc.Headers {
			req.Header.Set(k, v)
		}
	} else if rc.Method == "POST" || rc.Method == "PUT" {
		if strings.HasPrefix(rc.Body, "{") {
			req.Header.Set("Content-Type", "application/json")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	r.Handle(rc.Method, rc.Path, func(c *gin.Context) {
		//change argument if necessary here
		rc.Handler(c)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if rc.Finaliser != nil {
		rc.Finaliser(w)
	}
}

func RunSimpleGet(path string, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:    "GET",
		Path:      path,
		Handler:   handler,
		Finaliser: reply,
	}
	RunRequest(rc)
}

func RunSimpleDelete(path string, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:    "DELETE",
		Path:      path,
		Handler:   handler,
		Finaliser: reply,
	}
	RunRequest(rc)
}

func RunSimplePost(path string, body string,
	handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:    "POST",
		Path:      path,
		Body:      body,
		Handler:   handler,
		Finaliser: reply,
	}
	RunRequest(rc)
}

func TCreatePerson(data string) (Person, error) {
	var p1 Person
	var err error

	RunSimplePost("/person", data,
		func(c *gin.Context) {
			HPersonCreate(c)
		},
		func(r *httptest.ResponseRecorder) {
			err = json.Unmarshal(r.Body.Bytes(), &p1)
		})

	if err != nil {
		return Person{}, err
	}

	return p1, nil
}

func TCreateCriterion(decision_id int, data string) (Criterion, error) {
	var c1 Criterion
	var err error

	RunSimplePost("/decision/:decision_id/criterion", data,
		func(c *gin.Context) {
			c.Params = gin.Params{
				gin.Param{
					Key:   "decision_id",
					Value: strconv.Itoa(decision_id),
				},
			}
			HCriterionCreate(c)
		},
		func(r *httptest.ResponseRecorder) {
			err = json.Unmarshal(r.Body.Bytes(), &c1)
		})

	if err != nil {
		return Criterion{}, err
	}

	return c1, nil
}

func TCreateDecision(data string) (Decision, error) {
	var d1 Decision
	var err error

	RunSimplePost("/decision", data,
		func(c *gin.Context) {
			HDecisionCreate(c)
		},
		func(r *httptest.ResponseRecorder) {
			err = json.Unmarshal(r.Body.Bytes(), &d1)
		})

	if err != nil {
		return Decision{}, err
	}

	return d1, nil
}

func TCreateBallot(did int, data string) (Ballot, error) {
	var b1 Ballot
	var err error

	RunSimplePost("/decision/:did/ballot", data,
		func(c *gin.Context) {
			c.Params = gin.Params{
				gin.Param{
					Key:   "decision_id",
					Value: strconv.Itoa(did),
				},
			}
			HBallotCreate(c)
		},
		func(r *httptest.ResponseRecorder) {
			err = json.Unmarshal(r.Body.Bytes(), &b1)
		})

	if err != nil {
		return Ballot{}, err
	}

	return b1, nil
}

func TDeleteCriterion(decision_id int, criterion_id int) (JResult, error) {
	var res JResult
	var err error

	RunSimpleDelete("/decision/:decision_id/criterion/:ballot_id",
		func(c *gin.Context) {
			c.Params = gin.Params{
				gin.Param{
					Key:   "decision_id",
					Value: strconv.Itoa(decision_id),
				},
				gin.Param{
					Key:   "criterion_id",
					Value: strconv.Itoa(criterion_id),
				},
			}
			HCriterionDelete(c)
		},
		func(r *httptest.ResponseRecorder) {
			err = json.Unmarshal(r.Body.Bytes(), &res)
		})
	if err != nil {
		return JResult{}, err
	}

	return res, nil
}

func TDeleteBallot(decision_id int, ballot_id int) (JResult, error) {
	var res JResult
	var err error

	RunSimpleDelete("/decision/:decision_id/ballot/:ballot_id",
		func(c *gin.Context) {
			c.Params = gin.Params{
				gin.Param{
					Key:   "decision_id",
					Value: strconv.Itoa(decision_id),
				},
				gin.Param{
					Key:   "ballot_id",
					Value: strconv.Itoa(ballot_id),
				},
			}
			HBallotDelete(c)
		},
		func(r *httptest.ResponseRecorder) {
			err = json.Unmarshal(r.Body.Bytes(), &res)
		})
	if err != nil {
		return JResult{}, err
	}

	return res, nil
}

func TDeleteDecision(decision_id int) (JResult, error) {
	var res JResult
	var err error

	RunSimpleDelete("/decision/:decision_id",
		func(c *gin.Context) {
			c.Params = gin.Params{gin.Param{Key: "decision_id", Value: strconv.Itoa(decision_id)}}
			HDecisionDelete(c)
		},
		func(r *httptest.ResponseRecorder) {
			err = json.Unmarshal(r.Body.Bytes(), &res)
		})
	if err != nil {
		return JResult{}, err
	}

	return res, nil
}

func TDeletePerson(person_id int) (JResult, error) {
	var res JResult
	var err error

	RunSimpleDelete("/person/:person_id",
		func(c *gin.Context) {
			c.Params = gin.Params{gin.Param{Key: "person_id", Value: strconv.Itoa(person_id)}}
			HPersonDelete(c)
		},
		func(r *httptest.ResponseRecorder) {
			err = json.Unmarshal(r.Body.Bytes(), &res)
		})
	if err != nil {
		return JResult{}, err
	}

	return res, nil
}

func TInfoPerson(person_id int) (Person, error) {
	var res Person
	var err error

	RunSimpleDelete("/person/:person_id/info",
		func(c *gin.Context) {
			c.Params = gin.Params{gin.Param{Key: "person_id", Value: strconv.Itoa(person_id)}}
			HPersonInfo(c)
		},
		func(r *httptest.ResponseRecorder) {
			err = json.Unmarshal(r.Body.Bytes(), &res)
		})
	if err != nil {
		return Person{}, err
	}

	return res, nil

}

func TInfoBallot(decision_id int, ballot_id int) (Ballot, error) {
	var res Ballot
	var err error

	RunSimpleDelete("/decision/:decision_id/ballot/:ballot_id/info",
		func(c *gin.Context) {
			c.Params = gin.Params{
				gin.Param{
					Key:   "decision_id",
					Value: strconv.Itoa(decision_id),
				},
				gin.Param{
					Key:   "ballot_id",
					Value: strconv.Itoa(ballot_id),
				},
			}

			HBallotInfo(c)
		},
		func(r *httptest.ResponseRecorder) {
			err = json.Unmarshal(r.Body.Bytes(), &res)
		})

	if err != nil {
		return Ballot{}, err
	}

	return res, nil
}
