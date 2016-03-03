package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var mutex sync.Mutex

// Ballot represent a ballot that belong
// in a decision
type Ballot struct {
	BallotID   int    `db:"ballot_id" json:"ballot_id"`
	DecisionID int    `db:"decision_id" json:"decision_id"`
	Secret     string `db:"secret" json:"secret"`
	Name       string `db:"name" json:"name" binding:"required"`
	Email      string `db:"email" json:"email" binding:"required"`
	Sent       bool   `db:"sent" json:"sent"`
}

// BallotAllInfo is a struct to hold
// all necessary information of a ballot
type BallotAllInfo struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	URLDecision string   `json:"url"`
	Votes       []Vote   `json:"votes"`
	Ratings     []Rating `json:"rating"`
	Sent        bool     `json:"sent"`
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
	b.DecisionID = did // inherited

	mutex.Lock()
	defer mutex.Unlock()
	err = b.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"ballot": b}
	c.Writer.Header().Set("Location", fmt.Sprintf("/ballot/%d", b.BallotID))
	ServeResult(c, "ballot_create.js", result)

	// Send email to ballot, and set confirmation
	// flag on success
	title, body := GenerateInviteTemplate(b)
	go func(b *Ballot) {
		if b == nil {
			return
		}
		if err := Send(body, title, b.Email); err != nil {
			b.Sent = false
			return
		}
		b.Sent = true
		if _, err := dbmap.Update(b); err != nil {
			return
		}
	}(&b)
}

// HBallotCreateSilent create a ballot without sending
// it an invitation
func HBallotCreateSilent(c *gin.Context) {
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
	b.DecisionID = did // inherited

	mutex.Lock()
	defer mutex.Unlock()
	err = b.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"ballot": b}
	c.Writer.Header().Set("Location", fmt.Sprintf("/ballot/%d", b.BallotID))
	ServeResult(c, "ballot_create.js", result)
}

// GenerateInviteTemplate generate the template for email
// invitations
func GenerateInviteTemplate(b Ballot) (title string, body string) {
	title = fmt.Sprintf("%s's ballot", b.Name)
	body = fmt.Sprintf("<html><body>Hello %s, you have been invited to participate in a decision <a href=\"http://localhost:9999/decision/%d/ballot/%d/login/%s\">click here to vote</a>.</body></html>",
		b.Name, b.DecisionID, b.BallotID, b.Secret)
	return title, body
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

	c.JSON(http.StatusOK, gin.H{"result": "invited"})

	// Send email to ballot, and set confirmation
	// flag on success
	title, body := GenerateInviteTemplate(b)
	go func(b *Ballot) {
		if b == nil {
			return
		}
		if err := Send(body, title, b.Email); err != nil {
			b.Sent = false
			return
		}
		b.Sent = true
		if _, err := dbmap.Update(b); err != nil {
			return
		}
	}(&b)
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

	secret := HashPassword(fmt.Sprintf("b_%d_d_%d", bid, did))

	newBallot := Ballot{
		BallotID:   bid,
		DecisionID: did,
		Secret:     secret,
		Name:       json.Name,
		Email:      json.Email,
		Sent:       b.Sent,
	}
	_, err = dbmap.Update(&newBallot)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update ballot %d for decision %d", bid, did)})
		return
	}

	result := gin.H{"ballot": newBallot}
	ServeResult(c, "ballot_update.js", result)
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

	b := &Ballot{BallotID: bid, DecisionID: did}
	err = b.Destroy()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	ServeResult(c, "ballot_deleted.js", result)
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
	ServeResult(c, "ballot_info.js", result)
}

// Destroy removes a ballot from the database
// it also removes the dependencies of a ballots
// such as votes
func (b *Ballot) Destroy() error {
	if _, err := dbmap.Delete(b); err != nil {
		return fmt.Errorf("Unable to delete ballot %#v from database", b)
	}

	// Remove votes beloning to this ballot
	var votes []Vote
	_, _ = dbmap.Select(&votes, "SELECT * FROM vote WHERE ballot_id=$1", b.BallotID)
	for _, v := range votes {
		if err := v.Destroy(); err != nil {
			return err
		}
	}

	// Remove ratings belonging to this ballot
	var ratings []Rating
	_, _ = dbmap.Select(&ratings, "SELECT * FROM rating WHERE ballot_id=$1", b.BallotID)
	for _, r := range ratings {
		if err := r.Destroy(); err != nil {
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
	ballotIDStr := strconv.Itoa(ballot.BallotID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unable to parse ballot_id"})
		return
	}
	decisionIDStr := strconv.Itoa(ballot.DecisionID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unable to parse ballot_id"})
		return
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	bcookie := http.Cookie{
		Name:    "ballot_id",
		Value:   ballotIDStr,
		Path:    "/",
		Expires: expiration}
	dcookie := http.Cookie{
		Name:    "decision_id",
		Value:   decisionIDStr,
		Path:    "/",
		Expires: expiration}
	http.SetCookie(c.Writer, &bcookie)
	http.SetCookie(c.Writer, &dcookie)

	c.Redirect(http.StatusSeeOther, "http://localhost:9999/ballot.html")
}

// HBallotWhoami returns the current
// logged in ballot
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

	result := gin.H{"ballot": Ballot{BallotID: bval, DecisionID: dval}}
	ServeResult(c, "ballot_whoami.js", result)
}

// HBallotAllInfo show all of the information related to a ballot
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
	ai.Sent = ballot.Sent
	ai.URLDecision = fmt.Sprintf("/decision/%s", did)

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
	ServeResult(c, "ballots_all_info.js", result)
}

// Save inserts a ballot into the database
func (b *Ballot) Save() error {

	// Check if decision exists or not
	var d Decision
	err := dbmap.SelectOne(&d, "select * from decision where decision_id=$1", b.DecisionID)
	if err != nil {
		return fmt.Errorf("Decision %d does not exists, can't create ballot without a decision", b.DecisionID)
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

	b.Secret = HashPassword(fmt.Sprintf("b_%d_d_%d", b.BallotID, b.DecisionID))
	if _, err = dbmap.Update(b); err != nil {
		return fmt.Errorf("Unable to update ballot %#v to database", b)
	}

	return nil
}
