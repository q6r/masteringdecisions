package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Criterion represent a criterion for a decision
type Criterion struct {
	Criterion_ID int    `db:"criterion_id" json:"criterion_id"`
	Decision_ID  int    `db:"decision_id" json:"decision_id"` // inherited
	Name         string `db:"name" json:"name" binding:"required"`
	Weight       int    `db:"weight" json:"weight" binding:"required"`
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
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, cri)
}

// HCriterionDelete deletes a criterion from a decision
func HCriterionDelete(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	cid, err := strconv.Atoi(c.Param("criterion_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	cri := &Criterion{Criterion_ID: cid, Decision_ID: did}
	err = cri.Destroy()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

// HCriterionCreate creates a criterion for a decision
func HCriterionCreate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	var cri Criterion
	err = c.Bind(&cri)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid criterion object"})
		return
	}
	cri.Decision_ID = did // inherited

	err = cri.Save()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cri)
}

// Destroy removes a criterion from a decision
func (cri *Criterion) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM criterion WHERE criterion_id=$1 and decision_id=$2",
		cri.Criterion_ID, cri.Decision_ID)
	if err != nil {
		return err
	}

	return nil
}

// Save saves a criterion in the database
// Restrictions decision should exist
func (cri *Criterion) Save() error {

	// See if there's a decision this belongs to
	var d Decision
	err := dbmap.SelectOne(&d, "select * from decision where decision_id=$1", cri.Decision_ID)
	if err != nil {
		return fmt.Errorf("decision %d does not exist, criterion should belong to an existing decision", cri.Decision_ID)
	}

	if err := dbmap.Insert(cri); err != nil {
		return err
	}
	return nil
}
