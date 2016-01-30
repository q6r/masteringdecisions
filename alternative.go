package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Alternative represent an alternative in a decision
type Alternative struct {
	Alternative_ID int    `db:"alternative_id" json:"alternative_id"`
	Decision_ID    int    `db:"decision_id" json:"decision_id"`
	Name           string `db:"name" json:"name" required:"binding"`
	Rating         int    `db:"rating" json:"rating" required:"binding"`
}

// HAlternativeCreate create a ballot that belongs
// to a decision
func HAlternativeCreate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid decision_id"})
		return
	}

	var alt Alternative
	if err := c.Bind(&alt); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid alternative object"})
		return
	}
	alt.Decision_ID = did
	err = alt.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"alternative": alt}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "alternative_create.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// Save inserts a ballot into the database
func (alt *Alternative) Save() error {
	if err := dbmap.Insert(alt); err != nil {
		return fmt.Errorf("Unable to insert alternative %#v to database", alt)
	}
	return nil
}

// HAlternativeUpdate updates an alternative
func HBallotUpdate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	aid, err := strconv.Atoi(c.Param("alternative_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var alt Alternative
	err = dbmap.SelectOne(&alt, "SELECT * FROM alternative WHERE decision_id=$1 and alternative_id=$2", did, aid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("alternative %d for decision %d not found", aid, did)})
		return
	}

	var json Alternative
	err = c.Bind(&json)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": "Unable to parse alternative object"})
		return
	}

	new_alternative := Alternative{
		Alternative_ID: aid,
		Decision_ID:    did,
		Name:           json.Name,
		Rating:         json.Rating,
	}
	_, err = dbmap.Update(&new_alternative)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update alternative %d for decision %d", aid, did)})
		return
	}

	result := gin.H{"alternative": new_alternative}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "alternative_update.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// HAlternativeDelete deletes an alternative from a decision
func HAlternativeDelete(c *gin.Context) {

	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	aid, err := strconv.Atoi(c.Param("alternative_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	alt := &Alternative{Alternative_ID: aid, Decision_ID: did}
	err = alt.Destroy()
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "alternative_deleted.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// HAlternativeInfo gets the information for a specific
// alternative in a decision and returns an json object
// of the found alternative
func HAlternativeInfo(c *gin.Context) {
	did := c.Param("decision_id")
	aid := c.Param("alternative_id")
	var alt Alternative
	err := dbmap.SelectOne(&alt, "SELECT * FROM alternative where alternative_id=$1 and decision_id=$2", aid, did)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to find alternative %v for decision %v", aid, did)})
		return
	}

	result := gin.H{"alternative": alt}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "alternative_info.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}

}

// Destroy an alternative
func (alt *Alternative) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM alternative WHERE alternative_id=$1", alt.Alternative_ID)
	if err != nil {
		return fmt.Errorf("Unable to delete alternative %#v from database", alt)
	}

	// Remove votes beloning to this alternative
	var votes []Vote
	_, err = dbmap.Select(&votes, "SELECT * FROM vote WHERE alternative_id=$1", alt.Alternative_ID)
	if err != nil {
		return fmt.Errorf("Unable to find votes for alternative %#v", alt)
	}
	for _, v := range votes {
		err = v.Destroy()
		if err != nil {
			return err
		}
	}

	return nil
}
