package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/astaxie/beego/config"
	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/rageix/ginAuth"
)

var dbmap *gorp.DbMap

func main() {

	conf, err := config.NewConfig("ini", "config.conf")
	if err != nil {
		log.Fatalln(err)
	}

	dbmap = InitDatabase(conf)
	defer func() {
		err := dbmap.Db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	routes := gin.Default()

	// Debug Mode only clean route
	// used for testing purposes
	if gin.Mode() == "debug" {
		routes.GET("/clean", func(c *gin.Context) {
			err := dbmap.TruncateTables()
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err})
				return
			}
			c.JSON(http.StatusOK, gin.H{"result": "cleaned"})
		})
	}

	// enable logging into a file
	err = enableLoggerFile(routes, conf.String("logfile"))
	if err != nil {
		log.Fatalln(err)
	}

	// Templates
	////////////////
	routes.LoadHTMLGlob("templates/*")
	routes.Static("static/", "static/")

	// Roots of pages for the frontend
	/////////////////////////
	routes.GET("/", HRootHome)
	routes.GET("/login.html", HRootLogin)
	routes.GET("/logout.html", HRootLogout)
	routes.GET("/ballot.html", HRootBallot)
	routes.GET("/decision/:decision_id", HRootDecision)
	routes.GET("/results/:decision_id", HRootResults)

	// Person
	////////////////
	routes.POST("/person", ginAuth.Use, AuthAsAdmin, HPersonCreate)
	routes.GET("/persons", HPersonsList)
	routes.GET("/person/:person_id/info", HPersonInfo)
	routes.GET("/person/:person_id/decisions", HPersonDecisions)
	routes.DELETE("/person/:person_id", ginAuth.Use, AuthAsAdmin, HPersonDelete)
	routes.PUT("/person/:person_id", ginAuth.Use, AuthAsAll, HPersonUpdate)

	// Decision
	////////////////

	// decision homes
	routes.POST("/decision", ginAuth.Use, AuthAsAll, HDecisionCreate)
	routes.GET("/decision/:decision_id/duplicate", ginAuth.Use, AuthAsAll, HDecisionDuplicate)
	routes.GET("/decisions", HDecisionsList)
	routes.GET("/decision/:decision_id/info", HDecisionInfo)
	routes.DELETE("/decision/:decision_id", ginAuth.Use, AuthAsAll, HDecisionDelete)
	routes.PUT("/decision/:decision_id", ginAuth.Use, AuthAsAll, HDecisionUpdate)

	// decision's alternatives
	routes.POST("/decision/:decision_id/alternative", ginAuth.Use, AuthAsAll, HAlternativeCreate)
	routes.GET("/decision/:decision_id/alternatives", HDecisionAlternativesList)
	routes.GET("/decision/:decision_id/alternative/:alternative_id/info", HAlternativeInfo)
	routes.DELETE("/decision/:decision_id/alternative/:alternative_id", ginAuth.Use, AuthAsAll, HAlternativeDelete)
	routes.PUT("/decision/:decision_id/alternative/:alternative_id", ginAuth.Use, AuthAsAll, HAlternativeUpdate)

	// decision's ballots
	routes.GET("/decision/:decision_id/ballot/:ballot_id", HBallotAllInfo)
	routes.GET("/decision/:decision_id/ballot/:ballot_id/invite", HBallotInvite)
	routes.GET("/decision/:decision_id/ballots", HDecisionBallotsList)
	routes.POST("/decision/:decision_id/ballot", HBallotCreate)
	routes.POST("/decision/:decision_id/ballot_silent", HBallotCreateSilent)
	routes.GET("/decision/:decision_id/ballot/:ballot_id/info", HBallotInfo)
	routes.DELETE("/decision/:decision_id/ballot/:ballot_id", HBallotDelete)
	routes.PUT("/decision/:decision_id/ballot/:ballot_id", HBallotUpdate)

	routes.GET("/decision/:decision_id/ballot/:ballot_id/login/:secret", HBallotLogin)
	routes.GET("/ballot_whoami", HBallotWhoami)

	// decision's ballot's votes
	routes.GET(
		"/decision/:decision_id/ballot/:ballot_id/alternative/:alternative_id/criterion/:criterion_id/vote/:weight", HVoteCreate)
	routes.GET(
		"/decision/:decision_id/ballot/:ballot_id/votes", HVotesBallotList)
	routes.DELETE(
		"/decision/:decision_id/ballot/:ballot_id/alternative/:alternative_id/criterion/:criterion_id/vote",
		HVoteDelete)
	routes.PUT(
		"/decision/:decision_id/ballot/:ballot_id/alternative/:alternative_id/criterion/:criterion_id/vote/:weight",
		HVoteUpdate)

	// decision's ballot's rating criterion
	routes.GET(
		"/decision/:decision_id/ballot/:ballot_id/criterion/:criterion_id/vote/:rating", HRatingCreate)
	routes.GET(
		"/decision/:decision_id/criterion/:criterion_id/votes", HRatingBallots)
	routes.DELETE(
		"/decision/:decision_id/ballot/:ballot_id/criterion/:criterion_id/vote",
		HRatingDelete)

	routes.PUT(
		"/decision/:decision_id/ballot/:ballot_id/criterion/:criterion_id/vote/:rating",
		HRatingUpdate)

	// decision's criterions
	routes.GET("/decision/:decision_id/criterions", HDecisionCriterionsList)
	routes.POST("/decision/:decision_id/criterion", HCriterionCreate)
	routes.GET("/decision/:decision_id/criterion/:criterion_id/info", HCriterionInfo)
	routes.DELETE("/decision/:decision_id/criterion/:criterion_id", HCriterionDelete)
	routes.PUT("/decision/:decision_id/criterion/:criterion_id", HCriterionUpdate)

	// Login/Logout
	routes.POST("/login", HAuthLogin)
	routes.GET("/logout", HAuthLogout)
	routes.GET("/whoami", ginAuth.Use, AuthAsAll, HAuthWhoAmI)

	// Setup the authentication
	ginAuth.ConfigPath = "./config.conf"
	ginAuth.Unauthorized = HAuthUnauthenticated
	ginAuth.Authorized = HAuthAuthenticated
	err = ginAuth.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	httpPort, err := conf.Int("http_port")
	if err != nil {
		log.Fatalln(err)
	}

	err = routes.Run(fmt.Sprintf(":%d", httpPort))
	if err != nil {
		log.Fatalln(err)
	}
}

// enableLoggerFile enables logging to a file path
// empty logFilePath to disable logging to a file
func enableLoggerFile(routes *gin.Engine, logFilePath string) error {
	// Don't want logger
	if logFilePath == "" {
		return nil
	}

	logFile, err := os.OpenFile(logFilePath,
		os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}

	routes.Use(gin.LoggerWithWriter(logFile))

	return nil
}
