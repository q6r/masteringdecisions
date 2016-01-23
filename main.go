package main

// TODO : Vote weight should not be more than criterion weight

// TODO : Problem with binding it caches default value as 'unset'
// eg:sending integer 0 is considered invalid
// eg:sending an empty string is considered invalid

// TODO : Make sure duplicates are not allowed in Save of all

// TODO : Write more test cover at least 80% of code

// TODO : In all implement Check that checks if one with a primary key
// exists, good for code refactoring and simplifying

// TODO : Person pwhash should get encrypted

// TODO : Person should not send the hash back <implement after authorization>

// TODO : Write tests

// TODO : Review and test Save/Destroy restriction and write tests if possible :)

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/itsjamie/gin-cors"
)

var dbmap *gorp.DbMap

func main() {
	dbmap = InitDatabase()
	defer dbmap.Db.Close()

	routes := gin.Default()

	// Middlewares
	/////////////////
	// Apply the middleware to the router (works with groups too)
	routes.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	// Debug routes
	/////////////////
	if gin.Mode() == "debug" {
		routes.GET("/clean", func(c *gin.Context) {
			err := dbmap.TruncateTables()
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err})
				return
			}
			c.JSON(http.StatusOK, gin.H{"result": "cleaned"})
		})
	}

	// Person
	////////////////
	routes.POST("/person", HPersonCreate)
	routes.GET("/persons", HPersonsList)
	routes.GET("/person/:person_id/info", HPersonInfo)
	routes.GET("/person/:person_id/decisions", HPersonDecisions)
	routes.DELETE("/person/:person_id", HPersonDelete)

	// Decision
	////////////////

	// decision homes
	routes.POST("/decision", HDecisionCreate)
	routes.GET("/decisions", HDecisionsList)
	routes.GET("/decision/:decision_id/info", HDecisionInfo)
	routes.GET("/decision/:decision_id/stats", HStats)
	routes.DELETE("/decision/:decision_id", HDecisionDelete)

	// decision's ballots
	routes.GET("/decision/:decision_id/ballots", HDecisionBallotsList)
	routes.POST("/decision/:decision_id/ballot", HBallotCreate)
	routes.GET("/decision/:decision_id/ballot/:ballot_id/info", HBallotInfo)
	routes.DELETE("/decision/:decision_id/ballot/:ballot_id", HBallotDelete)

	// decision's ballot's votes
	routes.GET(
		"/decision/:decision_id/ballot/:ballot_id/criterion/:criterion_id/vote/:weight", HVoteCreate)
	routes.GET(
		"/decision/:decision_id/ballot/:ballot_id/votes", HVotesBallotList)
	routes.DELETE(
		"/decision/:decision_id/ballot/:ballot_id/criterion/:criterion_id/vote",
		HVoteDelete)

	// decision's criterions
	routes.GET("/decision/:decision_id/criterions", HDecisionCriterionsList)
	routes.POST("/decision/:decision_id/criterion", HCriterionCreate)
	routes.GET("/decision/:decision_id/criterion/:criterion_id/info", HCriterionInfo)
	routes.DELETE("/decision/:decision_id/criterion/:criterion_id", HCriterionDelete)

	routes.Run(":9999")
}
