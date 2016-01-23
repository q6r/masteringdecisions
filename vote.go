package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Vote represent a vote by a ballot for a specific
// criterion
type Vote struct {
	Criterion_ID int `db:"criterion_id" json:"criterion_id" required:"binding"`
	Ballot_ID    int `db:"ballot_id" json:"ballot_id" required:"binding"`
	Weight       int `db:"weight" json:"weight" required:"binding"`
}

// HVoteCreate a ballot votes on a criterion
// TODO : Force weight checking on criterion
// the weight in the vote should not be higher than the
// weight defined in the criterion
func HVoteCreate(c *gin.Context) {

	cid, err := strconv.Atoi(c.Param("criterion_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	weight, err := strconv.Atoi(c.Param("weight"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	v := Vote{Criterion_ID: cid, Ballot_ID: bid, Weight: weight}

	err = v.Save()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, v)
}

// HVoteDelete deletes a vote by a ballot
func HVoteDelete(c *gin.Context) {
	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	cid, err := strconv.Atoi(c.Param("criterion_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	v := Vote{Ballot_ID: bid, Criterion_ID: cid}
	err = v.Destroy()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

// HVotesBallotList list all votes made by a ballot
func HVotesBallotList(c *gin.Context) {
	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var vs []Vote
	_, err = dbmap.Select(&vs, "select * from vote WHERE ballot_id=$1", bid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vs)
}

// Destroy removes a vote from the database
func (v *Vote) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM vote WHERE ballot_id=$1 and criterion_id=$2", v.Ballot_ID, v.Criterion_ID)
	if err != nil {
		return err
	}
	return nil
}

// Save a vote in the database
// Restriction : Criterion should exists
// Restriction : Ballot should exists
// Restriction : Don't allow duplicates on ballot_id, criterion_id
// Restriction : Make sure the criterion and ballot we're voting for belongs to the same decision
func (v *Vote) Save() error {

	// No duplicate votes
	n, err := dbmap.SelectInt("select count(*) from vote where ballot_id=$1 and criterion_id=$2", v.Ballot_ID, v.Criterion_ID)
	if n >= 1 {
		return fmt.Errorf("vote %#v already exists", v)
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
		return fmt.Errorf("Unable to insert vote %#v to database", v)
	}

	return nil
}

// FindVotesByKeys find votes by keys
func FindVotesByKeys(criterion_id, ballot_id int) (Vote, error) {
	var vote Vote
	err := dbmap.SelectOne(&vote, "select * from vote where criterion_id=$1 and ballot_id=$2", criterion_id, ballot_id)
	if err != nil {
		return Vote{},
			fmt.Errorf("Unable to find vote for criterion_id=%v ballot_id=%v", criterion_id, ballot_id)
	}
	return vote, nil
}
