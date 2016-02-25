package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Criterion represent a criterion for a decision
type Criterion struct {
	CriterionID int       `db:"criterion_id" json:"criterion_id"`
	DecisionID  int       `db:"decision_id" json:"decision_id"` // inherited
	Name        string    `db:"name" json:"name" binding:"required"`
	Description string    `db:"description" json:"description"`
	Order       float32   `db:"order" json:"order"`
}

// HCriterionInfo get the information of a specific
// criterion in a decision and return it as a json
// object
func HCriterionInfo(c *gin.Context) {
	did := c.Param("decision_id")
	cid := c.Param("criterion_id")

	var cri Criterion
	err := dbmap.SelectOne(&cri, "select * from criterion where criterion_id=$1 and decision_id=$2", cid, did)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("criterion id %v for decision id %v not found", cid, did)})
		return
	}

	result := gin.H{"criterion": cri}
	ServeResult(c, "criterion_info.js", result)
}

// HCriterionDelete deletes a criterion from a decision
func HCriterionDelete(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	cid, err := strconv.Atoi(c.Param("criterion_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	cri := &Criterion{CriterionID: cid, DecisionID: did}
	err = cri.Destroy()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	ServeResult(c, "criterion_deleted.js", result)
}

// HCriterionCreate creates a criterion for a decision
func HCriterionCreate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var cri Criterion
	err = c.Bind(&cri)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid criterion object"})
		return
	}
	cri.DecisionID = did // inherited

	err = cri.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"criterion": cri}
	c.Writer.Header().Set("Location",
		fmt.Sprintf("/decision/%d/criterion/%d", did, cri.CriterionID))
	ServeResult(c, "criterion_create.js", result)
}

// HCriterionUpdate updates a criterion
func HCriterionUpdate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	cid, err := strconv.Atoi(c.Param("criterion_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var cri Criterion
	err = dbmap.SelectOne(&cri, "SELECT * FROM criterion WHERE decision_id=$1 and criterion_id=$2", did, cid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("criterion %d for decision %d not found", cid, did)})
		return
	}

	var json Criterion
	err = c.Bind(&json)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": "Unable to parse decision object"})
		return
	}

	newCriterion := Criterion{
		CriterionID: cid,
		DecisionID:  did,
		Name:        json.Name,
		Description: json.Description,
		Order:       json.Order,
	}
	_, err = dbmap.Update(&newCriterion)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update criterion %d for decision %d", cid, did)})
		return
	}

	result := gin.H{"criterion": newCriterion}
	ServeResult(c, "criterion_update.js", result)
}

// Destroy removes a criterion from a decision
func (cri *Criterion) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM criterion WHERE criterion_id=$1 and decision_id=$2",
		cri.CriterionID, cri.DecisionID)
	if err != nil {
		return fmt.Errorf("Unable to delete criterion %#v from database", cri)
	}

	// Remove the votes that have a vote for this destroied
	// criterion and destroy them
	var votes []Vote
	_, _ = dbmap.Select(&votes, "select * from vote where criterion_id=$1", cri.CriterionID)
	for _, vote := range votes {
		if err := vote.Destroy(); err != nil {
			return err
		}
	}

	// Remove the ratings that have a rate for this destroied
	// criterion and destroy them
	var ratings []Rating
	_, _ = dbmap.Select(&ratings, "select * from rating where criterion_id=$1", cri.CriterionID)
	for _, rating := range ratings {
		if err := rating.Destroy(); err != nil {
			return err
		}
	}

	return nil
}

// Save saves a criterion in the database
// Restrictions decision should exist
func (cri *Criterion) Save() error {
	// See if there's a decision this belongs to
	cobj, err := dbmap.Get(Decision{}, cri.DecisionID)
	if err != nil || cobj == nil {
		return fmt.Errorf("decision %d does not exist, criterion should belong to an existing decision", cri.DecisionID)
	}

	if err := dbmap.Insert(cri); err != nil {
		return fmt.Errorf("Unable to insert criterion %#v to database", cri)
	}

	return nil
}
