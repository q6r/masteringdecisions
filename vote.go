package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Vote represent a vote by a ballot for a specific
// criterion
type Vote struct {
	Alternative_ID int `db:"alternative_id" json:"alternative_id" binding:"required"`
	Criterion_ID   int `db:"criterion_id" json:"criterion_id" binding:"required"`
	Ballot_ID      int `db:"ballot_id" json:"ballot_id" binding:"required"`
	Weight         int `db:"weight" json:"weight" binding:"required"`
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

	v := Vote{Alternative_ID: aid, Criterion_ID: cid, Ballot_ID: bid, Weight: weight}

	err = v.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"vote": v}
	c.Writer.Header().Set("Location", fmt.Sprintf("/decision/%d/ballot/%d/alternative/%d/criterion/%d/vote/",
		b.Decision_ID, bid, aid, cid))
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "vote_create.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

/*
	PUT /ballot/:ballot_id/alternative/:alternative_id/criterion/:criterion_id/vote/:weight
*/
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

	new_vote := Vote{Alternative_ID: aid, Criterion_ID: cid, Ballot_ID: bid, Weight: weight}
	result := gin.H{"vote": new_vote}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "vote_update.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
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

	v := Vote{Alternative_ID: aid, Ballot_ID: bid, Criterion_ID: cid}
	err = v.Destroy()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "vote_deleted.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
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
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "vote_ballot_list.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// Destroy removes a vote from the database
func (v *Vote) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM vote WHERE ballot_id=$1 and criterion_id=$2 and alternative_id=$3",
		v.Ballot_ID, v.Criterion_ID, v.Alternative_ID)
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
	n, err := dbmap.SelectInt("select count(*) from vote where ballot_id=$1 and criterion_id=$2 and alternative_id=$3",
		v.Ballot_ID, v.Criterion_ID, v.Alternative_ID)
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

	// See if there's an alternative that this vote belongs to
	var alt Alternative
	err = dbmap.SelectOne(&alt, "select * from alternative where alternative_id=$1",
		v.Alternative_ID)
	if err != nil {
		return fmt.Errorf("alternative %d does not exists, can't create a vote that doesn't belong to an alternative",
			v.Alternative_ID)
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
