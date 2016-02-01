package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Rating implement ratinging an altenrative
// by a ballot
type Rating struct {
	Criterion_ID int `db:"criterion_id" json:"criterion_id" binding:"required"`
	Ballot_ID    int `db:"ballot_id" json:"ballot_id" binding:"required"`
	Rating       int `db:"rating" json:"rating" binding:"required"`
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

	r := Rating{Criterion_ID: cid, Ballot_ID: bid, Rating: rating}
	err = r.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"rating": r}
	c.Writer.Header().Set("Location",
		fmt.Sprintf("/decision/%d/ballot/%d/criterion/%d/vote/",
			b.Decision_ID, bid, cid))
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "rating_create.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}

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
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "rating_ballot_list.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// HRatingDelete delete a specific rating on an alternative
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

	r := Rating{Criterion_ID: cid, Ballot_ID: bid}
	err = r.Destroy()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "rating_deleted.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
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

	if rating > cri.Weight {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("New rating can't be more than %d", alt.Weight)})
		return
	}

	_, err = dbmap.Exec("UPDATE rating SET rating=$1 WHERE ballot_id=$2 and criterion_id=$3", rating, bid, cid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update rating for ballot %d and criterion %d", bid, cid)})
		return
	}

	new_rating := Rating{Criterion_ID: cid, Ballot_ID: bid, Rating: rating}
	result := gin.H{"rating": new_rating}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "rating_update.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// Save saves a rating in the database
func (w *Rating) Save() error {

	// No duplicate ratings
	n, err := dbmap.SelectInt("select count(*) from rating where ballot_id=$1 and criterion_id=$2", w.Ballot_ID, w.Criterion_ID)
	if n >= 1 {
		return fmt.Errorf("rating %#v already exists", w)
	}

	var b Ballot
	err = dbmap.SelectOne(&b, "select * from ballot where ballot_id=$1", w.Ballot_ID)
	if err != nil {
		return fmt.Errorf("ballot %d does not exists, can't create a rating without an owner", w.Ballot_ID)
	}

	var cri Criterion
	err = dbmap.SelectOne(&alt, "select * from criterion where criterion_id=$1", w.Criterion_ID)
	if err != nil {
		return fmt.Errorf("criterion %d does not exists, can't create a vote that doesn't belong to a criterion",
			w.Criterion_ID)
	}

	// Make sure the criterion and ballot belong to the same decision
	if cri.Decision_ID != b.Decision_ID {
		return fmt.Errorf("The criterion and ballot don't belong to this decision")
	}

	// Make sure the rating is not more than the alternative rating
	if w.Rating > cri.Weight {
		return fmt.Errorf("The rating is more than the maximum defined %d", cri.Weight)
	}

	err = dbmap.Insert(w)
	if err != nil {
		return fmt.Errorf("unable to insert rating %#v to database", w)
	}

	return nil
}

// Destroy removes a rating from the database
func (w *Rating) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM rating WHERE ballot_id=$1 and criterion_id=$2",
		w.Ballot_ID, w.Criterion_ID)
	if err != nil {
		return err
	}
	return nil
}
