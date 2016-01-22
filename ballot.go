package main

import (
	"net/http"
	"strconv"

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

	err = b.Save()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func HBallotDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	b := &Ballot{Ballot_ID: id}
	err = b.Destroy()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
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

// Destroy removes a ballot from the database
func (b *Ballot) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM ballot WHERE ballot_id=$1", b.Ballot_ID)
	if err != nil {
		return err
	}

	// Remove votes beloning to this ballot
	var votes []Vote
	_, err = dbmap.Select(&votes, "SELECT * FROM vote WHERE ballot_id=$1", b.Ballot_ID)
	if err != nil {
		return err
	}

	for _, v := range votes {
		err = v.Destroy()
		if err != nil {
			return err
		}
	}
	return nil
}

// Save inserts a ballot into the database
func (b *Ballot) Save() error {
	if err := dbmap.Insert(b); err != nil {
		return err
	}
	return nil
}
