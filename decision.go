package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Decision struct {
	Decision_ID            int    `db:"decision_id"`
	Person_ID              int    `db:"person_id" binding:"required"`
	Name                   string `db:"name" binding:"required"`
	Description            string `db:"description" binding:"required"`
	Owner_ID               int    `db:"owner_id" binding:"required"`
	Stage                  int    `db:"stage" binding:"required"`
	Criterion_Vote_Style   string `db:"criterion_vote_style" binding:"required"`
	Alternative_Vote_Style string `db:"alternative_vote_style" binding:"required"`
	Client_Settings        string `db:"client_settings" binding:"required"`
}

func HDecisionBallots(c *gin.Context) {
	did := c.Param("decision_id")
	var ballots []Ballot
	_, err := dbmap.Select(&ballots, "SELECT * FROM ballot WHERE decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, ballots)
}

func HDecisionsList(c *gin.Context) {
	var decisions []Decision
	_, err := dbmap.Select(&decisions, "SELECT * FROM decision")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, decisions)
}

func HDecisionInfo(c *gin.Context) {
	id := c.Param("decision_id")
	var decision Decision
	_, err := dbmap.Select(&decision, "SELECT * FROM decision where decision_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, decision)
}

func HDecisionCreate(c *gin.Context) {
	var decision Decision
	err := c.Bind(&decision)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	err = decision.CreateDecision()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, decision)
}

func HDecisionDelete(c *gin.Context) {
	id := c.Param("decision_id")
	_, err := dbmap.Exec("DELETE FROM decision WHERE decision_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

func (d *Decision) CreateDecision() error {
	err := dbmap.Insert(d)
	if err != nil {
		return err
	}
	return nil
}
