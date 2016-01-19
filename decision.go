package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Decision struct {
	Decision_ID            int    `binding:"required"`
	Person_ID              int    `binding:"required"`
	Name                   string `binding:"required"`
	Description            string `binding:"required"`
	Owner_ID               int    `binding:"required"`
	Stage                  int    `binding:"required"`
	Criterion_Vote_Style   string `binding:"required"`
	Alternative_Vote_Style string `binding:"required"`
	Client_Settings        string `binding:"required"`
}

func HDecisionBallots(c *gin.Context) {
	did := c.Param("decision_id")
	rows, err := database.DB.Query("SELECT * FROM ballot WHERE decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	defer rows.Close()

	var ballots []Ballot
	for rows.Next() {
		var b Ballot
		err := rows.Scan(&b.Ballot_ID, &b.Decision_ID, &b.Secret, &b.Name, &b.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		ballots = append(ballots, b)
	}

	if len(ballots) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": ballots})
}

func HDecisionsList(c *gin.Context) {
	rows, err := database.DB.Query("SELECT * FROM decision")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	defer rows.Close()

	var decisions []Decision
	for rows.Next() {
		var d Decision
		err = rows.Scan(&d.Decision_ID, &d.Person_ID, &d.Name,
			&d.Description, &d.Owner_ID, &d.Stage, &d.Criterion_Vote_Style,
			&d.Alternative_Vote_Style, &d.Client_Settings)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		decisions = append(decisions, d)
	}

	if len(decisions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": decisions})
}

func HDecisionInfo(c *gin.Context) {
	id := c.Param("decision_id")
	rows, err := database.DB.Query("SELECT * FROM decision where decision_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	defer rows.Close()

	var d Decision
	for rows.Next() {
		err = rows.Scan(&d.Decision_ID, &d.Person_ID, &d.Name,
			&d.Description, &d.Owner_ID, &d.Stage, &d.Criterion_Vote_Style,
			&d.Alternative_Vote_Style, &d.Client_Settings)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"result": d})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func HDecisionCreate(c *gin.Context) {
	var decision Decision
	err := c.Bind(&decision)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = decision.CreateDecision()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": decision})
}

func HDecisionDelete(c *gin.Context) {
	id := c.Param("decision_id")
	_, err := database.DB.Exec("DELETE FROM decision WHERE decision_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

func (d *Decision) CreateDecision() error {
	_, err := database.DB.Exec(
		"INSERT INTO decision VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
		d.Decision_ID, d.Person_ID, d.Name, d.Description, d.Owner_ID,
		d.Stage, d.Criterion_Vote_Style, d.Alternative_Vote_Style, d.Client_Settings)
	if err != nil {
		return err
	}
	return nil
}
