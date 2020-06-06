package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/MichaelS11/go-dht"
)

type reading struct {
	TimeStamp   string
	Temperature float64
	Humidity    float64
}

const GPIO = "GPIO17"
const Endpoint = "http://golangworkshop.com:5000/reading"

func main() {

	// grab our timestamp
	timeStamp := time.Now()

	hosterr := dht.HostInit()
	if hosterr != nil {
		fmt.Println("HostInit error:", hosterr)
		return
	}

	dht, dhterr := dht.NewDHT(GPIO, dht.Fahrenheit, "")
	if dhterr != nil {
		fmt.Println("NewDHT error:", dhterr)
		return
	}

	humidity, temperature, readerr := dht.Read()

	if readerr != nil {
		fmt.Println("Reader error:", readerr)
		return
	}

	newReading := reading{TimeStamp: timeStamp.Format("2006-01-02T15:04:05-0700"), Temperature: temperature, Humidity: humidity}

	fmt.Printf("Our Reading was: \nTemperature: %v\nHumidity:%v\n", temperature, humidity)

	var requestBody, reqerr = json.Marshal(newReading)

	if reqerr != nil {
		fmt.Println("Request error:", readerr)
		return
	}

	resp, resperror := http.Post(Endpoint, "application/json", bytes.NewBuffer(requestBody))

	if resperror != nil {
		fmt.Println("Response error:", resperror)
		return
	}

	defer resp.Body.Close()
}
