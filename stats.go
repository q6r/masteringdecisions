package main

// TODO : Finish implementation
// TODO : Add stats to database

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Stat struct {
	Ballot_ID         int     `json:"ballot_id"`
	Criterion_ID      int     `json:"criterion_id"`
	Decision_ID       int     `json:"decision_id"`
	User_Weight       float64 `json:"user_weight"`
	Criterion_Weight  float64 `json:"criterion_weight"`
	Percentage_Weight float64 `json:"percentage_weight"`
}

func HStats(c *gin.Context) {
	did := c.Param("decision_id")

	// Get the decision
	var decision Decision
	err := dbmap.SelectOne(&decision,
		"select * from decision where decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	// Get ballots for that decision
	var ballots []Ballot
	_, err = dbmap.Select(&ballots,
		"select * from ballot where decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	// Get criterions for that decisions
	var criterions []Criterion
	_, err = dbmap.Select(&criterions,
		"select * from criterion where decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	// find votes for this decision by ballot
	var all_votes [][]Vote
	for _, ballot := range ballots {
		var votes []Vote
		_, err = dbmap.Select(&votes,
			"select * from vote where ballot_id=$1",
			ballot.Ballot_ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		all_votes = append(all_votes, votes)
	}

	// At this points we have
	// 1. The decision
	// 2. Its ballots
	// 3. Its criterions
	// 4. Its votes grouped by (by ballots)

	// TODO : Do the math save in stats insert stats
	// ..etc

	c.JSON(http.StatusOK, "ok")
}
