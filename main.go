package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	ID        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
}

var db *sqlx.DB

func initDB() {
	dsn := "root:salam@tcp(localhost:3306)/emp"
	var err error
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connected to the database")
}

func main() {
	initDB()
	defer db.Close()

	router := gin.Default()

	router.GET("/employees", getEmployees)
	router.GET("/employees/:id", getEmployee)
	router.POST("/employees", createEmployee)
	router.PUT("/employees/:id", updateEmployee)
	router.DELETE("/employees/:id", deleteEmployee)

	router.Run(":8080")
}

func getEmployees(c *gin.Context) {
	var employees []Employee
	err := db.Select(&employees, "SELECT * FROM employees")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employees)
}

func getEmployee(c *gin.Context) {
	id := c.Param("id")

	var employee Employee
	err := db.Get(&employee, "SELECT * FROM employees WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func createEmployee(c *gin.Context) {
	var employee Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO employees (first_name, last_name, email) VALUES (?, ?, ?)",
		employee.FirstName, employee.LastName, employee.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	employee.ID = int(id)

	c.JSON(http.StatusOK, employee)
}

func updateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE employees SET first_name = ?, last_name = ?, email = ? WHERE id = ?",
		employee.FirstName, employee.LastName, employee.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

func deleteEmployee(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
