package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Person_ID  int    `binding:"required"`
	Email      string `binding:"required"`
	PW_hash    string `binding:"required"`
	Name_First string `binding:"required"`
	Name_Last  string `binding:"required"`
}

func HPersonsList(c *gin.Context) {
	rows, err := database.DB.Query("SELECT * FROM person;")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	defer rows.Close()

	var persons []Person
	for rows.Next() {
		var p Person
		err = rows.Scan(&p.Person_ID, &p.Email, &p.PW_hash, &p.Name_First, &p.Name_Last)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		persons = append(persons, p)
	}

	if len(persons) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": persons})
}

func HPersonCreate(c *gin.Context) {
	var person Person
	err := c.Bind(&person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid person object"})
		return
	}

	err = person.CreatePerson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": person})
}

func HPersonDelete(c *gin.Context) {
	id := c.Param("person_id")
	_, err := database.DB.Exec("DELETE FROM person WHERE person_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

func HPersonInfo(c *gin.Context) {
	id := c.Param("person_id")
	rows, err := database.DB.Query("SELECT * FROM person where person_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	defer rows.Close()

	persons := []Person{}
	for rows.Next() {
		var p Person
		err = rows.Scan(&p.Person_ID, &p.Email, &p.PW_hash, &p.Name_First, &p.Name_Last)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		persons = append(persons, p)
	}

	if len(persons) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": persons})
}

func HPersonDecisions(c *gin.Context) {
	id := c.Param("person_id")
	rows, err := database.DB.Query("SELECT * FROM decision WHERE person_id=$1", id)
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

func (p *Person) CreatePerson() error {
	res, err := database.DB.Exec("INSERT INTO person VALUES($1,$2,$3,$4,$5)",
		p.Person_ID, p.Email, p.PW_hash, p.Name_First, p.Name_Last)
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}
