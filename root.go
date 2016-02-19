package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HRootHome is the main index page
func HRootHome(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "page_index.js", "body": gin.H{}})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

// HRootLogin is the main login page
func HRootLogin(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "page_login.js", "body": gin.H{}})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

// HRootLogout is the main logout page
func HRootLogout(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "page_logout.js", "body": gin.H{}})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

// HRootBallot is the main ballot page
func HRootBallot(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "page_ballot.js", "body": gin.H{}})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

// HRootDecision is the main decision page
func HRootDecision(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "page_decision.js", "body": gin.H{}})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}
