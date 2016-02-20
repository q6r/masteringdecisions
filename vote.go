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
	AlternativeID int `db:"alternative_id" json:"alternative_id" binding:"required"`
	CriterionID   int `db:"criterion_id" json:"criterion_id" binding:"required"`
	BallotID      int `db:"ballot_id" json:"ballot_id" binding:"required"`
	Weight        int `db:"weight" json:"weight" binding:"required"`
}

// HVoteCreate a ballot votes on a criterion on an alternative
func HVoteCreate(c *gin.Context) {

	aid, err := strconv.Atoi(c.Param("alternative_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	cid, err := strconv.Atoi(c.Param("criterion_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	weight, err := strconv.Atoi(c.Param("weight"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var b Ballot
	err = dbmap.SelectOne(&b, "select * from ballot where ballot_id=$1", bid)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	v := Vote{AlternativeID: aid, CriterionID: cid, BallotID: bid, Weight: weight}

	err = v.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"vote": v}
	c.Writer.Header().Set("Location", fmt.Sprintf("/decision/%d/ballot/%d/alternative/%d/criterion/%d/vote/",
		b.DecisionID, bid, aid, cid))
	ServeResult(c, "vote_create.js", result)
}

// HVoteUpdate updates a vote
func HVoteUpdate(c *gin.Context) {
	aid, err := strconv.Atoi(c.Param("alternative_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	cid, err := strconv.Atoi(c.Param("criterion_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	weight, err := strconv.Atoi(c.Param("weight"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var cri Criterion
	err = dbmap.SelectOne(&cri, "SELECT * FROM criterion WHERE criterion_id=$1", cid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update vote for ballot %d and criterion %d", bid, cid)})
		return
	}

	_, err = dbmap.Exec("UPDATE vote SET weight=$1 WHERE criterion_id=$2 and ballot_id=$3 and alternative_id=$4", weight, cid, bid, aid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update vote for ballot %d and criterion %d", bid, cid)})
		return
	}

	newVote := Vote{AlternativeID: aid, CriterionID: cid, BallotID: bid, Weight: weight}
	result := gin.H{"vote": newVote}
	ServeResult(c, "vote_update.js", result)
}

// HVoteDelete deletes a vote by a ballot
func HVoteDelete(c *gin.Context) {
	aid, err := strconv.Atoi(c.Param("alternative_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	cid, err := strconv.Atoi(c.Param("criterion_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	v := Vote{AlternativeID: aid, BallotID: bid, CriterionID: cid}
	err = v.Destroy()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	ServeResult(c, "vote_deleted.js", result)
}

// HVotesBallotList list all votes made by a ballot
func HVotesBallotList(c *gin.Context) {
	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var vs []Vote
	_, err = dbmap.Select(&vs, "select * from vote WHERE ballot_id=$1", bid)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"votes": vs}
	ServeResult(c, "vote_ballot_list.js", result)
}

// Destroy removes a vote from the database
func (v *Vote) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM vote WHERE ballot_id=$1 and criterion_id=$2 and alternative_id=$3",
		v.BallotID, v.CriterionID, v.AlternativeID)
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
	n, _ := dbmap.SelectInt("select count(*) from vote where ballot_id=$1 and criterion_id=$2 and alternative_id=$3",
		v.BallotID, v.CriterionID, v.AlternativeID)
	if n >= 1 {
		return fmt.Errorf("vote %#v already exists", v)
	}

	// See if there's a criterion that this vote belongs to
	var cri Criterion
	err := dbmap.SelectOne(&cri, "select * from criterion where criterion_id=$1",
		v.CriterionID)
	if err != nil {
		return fmt.Errorf("criterion %d does not exist, can't create a vote without an owner", v.CriterionID)
	}

	// See if there's a ballot that this vote belongs to
	var b Ballot
	err = dbmap.SelectOne(&b, "select * from ballot where ballot_id=$1",
		v.BallotID)
	if err != nil {
		return fmt.Errorf("ballot %d does not exists, can't create a vote without an owner", v.BallotID)
	}

	// See if there's an alternative that this vote belongs to
	var alt Alternative
	err = dbmap.SelectOne(&alt, "select * from alternative where alternative_id=$1",
		v.AlternativeID)
	if err != nil {
		return fmt.Errorf("alternative %d does not exists, can't create a vote that doesn't belong to an alternative",
			v.AlternativeID)
	}

	// Make sure the criterion and ballot belong to the same decision
	if cri.DecisionID != b.DecisionID {
		return fmt.Errorf("criterion belongs to decision %d while ballot belongs to decision %d", cri.DecisionID, b.DecisionID)
	}

	if err = dbmap.Insert(v); err != nil {
		return fmt.Errorf("Unable to insert vote %#v to database", v)
	}

	return nil
}
