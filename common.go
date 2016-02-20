package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServeResults serve the results at the end
// of most handlers
func ServeResult(c *gin.Context, scriptName string, result gin.H) {
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": scriptName, "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
