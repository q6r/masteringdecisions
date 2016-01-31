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
	Alternative_ID int `db:"alternative_id" json:"alternative_id" binding:"required"`
	Ballot_ID      int `db:"ballot_id" json:"ballot_id" binding:"required"`
	Rating         int `db:"rating" json:"rating" binding:"required"`
}

// HRatingCreate creates rating by a ballot on a
// specific alternative
// eg : Ballot 1 ratings Alternative  2 with rating 20
// GET /decision/<n>/ballot/1/alternative/2/vote/20
func HRatingCreate(c *gin.Context) {
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

	r := Rating{Alternative_ID: aid, Ballot_ID: bid, Rating: rating}
	err = r.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"rating": r}
	c.Writer.Header().Set("Location",
		fmt.Sprintf("/decision/%d/ballot/%d/alternative/%d/vote/",
			b.Decision_ID, bid, aid))
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "rating_create.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}

}

// HRatingBallots shows all the ratings on an alternative
// made by all the ballots belonging to that decision
// GET /decision/<n>/ballot/1/alternative/2/votes
func HRatingBallots(c *gin.Context) {
	aid, err := strconv.Atoi(c.Param("alternative_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var rs []Rating
	_, err = dbmap.Select(&rs, "select * from rating WHERE alternative_id=$1",
		aid)
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

	r := Rating{Alternative_ID: aid, Ballot_ID: bid}
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

// HRatingUpdate update a specific rating on an alternative
// by a ballot
func HRatingUpdate(c *gin.Context) {
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
	rating, err := strconv.Atoi(c.Param("rating"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var alt Alternative
	err = dbmap.SelectOne(&alt,
		"SELECT * FROM alternative WHERE alternative_id=$1", aid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update rating for ballot %d and alternative %d", bid, aid)})
		return
	}

	fmt.Printf("Alt : %#v\n", alt)

	if rating > alt.Rating {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("New rating can't be more than %d", alt.Rating)})
		return
	}

	_, err = dbmap.Exec("UPDATE rating SET rating=$1 WHERE ballot_id=$2 and alternative_id=$3", rating, bid, aid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update rating for ballot %d and alternative %d", bid, aid)})
		return
	}

	new_rating := Rating{Alternative_ID: aid, Ballot_ID: bid, Rating: rating}
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
	n, err := dbmap.SelectInt("select count(*) from rating where ballot_id=$1 and alternative_id=$2", w.Ballot_ID, w.Alternative_ID)
	if n >= 1 {
		return fmt.Errorf("rating %#v already exists", w)
	}

	var b Ballot
	err = dbmap.SelectOne(&b, "select * from ballot where ballot_id=$1", w.Ballot_ID)
	if err != nil {
		return fmt.Errorf("ballot %d does not exists, can't create a rating without an owner", w.Ballot_ID)
	}

	var alt Alternative
	err = dbmap.SelectOne(&alt, "select * from alternative where alternative_id=$1", w.Alternative_ID)
	if err != nil {
		return fmt.Errorf("alternative %d does not exists, can't create a vote that doesn't belong to an alternative",
			w.Alternative_ID)
	}

	// Make sure the alternative and ballot belong to the same decision
	if alt.Decision_ID != b.Decision_ID {
		return fmt.Errorf("The alternative and ballot don't belong to this decision")
	}

	// Make sure the rating is not more than the alternative rating
	if w.Rating > alt.Rating {
		return fmt.Errorf("The rating is more than the maximum defined %d", alt.Rating)
	}

	err = dbmap.Insert(w)
	if err != nil {
		return fmt.Errorf("unable to insert rating %#v to database", w)
	}

	return nil
}

// Destroy removes a rating from the database
func (w *Rating) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM rating WHERE ballot_id=$1 and alternative_id=$2",
		w.Ballot_ID, w.Alternative_ID)
	if err != nil {
		return err
	}
	return nil
}
