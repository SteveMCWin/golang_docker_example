package main

import (
	"blog_crud/persons"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	persons_db := &persons.Db{}
	err := persons_db.InitDb()
	if err != nil {
		log.Fatal("Error opening db:", err)
	}

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("person", getPersons(persons_db))
		v1.GET("person/:id", getPersonById(persons_db))
		v1.GET("person/search", searchForPersons(persons_db))
		v1.POST("person", addPerson(persons_db))
		v1.PUT("person/:id", udpatePerson(persons_db))
		v1.DELETE("person/:id", deletePerson(persons_db))
		v1.OPTIONS("person", options)
	}

	router.Run(":8080")
}

func getPersons(persons_db *persons.Db) func(c *gin.Context){
	return func(c *gin.Context) {
		persons, err := persons_db.GetPersons()
		if err != nil {
			log.Fatal("Error getting persons:", err)
		}

		if persons == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No records found"})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": persons})
		}
	}
}

func getPersonById(persons_db *persons.Db) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		person, err := persons_db.GetPersonById(id)
		if err != nil {
			log.Fatal(err)
		}

		if person.FirstName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No records found"})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": person})
		}
	}
}

func searchForPersons(persons_db *persons.Db) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Query("name")
		persons, err := persons_db.FindPeopleByName(name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
		} else {
			c.JSON(http.StatusOK, gin.H{ "data": persons })
		}
	}
}

func addPerson(persons_db *persons.Db) func(c *gin.Context) {
	return func(c *gin.Context) {
		var json persons.Person

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		success, err := persons_db.AddPerson(json)
		if success {
			c.JSON(http.StatusOK, gin.H{"message": "Success"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}

func udpatePerson(persons_db *persons.Db) func(c *gin.Context) {
	return func(c *gin.Context) {
		var json persons.Person
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id_to_update, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		}

		success, err := persons_db.UpdatePerson(json, id_to_update)
		if success {
			c.JSON(http.StatusOK, gin.H{"message": "Success"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}

func deletePerson(persons_db *persons.Db) func(c *gin.Context) {
	return func(c *gin.Context) {

		id_to_delete, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		}

		success, err := persons_db.DeletePerson(id_to_delete)
		if success {
			c.JSON(http.StatusOK, gin.H{"message": "Success"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}

func options(c *gin.Context) {

	ourOptions := "HTTP/1.1 200 OK\n" +
		"Allow: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Origin: http://locahost:8080\n" +
		"Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Headers: Content-Type\n"

	c.String(200, ourOptions)
}

