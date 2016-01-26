package main

// TODO : Start applying the authentication middlewares to
// required routes
// TODO : Write more test cover at least 80% of code
// TODO : Person should not send the hash back <implement after authorization>
// TODO : Write tests

import (
	"log"
	"net/http"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/itsjamie/gin-cors"
	"github.com/rageix/ginAuth"
)

var dbmap *gorp.DbMap

func main() {

	conf, err := config.NewConfig("ini", "config.conf")
	if err != nil {
		log.Fatalln(err)
	}

	dbmap = InitDatabase(conf)
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

	// Login/Logout
	routes.POST("/login", HAuthLogin)
	routes.GET("/logout", HAuthLogout)
	// TODO : an example that requires authenticated users only
	routes.GET("/checklogin", ginAuth.Use, AuthAsAll, HAuthCheck)

	// Setup the authentication
	ginAuth.ConfigPath = "./config.conf"
	ginAuth.Unauthorized = HAuthUnauthenticated
	ginAuth.Authorized = HAuthAuthenticated
	err = ginAuth.LoadConfig()
	if err != nil {
		log.Fatalln(err)
		return
	}

	routes.Run(":9999")
}
