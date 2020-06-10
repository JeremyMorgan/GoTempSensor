package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Reading struct {
	TimeStamp   string
	Temperature float64
	Humidity    float64
}

var db *sql.DB

func init() {
	// set up the database
	fmt.Println("Connected!! ")
	db, dbError = sql.Open("sqlite3", "./readings.db")
	Check(dbError)
	statement, prepError := db.Prepare("CREATE TABLE IF NOT EXISTS reading (TimeStamp TEXT, Temperature NUMERIC, Humidity NUMERIC)")
	Check(prepError)
	statement.Exec()

}

func main() {

	r := gin.Default()

	// this is our GET that returns 10 temperature readings.

	r.GET("/reading", func(c *gin.Context) {
		lastTen := getLastTen()
		// stuff into a JSON object and return it
		c.JSON(200, gin.H{"message": lastTen})
	})

	r.POST("/reading", tempData)
	r.Run(":5000")
}

func tempData(c *gin.Context) {
	// pull from original post and put into our struct

	if c.Request.Method == "POST" {

		fmt.Println("We got here")

		var r Reading
		c.BindJSON(&r)

		// save to database here
		saveToDatabase(r.TimeStamp, r.Temperature, r.Humidity)

		c.JSON(http.StatusOK, gin.H{
			"status":  "Posted!",
			"Message": "This worked!",
		})
	}
}

func saveToDatabase(TimeStamp string, Temperature float64, Humidity float64) {

	statement, err := db.Prepare("INSERT INTO reading (TimeStamp, Temperature, Humidity) VALUES (?,?,?)")
	Check(err)

	_, err = statement.Exec(TimeStamp, Temperature, Humidity)
	Check(err)

}

func getLastTen() []Reading {

	// query the database for readings
	rows, _ := db.Query("SELECT TimeStamp, Temperature, Humidity from reading LIMIT 20")

	// create some temp variables
	var TimeStamp string
	var Temperature float64
	var Humidity float64

	// make a slice
	lastTen := make([]Reading, 10)

	// insert data into slice
	for rows.Next() {
		rows.Scan(&TimeStamp, &Temperature, &Humidity)
		lastTen = append(lastTen, Reading{TimeStamp: TimeStamp, Temperature: Temperature, Humidity: Humidity})
	}
	// return it
	return lastTen
}

func Check(e error) {

	if e != nil {
		panic(e)
	}
}
