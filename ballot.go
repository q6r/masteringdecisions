package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Ballot struct {
	Ballot_ID   int    `binding:"required"`
	Decision_ID int    `binding:"required"`
	Secret      int    `binding:"required"`
	Name        string `binding:"required"`
	Email       string `binding:"required"`
}

func HBallotCreate(c *gin.Context) {
	var b Ballot
	if err := c.BindJSON(&b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	log.Printf("%#v\n", b)

	err := b.CreateBallot()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": b})
}

func HBallotList(c *gin.Context) {
	rows, err := database.DB.Query("SELECT * FROM ballot")
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

func HBallotInfo(c *gin.Context) {
	bid := c.Param("ballot_id")
	row, err := database.DB.Query("SELECT * FROM ballot where ballot_id=$1", bid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	defer row.Close()

	var b Ballot
	for row.Next() {
		err := row.Scan(&b.Ballot_ID, &b.Decision_ID, &b.Secret, &b.Name, &b.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"result": b})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func (b *Ballot) CreateBallot() error {
	_, err := database.DB.Exec("INSERT INTO ballot VALUES($1,$2,$3,$4,$5)",
		b.Ballot_ID, b.Decision_ID, b.Secret, b.Name, b.Email)
	if err != nil {
		return err
	}
	return nil
}
