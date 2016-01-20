package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Ballot struct {
	Ballot_ID   int    `db:"ballot_id" json:"ballot_id"`
	Decision_ID int    `db:"decision_id" json:"decision_id" binding:"required"`
	Secret      int    `db:"secret" json:"secret" binding:"required"`
	Name        string `db:"name" json:"name" binding:"required"`
	Email       string `db:"email" json:"email" binding:"required"`
}

func HBallotCreate(c *gin.Context) {
	var b Ballot
	if err := c.Bind(&b); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	// Make sure the Ballot's decision exists otherwise we quit
	var d Decision
	err := dbmap.SelectOne(&d, "SELECT * from decision WHERE decision_id=$1", b.Decision_ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "the decision this ballot belong to does not exist, create it first."})
		return
	}

	err = b.CreateBallot()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, b)
}

func HBallotDelete(c *gin.Context) {
	id := c.Param("ballot_id")
	_, err := dbmap.Exec("DELETE FROM ballot WHERE ballot_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "deleted ballot"})
}

func HBallotList(c *gin.Context) {
	var ballots []Ballot
	_, err := dbmap.Select(&ballots, "SELECT * FROM ballot")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, ballots)

}

func HBallotInfo(c *gin.Context) {
	bid := c.Param("ballot_id")
	var ballot Ballot
	err := dbmap.SelectOne(&ballot, "SELECT * FROM ballot where ballot_id=$1", bid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, ballot)
}

func (b *Ballot) CreateBallot() error {
	if err := dbmap.Insert(b); err != nil {
		return err
	}
	return nil
}
