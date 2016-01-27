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

	result := gin.H{"persons": persons}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "persons_list.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// HPersonUpdate updates a perosn
func HPersonUpdate(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("person_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var p Person
	err = dbmap.SelectOne(&p, "SELECT * FROM person WHERE person_id=$1", pid)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": fmt.Sprintf("person %d not found", pid)})
		return
	}

	var json Person
	err = c.Bind(&json)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": "Unable to parse person object"})
		return
	}

	new_hash := HashPassword(json.PW_hash)
	new_person := Person{
		Person_ID:  pid,
		Email:      json.Email,
		PW_hash:    new_hash,
		Name_First: json.Name_First,
		Name_Last:  json.Name_Last,
	}
	_, err = dbmap.Update(&new_person)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": fmt.Sprintf("Unable to update person %d", pid)})
		return
	}

	new_person.PW_hash = "<hidden>"
	result := gin.H{"person": new_person}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "person_update.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
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
	result := gin.H{"person": person}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "person_create.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
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

	result := gin.H{"result": "deleted"}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "person_deleted.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
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

	result := gin.H{"person": person}
	c.JSON(http.StatusOK, result)
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

	result := gin.H{"decisions": decisions}
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "htmlwrapper.tmpl", gin.H{"scriptname": "person_decisions.js", "body": result})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// Destroy a person from the database and remove its dependencies
func (p *Person) Destroy() error {
	rd, err := dbmap.Delete(p)
	if err != nil {
		return fmt.Errorf("Unable to delete %#v from database : %#v", p, err)
	}
	if rd == 0 {
		return fmt.Errorf("Nothing deleted")
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
