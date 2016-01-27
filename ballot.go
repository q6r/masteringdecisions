package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Ballot represent a ballot that belong
// in a decision
type Ballot struct {
	Ballot_ID   int    `db:"ballot_id" json:"ballot_id"`
	Decision_ID int    `db:"decision_id" json:"decision_id"`
	Secret      string `db:"secret" json:"secret"`
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
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid ballot object"})
		return
	}
	b.Decision_ID = did // inherited

	err = b.Save()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"ballot": b}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "ballot_create.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}

}

// HBallotUpdate updates a ballot
func HBallotUpdate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var b Ballot
	err = dbmap.SelectOne(&b, "SELECT * FROM ballot WHERE decision_id=$1 and ballot_id=$2", did, bid)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": fmt.Sprintf("ballot %d for decision %d not found", bid, did)})
		return
	}

	var json Ballot
	err = c.Bind(&json)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": "Unable to parse decision object"})
		return
	}

	new_ballot := Ballot{
		Ballot_ID:   bid,
		Decision_ID: did,
		Secret:      json.Secret,
		Name:        json.Name,
		Email:       json.Email,
	}
	_, err = dbmap.Update(&new_ballot)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": fmt.Sprintf("Unable to update ballot %d for decision", bid, did)})
		return
	}

	result := gin.H{"ballot": new_ballot}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "ballot_update.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// HBallotDelete deletes a ballot from a decision
func HBallotDelete(c *gin.Context) {

	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	b := &Ballot{Ballot_ID: bid, Decision_ID: did}
	err = b.Destroy()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "ballot_deleted.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
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
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Unable to find ballot %v for decision %v", bid, did)})
		return
	}

	result := gin.H{"ballot": ballot}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "ballot_info.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}

}

// Destroy removes a ballot from the database
// it also removes the dependencies of a ballots
// such as votes
func (b *Ballot) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM ballot WHERE ballot_id=$1", b.Ballot_ID)
	if err != nil {
		return fmt.Errorf("Unable to delete ballot %#v from database", b)
	}

	// Remove votes beloning to this ballot
	var votes []Vote
	_, err = dbmap.Select(&votes, "SELECT * FROM vote WHERE ballot_id=$1", b.Ballot_ID)
	if err != nil {
		return fmt.Errorf("Unable to find votes for ballot %#v", b)
	}

	for _, v := range votes {
		err = v.Destroy()
		if err != nil {
			return err
		}
	}
	return nil
}

// HBallotLogin is used to login users to their ballot
// it only sets the cookie for the user
// eg : A ballot is created by a facilitator
// The ballot has a unique secret pkbdf2 hashed
// A link is sent to the user in the form
// /decision/123/ballot/222/login/:secret
// user click on the link and is redirected to .... some url
func HBallotLogin(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "unable to parse decision id"})
		return
	}
	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "unable to parse ballot id"})
		return
	}
	secret := c.Param("secret")

	// find the ballot
	var ballot Ballot
	err = dbmap.SelectOne(&ballot, "SELECT * FROM ballot where ballot_id=$1 and decision_id=$2", bid, did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Unable to find ballot %v for decision %v", bid, did)})
		return
	}

	// TODO : remove this if above todo is done
	if ballot.Secret != secret {
		c.JSON(http.StatusNotFound, gin.H{"error": "Secret does not belong to this ballot"})
		return
	}

	// set the cookies
	ballot_id_str := strconv.Itoa(ballot.Ballot_ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "unable to parse ballot_id"})
		return
	}
	decision_id_str := strconv.Itoa(ballot.Decision_ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "unable to parse ballot_id"})
		return
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	bcookie := http.Cookie{
		Name:    "ballot_id",
		Value:   ballot_id_str,
		Path:    "/",
		Expires: expiration}
	dcookie := http.Cookie{
		Name:    "decision_id",
		Value:   decision_id_str,
		Path:    "/",
		Expires: expiration}
	http.SetCookie(c.Writer, &bcookie)
	http.SetCookie(c.Writer, &dcookie)
	// TODO : Change url
	c.Redirect(http.StatusSeeOther, "http://localhost/ballot.html")
}

func HBallotWhoami(c *gin.Context) {
	bcookie, err := c.Request.Cookie("ballot_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "unable to find ballot cookie"})
		return
	}
	dcookie, err := c.Request.Cookie("decision_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "unable to find decision cookie"})
		return
	}

	dval, err := strconv.Atoi(dcookie.Value)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to parse cookie decision id"})
		return
	}
	bval, err := strconv.Atoi(bcookie.Value)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to parse cookie ballot id"})
		return
	}

	result := gin.H{"ballot": Ballot{Ballot_ID: bval, Decision_ID: dval}}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "ballot_whoami.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}

}

// Save inserts a ballot into the database
func (b *Ballot) Save() error {

	b.SetupSecret()

	if err := dbmap.Insert(b); err != nil {
		return fmt.Errorf("Unable to insert ballot %#v to database", b)
	}
	return nil
}

// SetupSecret sets up the secret for a ballot
// it's unique :<
func (b *Ballot) SetupSecret() {
	b.Secret = HashPassword(fmt.Sprintf("b_%d_d_%d",
		b.Ballot_ID,
		b.Decision_ID))
}
