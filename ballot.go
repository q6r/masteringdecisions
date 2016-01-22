package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Ballot represent a ballot that belong
// in a decision
type Ballot struct {
	Ballot_ID   int    `db:"ballot_id" json:"ballot_id"`
	Decision_ID int    `db:"decision_id" json:"decision_id"`
	Secret      int    `db:"secret" json:"secret" binding:"required"`
	Name        string `db:"name" json:"name" binding:"required"`
	Email       string `db:"email" json:"email" binding:"required"`
}

// HBallotCreate create a ballot that belongs
// to a decision
func HBallotCreate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid decision_id"})
		return
	}

	var b Ballot
	if err := c.Bind(&b); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	b.Decision_ID = did // inherited

	err = b.Save()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

// HBallotDelete deletes a ballot from a decision
func HBallotDelete(c *gin.Context) {

	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	b := &Ballot{Ballot_ID: bid, Decision_ID: did}
	err = b.Destroy()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

// HBallotInfo gets the information for a specific
// ballot in a decision and retusn an json object
// of the found ballot
func HBallotInfo(c *gin.Context) {
	did := c.Param("decision_id")
	bid := c.Param("ballot_id")
	var ballot Ballot
	err := dbmap.SelectOne(&ballot, "SELECT * FROM ballot where ballot_id=$1 and decision_id=$2", bid, did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, ballot)
}

// Destroy removes a ballot from the database
// it also removes the dependencies of a ballots
// such as votes
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
