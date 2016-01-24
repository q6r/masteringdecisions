package main

// TODO : Change checklogin
// TODO : Improve all returned things
// TODO : .....

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rageix/ginAuth"
)

// AuthRequest is the request for logging in
type AuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// HAuthLogin expects {"email":<str>, "password":<str>}
// if we successed an encrypted cookie is created
func HAuthLogin(c *gin.Context) {
	var ar AuthRequest
	err := c.Bind(&ar)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid authrequest object"})
		return
	}

	// Find the user in the database
	hashed := HashPassword(ar.Password)
	var p Person
	err = dbmap.SelectOne(&p, "select * from person where email=$1 and pw_hash=$2", ar.Email, hashed)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wrong email or password"})
		return
	}

	extra := map[string]string{"email": ar.Email,
		"person_id": strconv.Itoa(p.Person_ID)}
	err = ginAuth.Login(c, extra)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "logged in"})
}

// HAuthLogout does a logout and removes
// the encrypted cookie
func HAuthLogout(c *gin.Context) {
	ginAuth.Logout(c)
}

// HAuthCheck returns the status of the current user
func HAuthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "authenticated"})
}

func HAuthUnauthenticated(c *gin.Context) {
	c.Abort()
}

func HAuthAuthenticated(c *gin.Context) {

}

// AuthAsAll is a middleware to be used
// after ginAuth.Use to assert authenticated users
// accepts all users
func AuthAsAll(c *gin.Context) {
	_, exists := c.Get("cookieData")
	if exists == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		c.Abort()
	}
}

// AuthAsAdmin is a middleware to be used
// after ginAuth.Use to assert authenticated
// users is only admin
func AuthAsAdmin(c *gin.Context) {
	got, exists := c.Get("cookieData")
	if exists == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		c.Abort()
	}

	cookie := got.(map[string]string)
	if cookie["person_id"] != "0" {
		c.JSON(http.StatusUnauthorized,
			gin.H{"error": fmt.Sprintf("%v is not an admin", cookie["email"])})
		c.Abort()
	}

}
