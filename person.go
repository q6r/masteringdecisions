package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Person represent a person in the database
type Person struct {
	Person_ID  int    `db:"person_id" json:"person_id"`
	Email      string `db:"email" json:"email" binding:"required"`
	PW_hash    string `db:"pw_hash" json:"pw_hash" binding:"required"`
	Name_First string `db:"name_first" json:"name_first" binding:"required"`
	Name_Last  string `db:"name_last" json:"name_last" binding:"required"`
}

// HPersonsList returns a list of persons as json object
func HPersonsList(c *gin.Context) {
	var persons []Person
	_, err := dbmap.Select(&persons, "select * from person order by person_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, persons)
}

// HPersonCreate creates a person in the database
func HPersonCreate(c *gin.Context) {
	var person Person
	err := c.Bind(&person)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid person object"})
		return
	}

	err = person.Save()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

// HPersonDelete deletes a person from the database
func HPersonDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("person_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	p := &Person{Person_ID: id}
	err = p.Destroy()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

// HPersonInfo return information for a person
func HPersonInfo(c *gin.Context) {
	id := c.Param("person_id")

	var person Person
	err := dbmap.SelectOne(&person, "select * from person where person_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, person)
}

// HPersonDecisions returns the decisions owned by that person
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

// Destroy a person from the database and remove its dependencies
func (p *Person) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM person WHERE person_id=$1", p.Person_ID)
	if err != nil {
		return err
	}

	// Remove the person's decisions
	var decisions []Decision
	_, err = dbmap.Select(&decisions, "SELECT * from decision WHERE person_id=$1", p.Person_ID)
	if err != nil {
		return err
	}

	for _, d := range decisions {
		err := d.Destroy()
		if err != nil {
			return err
		}
	}

	return nil
}

// Save a person to database
func (p *Person) Save() error {
	if err := dbmap.Insert(p); err != nil {
		return err
	}
	return nil
}
