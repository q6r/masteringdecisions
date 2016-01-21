// TODO : Problem with binding it caches default value as 'unset'
// eg:sending integer 0 is considered invalid
// eg:sending an empty string is considered invalid

// TODO : Write more test cover at least 80% of code
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
)

// TODO : Write tests

var dbmap *gorp.DbMap

func main() {
	dbmap = InitDatabase()
	defer dbmap.Db.Close()

	routes := gin.Default()

	// Person
	/*
		1. Create a person
		2. List all persons
		3. List all persons decisions
		4. Delete a person
		5. TODO : Change a person information (PUT/PATCH)
	*/
	routes.POST("/person", HPersonCreate)
	routes.GET("/persons", HPersonsList)
	routes.GET("/person/:person_id/info", HPersonInfo)
	routes.GET("/person/:person_id/decisions", HPersonDecisions)
	routes.DELETE("/person/:person_id", HPersonDelete)

	// Decision
	/*
		1. Create a decision
		2. List all decisions
		3. Delete a decision
		4. List decision ballots
		5. TODO : Change a decision information
	*/
	routes.POST("/decision", HDecisionCreate)
	routes.GET("/decisions", HDecisionsList)
	routes.GET("/decision/:decision_id/info", HDecisionInfo)
	routes.GET("/decision/:decision_id/ballots", HDecisionBallots)
	routes.DELETE("/decision/:decision_id", HDecisionDelete)

	// Ballot
	/*
		1. Create a ballot
		2. List all ballots
		3. Show ballot information
		4. Delete a ballot
		5. TODO : Change a ballot feature
	*/
	routes.POST("/ballot", HBallotCreate)
	routes.GET("/ballots", HBallotList)
	routes.GET("/ballot/:ballot_id", HBallotInfo)
	routes.POST("/ballot/:ballot_id", HBallotVote)
	routes.DELETE("/ballot/:ballot_id", HBallotDelete)

	// Criterion
	/*
		1. Create a criterion
		2. List all criterions
		3. Show criterion information
		4. Delete a criterion
		5. TODO : Change a criterion information
	*/
	routes.POST("/criterion", HCriterionCreate)
	routes.GET("/criterions", HCriterionList)
	routes.GET("/criterion/:criterion_id/info", HCriterionInfo)
	routes.DELETE("/criterion/:criterion_id", HCriterionDelete)

	routes.Run(":9999")
}
