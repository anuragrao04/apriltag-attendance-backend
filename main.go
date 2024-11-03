package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Student represents a student record
type Student struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	SRN      string `json:"srn"`
	PRN      string `json:"prn"`
	Detected bool   `json:"detected"`
}

// connectDB establishes a connection to the SQLite database
func connectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./pes-people-2024-11-01.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// getStudents is a handler function for fetching students from a given table
func getStudents(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Query("table")
		if tableName == "" {
			c.JSON(400, gin.H{"error": "Table name is required"})
			return
		}

		var students []Student
		rows, err := db.Query(fmt.Sprintf("SELECT ROWID, srn, prn, name FROM %s", tableName))
		if err != nil {
			c.JSON(500, gin.H{"error": "No Such Class Exists"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var s Student
			err = rows.Scan(&s.ID, &s.SRN, &s.PRN, &s.Name)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to scan row"})
				return
			}
			s.Detected = false // Always set detected to false
			students = append(students, s)
		}

		c.JSON(200, students)
	}
}

func getTag(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Query("table")
		prn := c.Query("prn")
		log.Println(tableName, prn)

		if tableName == "" || prn == "" {
			c.String(400, "Both table and prn parameters are required")
			return
		}

		var rowID int
		err := db.QueryRow(fmt.Sprintf("SELECT ROWID FROM %s WHERE prn = ?", tableName), prn).Scan(&rowID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.String(404, "Not Found")
			} else {
				log.Println(err)
				c.String(500, "Database error")
			}
			return
		}

		c.String(200, fmt.Sprint(rowID))
	}
}

func main() {
	db := connectDB()
	defer db.Close()

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:2233", "https://attendance.anuragrao.live"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	r.GET("/students", getStudents(db))
	r.GET("/get-tag", getTag(db))

	log.Println("Server is running on port 6969")
	r.Run(":6969")
}
