package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Decision represent a decision owned by Person_ID
type Decision struct {
	DecisionID             int    `db:"decision_id" json:"decision_id"`
	PersonID               int    `db:"person_id" json:"person_id" binding:"required"`
	Name                   string `db:"name" json:"name" binding:"required"`
	Description            string `db:"description" json:"description" binding:"required"`
	Stage                  int    `db:"stage" json:"stage" binding:"required"`
	CriterionVoteStyle     string `db:"criterion_vote_style" json:"criterion_vote_style" binding:"required"`
	AlternativeVoteStyle   string `db:"alternative_vote_style" json:"alternative_vote_style" binding:"required"`
	ClientSettings         string `db:"client_settings" json:"client_settings"`
	DisplayName            string `db:"display_name" json:"display_name"`
	CriteriaInstruction    string `db:"criteria_instruction" json:"criteria_instruction"`
	AlternativeInstruction string `db:"alternative_instruction" json:"alternative_instruction"`
	Image                  string `db:"image" json:"image"`
}

// HDecisionBallotsList returns a list of ballots beloning
// to a decision, show all their information Using
// an array of BallotAllInfo
func HDecisionBallotsList(c *gin.Context) {
	did := c.Param("decision_id")
	var ballots []Ballot
	_, err := dbmap.Select(&ballots, "SELECT * FROM ballot WHERE decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Unable to find ballots for decision id %v", did)})
		return
	}

	var ais []BallotAllInfo

	for _, b := range ballots {
		var ai BallotAllInfo
		ai.Name = b.Name
		ai.Email = b.Email
    ai.Sent = b.Sent
		ai.URLDecision = fmt.Sprintf("/decision/%s/ballot/%d", did, b.BallotID)
		// Get the votes for this ballot
		_, err = dbmap.Select(&ai.Votes, "SELECT * FROM vote where ballot_id=$1", b.BallotID)
		if err != nil {
			c.JSON(http.StatusForbidden,
				gin.H{"error": fmt.Sprintf("Unable to find votes for ballot %v", b.BallotID)})
			return
		}
		// Get the ratings for this ballot
		_, err = dbmap.Select(&ai.Ratings, "SELECT * FROM rating where ballot_id=$1", b.BallotID)
		if err != nil {
			c.JSON(http.StatusForbidden,
				gin.H{"error": fmt.Sprintf("Unable to find votes for ballot %v", b.BallotID)})
			return
		}
		ais = append(ais, ai)
	}

	result := gin.H{"ballots": ais}
	ServeResult(c, "decision_ballots.js", result)
}

// HDecisionAlternativesList returns a list of alternatives beloning
// to a decision
func HDecisionAlternativesList(c *gin.Context) {
	did := c.Param("decision_id")
	var alts []Alternative
	_, err := dbmap.Select(&alts, "SELECT * FROM alternative WHERE decision_id=$1 ORDER BY \"order\", \"name\" ASC", did)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to find alternatives for decision id %s", did)})
		return
	}

	result := gin.H{"alternatives": alts}
	ServeResult(c, "decision_alternatives.js", result)
}

// HDecisionCriterionsList returns a list of criterions beloning
// to a decision
func HDecisionCriterionsList(c *gin.Context) {
	did := c.Param("decision_id")
	var cris []Criterion
	_, err := dbmap.Select(&cris, "SELECT * FROM criterion WHERE decision_id=$1 ORDER BY \"order\", \"name\" ASC", did)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to find criterion for decision %s", did)})
		return
	}

	result := gin.H{"criterions": cris}
	ServeResult(c, "decision_criterions.js", result)
}

// HDecisionDuplicate duplicates a decision by cloning
// its information, criterions, and alternatives
func HDecisionDuplicate(c *gin.Context) {
	did := c.Param("decision_id")

	// Get the decision to duplicate
	// and its alternatives and criterions
	dobj, err := dbmap.Get(Decision{}, did)
	if err != nil || dobj == nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to find decision %s", did)})
		return
	}
	dsrc := dobj.(*Decision)
	var cris []Criterion
	_, _ = dbmap.Select(&cris, "select * from criterion where decision_id=$1", dsrc.DecisionID)
	var alts []Alternative
	_, _ = dbmap.Select(&alts, "select * from alternative where decision_id=$1", dsrc.DecisionID)

	// This is auto incr we need to inset and
	// determine its new value
	dsrc.DecisionID = 0
	if err := dbmap.Insert(dsrc); err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": "Unable to insert duplicated decision"})
		return
	}
	// Get the new auto incr id
	if err := dbmap.SelectOne(dsrc,
		"select * from decision where decision_id=(select max(decision_id) from decision) and person_id=$1 and name=$2",
		dsrc.PersonID, dsrc.Name); err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": "Unable to get duplicated decision"})
		return
	}

	// Now change the criterions and alternatives
	// decision ownership
	for _, cri := range cris {
		cri.DecisionID = dsrc.DecisionID
		if err := cri.Save(); err != nil {
			c.JSON(http.StatusForbidden,
				gin.H{"error": "Unable to save criterion of duplicated decision"})
			return
		}
	}
	for _, alt := range alts {
		alt.DecisionID = dsrc.DecisionID
		if err := alt.Save(); err != nil {
			c.JSON(http.StatusForbidden,
				gin.H{"error": "Unable to save alternative of duplicated decision"})
			return
		}
	}

	result := gin.H{"decision": dsrc}
	ServeResult(c, "decision_duplicate.js", result)
}

