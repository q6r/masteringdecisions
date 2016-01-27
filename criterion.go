package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("criterion id %v for decision id %v not found", cid, did)})
		return
	}

	result := gin.H{"criterion": cri}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "criterion_info.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// HCriterionDelete deletes a criterion from a decision
func HCriterionDelete(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	cid, err := strconv.Atoi(c.Param("criterion_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	cri := &Criterion{Criterion_ID: cid, Decision_ID: did}
	err = cri.Destroy()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "criterion_deleted.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// HCriterionCreate creates a criterion for a decision
func HCriterionCreate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

	result := gin.H{"criterion": cri}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "criterion_create.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// HCriterionUpdate updates a criterion
func HCriterionUpdate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	cid, err := strconv.Atoi(c.Param("criterion_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var cri Criterion
	err = dbmap.SelectOne(&cri, "SELECT * FROM criterion WHERE decision_id=$1 and criterion_id=$2", did, cid)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": fmt.Sprintf("criterion %d for decision %d not found", cid, did)})
		return
	}

	var json Criterion
	err = c.Bind(&json)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": "Unable to parse decision object"})
		return
	}

	new_criterion := Criterion{
		Criterion_ID: cid,
		Decision_ID:  did,
		Name:         json.Name,
		Weight:       json.Weight,
	}
	_, err = dbmap.Update(&new_criterion)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": fmt.Sprintf("Unable to update criterion %d for decision %d", cid, did)})
		return
	}

	result := gin.H{"criterion": new_criterion}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "criterion_update.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// Destroy removes a criterion from a decision
func (cri *Criterion) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM criterion WHERE criterion_id=$1 and decision_id=$2",
		cri.Criterion_ID, cri.Decision_ID)
	if err != nil {
		return fmt.Errorf("Unable to delete criterion %#v from database", cri)
	}

	return nil
}

// Save saves a criterion in the database
// Restrictions decision should exist
// TODO : Don't allow duplication ?
func (cri *Criterion) Save() error {
	// See if there's a decision this belongs to
	var d Decision
	err := dbmap.SelectOne(&d, "select * from decision where decision_id=$1", cri.Decision_ID)
	if err != nil {
		return fmt.Errorf("decision %d does not exist, criterion should belong to an existing decision", cri.Decision_ID)
	}

	if err := dbmap.Insert(cri); err != nil {
		return fmt.Errorf("Unable to insert criterion %#v to database", cri)
	}

	return nil
}
