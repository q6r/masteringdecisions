package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Unable to find persons")})
		return
	}

	for i := range persons {
		persons[i].PW_hash = "<hidden>"
	}

	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "persons_list.js", "body": persons})
	} else {
		c.JSON(http.StatusOK, persons)
	}
}

// HPersonCreate creates a person in the database
func HPersonCreate(c *gin.Context) {
	var person Person
	err := c.Bind(&person)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid person object"})
		return
	}

	// Encrypt the plaintext before saving
	person.PW_hash = HashPassword(person.PW_hash)

	err = person.Save()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	person.PW_hash = "<hidden>"

	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "person_create.js", "body": person})
	} else {
		c.JSON(http.StatusOK, person)
	}
}

// HPersonDelete deletes a person from the database
func HPersonDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("person_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	p := &Person{Person_ID: id}
	err = p.Destroy()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "person_deleted.js", "body": gin.H{"result": "deleted"}})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": "deleted"})
	}
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

	person.PW_hash = "<hidden>"
	c.JSON(http.StatusOK, person)
}

// HPersonDecisions returns the decisions owned by that person
func HPersonDecisions(c *gin.Context) {
	id := c.Param("person_id")
	var decisions []Decision
	_, err := dbmap.Select(&decisions, "select * from decision where person_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "person_decisions.js", "body": decisions})
	} else {
		c.JSON(http.StatusOK, decisions)
	}
}

// Destroy a person from the database and remove its dependencies
func (p *Person) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM person WHERE person_id=$1", p.Person_ID)
	if err != nil {
		return fmt.Errorf("Unable to delete %#v from database : %#v", p, err)
	}

	// Remove the person's decisions
	var decisions []Decision
	_, err = dbmap.Select(&decisions, "SELECT * from decision WHERE person_id=$1", p.Person_ID)
	if err != nil {
		return fmt.Errorf("Unable to find decisions for person %#v", p)
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

	// disallow duplicate emails
	n, err := dbmap.SelectInt("select count(*) from person where email=$1", p.Email)
	if err != nil {
		return fmt.Errorf("Unable to check person database for duplicate email")
	}
	if n != 0 {
		return fmt.Errorf("Email %s already exists", p.Email)
	}

	if err := dbmap.Insert(p); err != nil {
		return fmt.Errorf("Unable to insert person %d into database:%#v",
			p.Person_ID, err)
	}
	return nil
}
