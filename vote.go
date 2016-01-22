package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Vote struct {
	Criterion_ID int `db:"criterion_id" json:"criterion_id" required:"binding"`
	Ballot_ID    int `db:"ballot_id" json:"ballot_id" required:"binding"`
	Weight       int `db:"weight" json:"weight" required:"binding"`
}

// TODO : Force weight checking on criterion
// the weight in the vote should not be higher than the
// weight defined in the criterion
func HVoteCreate(c *gin.Context) {
	var v Vote
	err := c.Bind(&v)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	err = v.Save()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, v)
}

// requires ballot_id, vote_id
func HVoteDelete(c *gin.Context) {
	type DelReq struct {
		Ballot_ID    int `json:"ballot_id" binding:"required"`
		Criterion_ID int `json:"criterion_id" binding:"required"`
	}
	var delreq DelReq

	err := c.Bind(&delreq)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	v := Vote{Ballot_ID: delreq.Ballot_ID, Criterion_ID: delreq.Criterion_ID}
	err = v.Destroy()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

func HVotesList(c *gin.Context) {
	var vs []Vote
	_, err := dbmap.Select(&vs, "select * from vote")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, vs)
}

// Requirements : none ?
// TODO : should be called inside ballot eg : removing a ballot removes its votes
func (v *Vote) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM vote WHERE ballot_id=$1 and criterion_id=$2", v.Ballot_ID, v.Criterion_ID)
	if err != nil {
		return err
	}
	return nil
}

// Restriction : Criterion should exists
// Restriction : Ballot should exists
// Restriction : Don't allow duplicates on ballot_id, criterion_id
// Restriction : Make sure the criterion and ballot we're voting for belongs to the same decision
func (v *Vote) Save() error {

	// No duplicate votes
	n, err := dbmap.SelectInt("select count(*) from vote where ballot_id=$1 and criterion_id=$2", v.Ballot_ID, v.Criterion_ID)
	if n >= 1 {
		return fmt.Errorf("vote already exists.")
	}

	// See if there's a criterion that this vote belongs to
	var cri Criterion
	err = dbmap.SelectOne(&cri, "select * from criterion where criterion_id=$1",
		v.Criterion_ID)
	if err != nil {
		return fmt.Errorf("criterion %d does not exist, can't create a vote without an owner", v.Criterion_ID)
	}

	// See if there's a ballot that this vote belongs to
	var b Ballot
	err = dbmap.SelectOne(&b, "select * from ballot where ballot_id=$1",
		v.Ballot_ID)
	if err != nil {
		return fmt.Errorf("ballot %d does not exists, can't create a vote without an owner", v.Ballot_ID)
	}

	// Make sure the criterion and ballot belong to the same decision
	if cri.Decision_ID != b.Decision_ID {
		return fmt.Errorf("criterion belongs to decision %d while ballot belongs to decision %d", cri.Decision_ID, b.Decision_ID)
	}

	err = dbmap.Insert(v)
	if err != nil {
		return err
	}

	return nil
}
