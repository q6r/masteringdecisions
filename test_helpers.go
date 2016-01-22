package gobackend

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

type RequestFunc func(*gin.Context)
type ResponseFunc func(*httptest.ResponseRecorder)

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
