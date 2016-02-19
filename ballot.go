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

type BallotAllInfo struct {
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	URL_Decision string   `json:"url"`
	Votes        []Vote   `json:"votes"`
	Ratings      []Rating `json:"rating"`
}

// HBallotCreate create a ballot that belongs
// to a decision
func HBallotCreate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid decision_id"})
		return
	}

	var b Ballot
	if err := c.Bind(&b); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid ballot object"})
		return
	}
	b.Decision_ID = did // inherited

	err = b.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"ballot": b}
	c.Writer.Header().Set("Location", fmt.Sprintf("/ballot/%d", b.Ballot_ID))
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "ballot_create.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}

	// Send invitation
	title := fmt.Sprintf("%s's ballot", b.Name)
	body := fmt.Sprintf("Hello %s, you have been invited to participate in a decision at the following url : http://localhost:9999/decision/%d/ballot/%d/login/%s",
		b.Name, b.Decision_ID, b.Ballot_ID, b.Secret)

	// Non-blocking can improve error handling
	// but must block
	go Send(body, title, b.Email)
}

// HBallotInvite invites a specific ballot via email
// to participate in its decision
func HBallotInvite(c *gin.Context) {
	did := c.Param("decision_id")
	bid := c.Param("ballot_id")

	var b Ballot
	err := dbmap.SelectOne(&b,
		"SELECT * FROM ballot where ballot_id=$1 and decision_id=$2", bid, did)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to find ballot %v for decision %v", bid, did)})
		return
	}

	title := fmt.Sprintf("%s's ballot", b.Name)
	body := fmt.Sprintf("Hello %s, you have been invited to participate in a decision at the following url : http://localhost:9999/decision/%d/ballot/%d/login/%s",
		b.Name, b.Decision_ID, b.Ballot_ID, b.Secret)

	err = Send(body, title, b.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": "invited"})
	}
}

// HBallotUpdate updates a ballot
func HBallotUpdate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var b Ballot
	err = dbmap.SelectOne(&b, "SELECT * FROM ballot WHERE decision_id=$1 and ballot_id=$2", did, bid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("ballot %d for decision %d not found", bid, did)})
		return
	}

	var json Ballot
	err = c.Bind(&json)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": "Unable to parse ballot object"})
		return
	}

  var secret string = HashPassword(fmt.Sprintf("b_%d_d_%d", bid, did))
  
	new_ballot := Ballot{
		Ballot_ID:   bid,
		Decision_ID: did,
		Secret:      secret,
		Name:        json.Name,
		Email:       json.Email,
	}
	_, err = dbmap.Update(&new_ballot)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update ballot %d for decision %d", bid, did)})
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
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	b := &Ballot{Ballot_ID: bid, Decision_ID: did}
	err = b.Destroy()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Unable to find ballot %v for decision %v", bid, did)})
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
func HBallotLogin(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unable to parse decision id"})
		return
	}
	bid, err := strconv.Atoi(c.Param("ballot_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unable to parse ballot id"})
		return
	}
	secret := c.Param("secret")

	// Find the ballot
	var ballot Ballot
	err = dbmap.SelectOne(&ballot, "SELECT * FROM ballot where ballot_id=$1 and decision_id=$2", bid, did)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Unable to find ballot %v for decision %v", bid, did)})
		return
	}

	if ballot.Secret != secret {
		c.JSON(http.StatusForbidden, gin.H{"error": "Secret does not belong to this ballot"})
		return
	}

	// Set the cookies
	ballot_id_str := strconv.Itoa(ballot.Ballot_ID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unable to parse ballot_id"})
		return
	}
	decision_id_str := strconv.Itoa(ballot.Decision_ID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unable to parse ballot_id"})
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

	c.Redirect(http.StatusSeeOther, "http://localhost:9999/ballot.html")
}

func HBallotWhoami(c *gin.Context) {
	bcookie, err := c.Request.Cookie("ballot_id")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unable to find ballot cookie"})
		return
	}
	dcookie, err := c.Request.Cookie("decision_id")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unable to find decision cookie"})
		return
	}

	dval, err := strconv.Atoi(dcookie.Value)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unable to parse cookie decision id"})
		return
	}
	bval, err := strconv.Atoi(bcookie.Value)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unable to parse cookie ballot id"})
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

// HBallotBallotAllInfo show all of the information related to a ballot
func HBallotAllInfo(c *gin.Context) {
	did := c.Param("decision_id")
	bid := c.Param("ballot_id")

	var ballot Ballot
	err := dbmap.SelectOne(&ballot, "SELECT * FROM ballot where ballot_id=$1 and decision_id=$2", bid, did)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to find ballot %v for decision %v", bid, did)})
		return
	}

	// First get the ballot
	var ai BallotAllInfo
	ai.Name = ballot.Name
	ai.Email = ballot.Email
	ai.URL_Decision = fmt.Sprintf("/decision/%s", did)

	// Get the votes for this ballot
	_, err = dbmap.Select(&ai.Votes, "SELECT * FROM vote where ballot_id=$1", bid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to find votes for ballot %v", bid)})
		return
	}

	// Get the ratings for this ballot
	_, err = dbmap.Select(&ai.Ratings, "SELECT * FROM rating where ballot_id=$1", bid)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to find ratings for ballot %v", bid)})
		return
	}

	result := gin.H{"ballot": ai}

	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl",
			gin.H{"scriptname": "ballots_all_info.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// Save inserts a ballot into the database
func (b *Ballot) Save() error {

	// Check if decision exists or not
	var d Decision
	err := dbmap.SelectOne(&d, "select * from decision where decision_id=$1", b.Decision_ID)
	if err != nil {
		return fmt.Errorf("Decision %d does not exists, can't create ballot without a decision", b.Decision_ID)
	}

	if err = dbmap.Insert(b); err != nil {
		return fmt.Errorf("Unable to insert ballot %#v to database", b)
	}

	// Now recreate the hash correctly for
	// the current ballot after knowing
	// the id
	if err := dbmap.SelectOne(b, "select * from ballot where ballot_id=(select max(ballot_id) from ballot) and email=$1", b.Email); err != nil {
		return fmt.Errorf("Unable to set secret on ballot %#v to database", err)
	}

	b.Secret = HashPassword(fmt.Sprintf("b_%d_d_%d", b.Ballot_ID, b.Decision_ID))
	if _, err = dbmap.Update(b); err != nil {
		return fmt.Errorf("Unable to update ballot %#v to database", b)
	}

	return nil
}
