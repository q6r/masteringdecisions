package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Alternative represent an alternative in a decision
type Alternative struct {
	AlternativeID int     `db:"alternative_id" json:"alternative_id"`
	DecisionID    int     `db:"decision_id" json:"decision_id"`
	Name          string  `db:"name" json:"name" binding:"required"`
	Description   string  `db:"description" json:"description"`
	Cost          float32 `db:"cost" json:"cost"`
	Order         float32     `db:"order" json:"order"`
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
	alt.DecisionID = did
	err = alt.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"alternative": alt}
	c.Writer.Header().Set("Location", fmt.Sprintf("/decision/%d/alternative/%d", alt.DecisionID, alt.AlternativeID))
	ServeResult(c, "alternative_create.js", result)
}

// Save inserts a ballot into the database
func (alt *Alternative) Save() error {
	// Check if decision exists or not
	if dobj, err := dbmap.Get(Decision{}, alt.DecisionID); err != nil || dobj == nil {
		return fmt.Errorf("Decision %d does not exists, can't create alternative without a decision", alt.DecisionID)
	}

	if err := dbmap.Insert(alt); err != nil {
		return fmt.Errorf("Unable to insert alternative %#v to database", alt)
	}

	return nil
}

// HAlternativeUpdate updates an alternative
func HAlternativeUpdate(c *gin.Context) {
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

	newAlternative := Alternative{
		AlternativeID: aid,
		DecisionID:    did,
		Name:          json.Name,
		Description:   json.Description,
		Cost:          json.Cost,
		Order:         json.Order,
	}
	_, err = dbmap.Update(&newAlternative)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update alternative %d for decision %d", aid, did)})
		return
	}

	result := gin.H{"alternative": newAlternative}
	ServeResult(c, "alternative_update.js", result)
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

	alt := &Alternative{AlternativeID: aid, DecisionID: did}
	err = alt.Destroy()
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	ServeResult(c, "alternative_deleted.js", result)
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
	ServeResult(c, "alternative_info.js", result)
}

// Destroy an alternative
func (alt *Alternative) Destroy() error {
	if _, err := dbmap.Delete(alt); err != nil {
		return fmt.Errorf("Unable to delete alternative %#v from database", alt)
	}

	// Remove votes beloning to this alternative
	var votes []Vote
	_, _ = dbmap.Select(&votes, "SELECT * FROM vote WHERE alternative_id=$1", alt.AlternativeID)
	for _, v := range votes {
		if err := v.Destroy(); err != nil {
			return err
		}
	}

	return nil
}
