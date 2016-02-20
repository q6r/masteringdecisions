package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Rating implement ratinging an criterion
// by a ballot
type Rating struct {
	CriterionID int `db:"criterion_id" json:"criterion_id" binding:"required"`
	BallotID    int `db:"ballot_id" json:"ballot_id" binding:"required"`
	Rating      int `db:"rating" json:"rating" binding:"required"`
}

// HRatingCreate creates rating by a ballot on a
// specific criterion
// eg : Ballot 1 ratings Alternative  2 with rating 20
// GET /decision/<n>/ballot/1/criterion/2/vote/20
func HRatingCreate(c *gin.Context) {
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

	rating, err := strconv.Atoi(c.Param("rating"))
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

	r := Rating{CriterionID: cid, BallotID: bid, Rating: rating}
	err = r.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"rating": r}
	c.Writer.Header().Set("Location",
		fmt.Sprintf("/decision/%d/ballot/%d/criterion/%d/vote/",
			b.DecisionID, bid, cid))
	ServeResult(c, "rating_create.js", result)
}

// HRatingBallots shows all the ratings on a criterion
// made by all the ballots belonging to that decision
// GET /decision/<n>/ballot/1/criterion/2/votes
func HRatingBallots(c *gin.Context) {
	cid, err := strconv.Atoi(c.Param("criterion_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var rs []Rating
	_, err = dbmap.Select(&rs, "select * from rating WHERE criterion_id=$1",
		cid)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"ratings": rs}
	ServeResult(c, "rating_ballot_list.js", result)
}

// HRatingDelete delete a specific rating on an criterion
// by a ballot
func HRatingDelete(c *gin.Context) {
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

	r := Rating{CriterionID: cid, BallotID: bid}
	err = r.Destroy()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	ServeResult(c, "rating_deleted.js", result)
}

// HRatingUpdate update a specific rating on a criterion
// by a ballot
func HRatingUpdate(c *gin.Context) {
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
	rating, err := strconv.Atoi(c.Param("rating"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var cri Criterion
	err = dbmap.SelectOne(&cri,
		"SELECT * FROM criterion WHERE criterion_id=$1", cid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update rating for ballot %d and criterion %d", bid, cid)})
		return
	}

	_, err = dbmap.Exec("UPDATE rating SET rating=$1 WHERE ballot_id=$2 and criterion_id=$3", rating, bid, cid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update rating for ballot %d and criterion %d", bid, cid)})
		return
	}

	newRating := Rating{CriterionID: cid, BallotID: bid, Rating: rating}
	result := gin.H{"rating": newRating}
	ServeResult(c, "rating_update.js", result)
}

// Save saves a rating in the database
func (w *Rating) Save() error {

	// No duplicate ratings
	n, _ := dbmap.SelectInt("select count(*) from rating where ballot_id=$1 and criterion_id=$2", w.BallotID, w.CriterionID)
	if n >= 1 {
		return fmt.Errorf("rating %#v already exists", w)
	}

	var b Ballot
	err := dbmap.SelectOne(&b, "select * from ballot where ballot_id=$1", w.BallotID)
	if err != nil {
		return fmt.Errorf("ballot %d does not exists, can't create a rating without an owner", w.BallotID)
	}

	var cri Criterion
	err = dbmap.SelectOne(&cri, "select * from criterion where criterion_id=$1", w.CriterionID)
	if err != nil {
		return fmt.Errorf("criterion %d does not exists, can't create a vote that doesn't belong to a criterion",
			w.CriterionID)
	}

	// Make sure the criterion and ballot belong to the same decision
	if cri.DecisionID != b.DecisionID {
		return fmt.Errorf("The criterion and ballot don't belong to this decision")
	}

	if err = dbmap.Insert(w); err != nil {
		return fmt.Errorf("unable to insert rating %#v to database", w)
	}

	return nil
}

// Destroy removes a rating from the database
func (w *Rating) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM rating WHERE ballot_id=$1 and criterion_id=$2",
		w.BallotID, w.CriterionID)
	if err != nil {
		return err
	}
	return nil
}