// HDecisionsList returns a list of all decision defined
// in the database their name and url only
func HDecisionsList(c *gin.Context) {
	var decisions []Decision
	_, err := dbmap.Select(&decisions, "SELECT * FROM decision")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unable to find decisions in database"})
		return
	}

	type Link struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	var links []Link
	for _, d := range decisions {
		l := Link{Name: d.Name, URL: fmt.Sprintf("/decision/%d", d.DecisionID)}
		links = append(links, l)
	}

	result := gin.H{"decisions": links}
	ServeResult(c, "decisions_list.js", result)
}

// HDecisionInfo returns a decision information
// a decision object not it's stats
func HDecisionInfo(c *gin.Context) {
	did := c.Param("decision_id")
	var decision Decision
	err := dbmap.SelectOne(&decision, "SELECT * FROM decision where decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Unable to find decisions with id %v", did)})
		return
	}

	result := gin.H{"decision": decision}
	ServeResult(c, "decision_info.js", result)
}

// HDecisionUpdate updates a decision
func HDecisionUpdate(c *gin.Context) {
	did, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var d Decision
	err = dbmap.SelectOne(&d, "SELECT * FROM decision WHERE decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("decision %d not found", did)})
		return
	}

	var json Decision
	err = c.Bind(&json)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": "Unable to parse decision object"})
		return
	}

	newDecision := Decision{
		DecisionID:             did,
		PersonID:               json.PersonID,
		Name:                   json.Name,
		Description:            json.Description,
		Stage:                  json.Stage,
		CriterionVoteStyle:     json.CriterionVoteStyle,
		AlternativeVoteStyle:   json.AlternativeVoteStyle,
		ClientSettings:         json.ClientSettings,
		DisplayName:            json.DisplayName,
		CriteriaInstruction:    json.CriteriaInstruction,
		AlternativeInstruction: json.AlternativeInstruction,
		Image: json.Image,
	}

	_, err = dbmap.Update(&newDecision)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update decision %d", did)})
		return
	}

	result := gin.H{"decision": newDecision}
	ServeResult(c, "decision_update.js", result)
}

// HDecisionCreate creates a decision beloning to a specific
// person
func HDecisionCreate(c *gin.Context) {

	var decision Decision
	err := c.Bind(&decision)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid decision object"})
		return
	}

	err = decision.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"decision": decision}
	c.Writer.Header().Set("Location", fmt.Sprintf("/decision/%d", decision.DecisionID))
	ServeResult(c, "decision_create.js", result)
}

// HDecisionDelete deletes a decision from database
func HDecisionDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("decision_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	d := &Decision{DecisionID: id}
	err = d.Destroy()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	ServeResult(c, "decision_deleted.js", result)
}

// Destroy a decision from the database
// and remove it's dependencies such as ballots
// when destroying ballots they'll destroy their votes..etc
func (d *Decision) Destroy() error {
	if _, err := dbmap.Delete(d); err != nil {
		return fmt.Errorf("Unable to delete decision %#v from database", d)
	}

	// Remove the ballots of this decision
	// removes the votes
	var ballots []Ballot
	_, _ = dbmap.Select(&ballots, "SELECT * FROM ballot WHERE decision_id=$1", d.DecisionID)
	for _, b := range ballots {
		err := b.Destroy()
		if err != nil {
			return err
		}
	}

	// Remove criterions
	// Does not remove anything..
	var cris []Criterion
	_, _ = dbmap.Select(&cris, "select * from criterion where decision_id=$1", d.DecisionID)
	for _, cri := range cris {
		err := cri.Destroy()
		if err != nil {
			return err
		}
	}

	// Removing the alternatives remove the votes related to it
	var alts []Alternative
	_, _ = dbmap.Select(&alts, "select * from alternative where decision_id=$1", d.DecisionID)
	for _, alt := range alts {
		err := alt.Destroy()
		if err != nil {
			return err
		}
	}

	return nil
}

// Save saves decision in the database
// Restriction : Decision can't be created without an existing owner
// Restriction : Decision can't be owned by two different people
func (d *Decision) Save() error {

	// See if there's a person that this decision belongs to
	// otherwise we quit
	var p Person
	err := dbmap.SelectOne(&p, "SELECT * FROM person WHERE person_id=$1", d.PersonID)
	if err != nil {
		return fmt.Errorf("person %d does not exist, can't create a decision without an owner", d.PersonID)
	}

	// Check ownership of decisions
	var ds []Decision
	_, _ = dbmap.Select(&ds, "select * from decision where decision_id=$1", d.DecisionID)
	for _, i := range ds {
		if i.PersonID != d.PersonID {
			return fmt.Errorf("decision %d already owned by person %d", d.DecisionID, i.PersonID)
		}
	}

	if err = dbmap.Insert(d); err != nil {
		return fmt.Errorf("Unable to insert decision %#v to database", d)
	}
	return nil
}
