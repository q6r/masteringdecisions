package main

import "github.com/gin-gonic/gin"

// HRootHome is the main index page
func HRootHome(c *gin.Context) {
	ServeResult(c, "page_index.js", gin.H{})
}

// HRootLogin is the main login page
func HRootLogin(c *gin.Context) {
	ServeResult(c, "page_login.js", gin.H{})
}

// HRootLogout is the main logout page
func HRootLogout(c *gin.Context) {
	ServeResult(c, "page_logout.js", gin.H{})
}

// HRootBallot is the main ballot page
func HRootBallot(c *gin.Context) {
	ServeResult(c, "page_ballot.js", gin.H{})
}

// HRootDecision is the main decision page
func HRootDecision(c *gin.Context) {
	ServeResult(c, "page_decision.js", gin.H{})
}

// HRootResults is the main decision page
func HRootResults(c *gin.Context) {
	ServeResult(c, "page_results.js", gin.H{})
}
