package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Person_ID  int    `db:"person_id"`
	Email      string `db:"email" binding:"required"`
	PW_hash    string `db:"pw_hash" binding:"required"`
	Name_First string `db:"name_first" binding:"required"`
	Name_Last  string `db:"name_last" binding:"required"`
}

func HPersonsList(c *gin.Context) {
	var persons []Person
	_, err := dbmap.Select(&persons, "select * from person order by person_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, persons)
}

func HPersonCreate(c *gin.Context) {
	var person Person
	err := c.Bind(&person)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid person object"})
		return
	}

	err = person.CreatePerson()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, person)
}

func HPersonDelete(c *gin.Context) {
	id := c.Param("person_id")

	_, err := dbmap.Exec("DELETE FROM person WHERE person_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

func HPersonInfo(c *gin.Context) {
	id := c.Param("person_id")
	var person Person
	_, err := dbmap.Select(&person, "select * from person where person_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, person)
}

func HPersonDecisions(c *gin.Context) {
	id := c.Param("person_id")
	var decisions []Decision
	_, err := dbmap.Select(&decisions, "select * from decision where person_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, decisions)
}

func (p *Person) CreatePerson() error {
	if err := dbmap.Insert(p); err != nil {
		return err
	}
	return nil
}
