package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Ballot struct {
	Ballot_ID   int    `db:"ballot_id"`
	Decision_ID int    `db:"decision_id" binding:"required"`
	Secret      int    `db:"secret" binding:"required"`
	Name        string `db:"name" binding:"required"`
	Email       string `db:"email" binding:"required"`
}

func HBallotCreate(c *gin.Context) {
	var b Ballot
	if err := c.Bind(&b); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	err := b.CreateBallot()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, b)
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
	_, err := dbmap.Select(&ballot, "SELECT * FROM ballot where ballot_id=$1", bid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
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
