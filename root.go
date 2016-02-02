package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HRootHome(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "page_index.js", "body": gin.H{}})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func HRootLogin(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "page_login.js", "body": gin.H{}})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func HRootLogout(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "page_logout.js", "body": gin.H{}})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func HRootBallot(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "page_ballot.js", "body": gin.H{}})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func HRootDecision(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "page_decision.js", "body": gin.H{}})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}
