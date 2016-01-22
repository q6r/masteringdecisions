package gobackend

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Decision represent a decision owned by Person_ID
type Decision struct {
	Decision_ID          int    `db:"decision_id" json:"decision_id"`
	Person_ID            int    `db:"person_id" json:"person_id" binding:"required"`
	Name                 string `db:"name" json:"name" binding:"required"`
	Description          string `db:"description" json:"description" binding:"required"`
	Stage                int    `db:"stage" json:"stage" binding:"required"`
	Criterion_Vote_Style string `db:"criterion_vote_style" json:"criterion_vote_style" binding:"required"`
	// TODO : remove alternative?
	Alternative_Vote_Style string `db:"alternative_vote_style" json:"alternative_vote_style" binding:"required"`
	Client_Settings        string `db:"client_settings" json:"client_settings" binding:"required"`
}

// HDecisionStats Finds all ballots beloning to a decision
// find their voes, find the criterion of the decision
// do the math, return a stats object
// TODO : implement me / move to stats.go
func HDecisionStats(c *gin.Context) {
}

// HDecisionBallotsList returns a list of ballots beloning
// to a decision
func HDecisionBallotsList(c *gin.Context) {
	did := c.Param("decision_id")
	var ballots []Ballot
	_, err := dbmap.Select(&ballots, "SELECT * FROM ballot WHERE decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, ballots)
}

// HDecisionCriterionsList returns a list of criterions beloning
// to a decision
func HDecisionCriterionsList(c *gin.Context) {
	did := c.Param("decision_id")
	var cris []Criterion
	_, err := dbmap.Select(&cris, "select * from criterion where decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, cris)
}

// HDecisionsList returns a list of all decision defined
// in the database
func HDecisionsList(c *gin.Context) {
	var decisions []Decision
	_, err := dbmap.Select(&decisions, "SELECT * FROM decision")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, decisions)
}

// HDecisionInfo returns a decision information
// a decision object not it's stats
func HDecisionInfo(c *gin.Context) {
	id := c.Param("decision_id")
	var decision Decision
	err := dbmap.SelectOne(&decision, "SELECT * FROM decision where decision_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, decision)
}

// HDecisionCreate creates a decision beloning to a specific
// person
func HDecisionCreate(c *gin.Context) {
	var decision Decision
	err := c.Bind(&decision)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	err = decision.Save()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, decision)
}

// HDecisionDelete deletes a decision from database
func HDecisionDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	d := &Decision{Decision_ID: id}
	err = d.Destroy()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "deleted decision, its ballots"})
}

// Destroy a decision from the database
// and remove it's dependencies such as ballots
// when destroying ballots they'll destroy their votes..etc
func (d *Decision) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM decision WHERE decision_id=$1", d.Decision_ID)
	if err != nil {
		return err
	}

	var ballots []Ballot
	_, err = dbmap.Select(&ballots, "SELECT * FROM ballot WHERE decision_id=$1", d.Decision_ID)
	if err != nil {
		return err
	}

	for _, b := range ballots {
		err := b.Destroy()
		if err != nil {
			return err
		}
	}

	var cris []Criterion
	_, err = dbmap.Select(&cris, "select * from criterion where decision_id=$1", d.Decision_ID)
	if err != nil {
		return err
	}

	for _, cri := range cris {
		err := cri.Destroy()
		if err != nil {
			return err
		}
	}

	return nil
}

// Save saves decision in the database
// Restriction : Decision can't be created without an existing owner
// Restriction : Decision can't be owned by two different people
func (d *Decision) Save() error {

	// See if there's a person that this decision belongs to
	// otherwise we quit
	var p Person
	err := dbmap.SelectOne(&p, "SELECT * FROM person WHERE person_id=$1", d.Person_ID)
	if err != nil {
		return fmt.Errorf("person %d does not exist, can't create a decision without an owner", d.Person_ID)
	}

	// If someone else other than us owns the same
	// decision then abort
	var ds []Decision
	_, err = dbmap.Select(&ds, "select * from decision where decision_id=$1", d.Decision_ID)
	if err != nil {
		return err
	}
	for _, i := range ds {
		if i.Person_ID != d.Person_ID {
			return fmt.Errorf("decision %d already owned by person %d", d.Decision_ID, i.Person_ID)
		}
	}

	err = dbmap.Insert(d)
	if err != nil {
		return err
	}
	return nil
}
