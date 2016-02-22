package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Person represent a person in the database
type Person struct {
	PersonID  int    `db:"person_id" json:"person_id"`
	Email     string `db:"email" json:"email" binding:"required"`
	PWHash    string `db:"pw_hash" json:"pw_hash" binding:"required"`
	NameFirst string `db:"name_first" json:"name_first" binding:"required"`
	NameLast  string `db:"name_last" json:"name_last" binding:"required"`
}

// HPersonsList returns a list of persons as json object
func HPersonsList(c *gin.Context) {
	var persons []Person
	_, err := dbmap.Select(&persons, "select * from person order by person_id")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Unable to find persons")})
		return
	}

	for i := range persons {
		persons[i].PWHash = "<hidden>"
	}

	result := gin.H{"persons": persons}
	ServeResult(c, "persons_list.js", result)
}

// HPersonUpdate updates a perosn
func HPersonUpdate(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("person_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	pobj, err := dbmap.Get(Person{}, pid)
	if err != nil || pobj == nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("person %d not found", pid)})
		return
	}
	p := pobj.(*Person)

	// PersonEasy represent a person in the database
	// with no password requirement
	type PersonEasy struct {
		PersonID  int    `db:"person_id" json:"person_id"`
		Email     string `db:"email" json:"email" binding:"required"`
		PWHash    string `db:"pw_hash" json:"pw_hash"`
		NameFirst string `db:"name_first" json:"name_first" binding:"required"`
		NameLast  string `db:"name_last" json:"name_last" binding:"required"`
	}

	var json PersonEasy
	err = c.Bind(&json)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": "Unable to parse person object"})
		return
	}

	// disallow duplicate emails
	n, err := dbmap.SelectInt("select count(*) from person where email=$1 and person_id<>$2", json.Email, p.PersonID)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to check person database for duplicate email")})
		return
	}
	if n != 0 {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Email %s already exists.", json.Email)})
		return
	}

	var newHash string
	if json.PWHash != "" {
		newHash = HashPassword(json.PWHash)
	} else {
		newHash = p.PWHash
	}

	newPerson := Person{
		PersonID:  pid,
		Email:     json.Email,
		PWHash:    newHash,
		NameFirst: json.NameFirst,
		NameLast:  json.NameLast,
	}
	_, err = dbmap.Update(&newPerson)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{"error": fmt.Sprintf("Unable to update person %d", pid)})
		return
	}

	newPerson.PWHash = "<hidden>"
	result := gin.H{"person": newPerson}
	ServeResult(c, "person_update.js", result)
}

// HPersonCreate creates a person in the database
func HPersonCreate(c *gin.Context) {
	var person Person
	err := c.Bind(&person)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid person object"})
		return
	}

	// Encrypt the plaintext before saving
	person.PWHash = HashPassword(person.PWHash)

	err = person.Save()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	person.PWHash = "<hidden>"
	result := gin.H{"person": person}
	c.Writer.Header().Set("Location", fmt.Sprintf("/person/%d", person.PersonID))
	ServeResult(c, "person_create.js", result)
}

// HPersonDelete deletes a person from the database
func HPersonDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("person_id"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	p := &Person{PersonID: id}
	if err = p.Destroy(); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"result": "deleted"}
	ServeResult(c, "person_deleted.js", result)
}

// HPersonInfo return information for a person
func HPersonInfo(c *gin.Context) {
	pid := c.Param("person_id")

	pobj, err := dbmap.Get(Person{}, pid)
	if err != nil || pobj == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "not found"})
		return
	}
	p := pobj.(*Person)

	p.PWHash = "<hidden>"

	result := gin.H{"person": p}
	c.JSON(http.StatusOK, result)
}

// HPersonDecisions returns the decisions owned by that person
func HPersonDecisions(c *gin.Context) {
	id := c.Param("person_id")
	var decisions []Decision
	_, err := dbmap.Select(&decisions, "select * from decision where person_id=$1", id)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{"decisions": decisions}
	ServeResult(c, "person_decisions.js", result)
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
	_, _ = dbmap.Select(&decisions, "SELECT * from decision WHERE person_id=$1", p.PersonID)
	for _, d := range decisions {
		if err := d.Destroy(); err != nil {
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
			p.PersonID, err)
	}

	return nil
}
