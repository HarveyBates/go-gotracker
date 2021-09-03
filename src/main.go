package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"github.com/joho/godotenv"
	"net/http"
	"io/ioutil"
)

// return the value of the key
func getEnvVal(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getStravaRefreshToken() string {
	data := map[string]string{
		"client_id": getEnvVal("STRAVA_ID"),
		"client_secret": getEnvVal("STRAVA_SECRET"),
		"refresh_token": getEnvVal("STRAVA_REFRESH"),
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	request := "https://www.strava.com/oauth/token?grant_type=refresh_token"

	response, err := http.Post(request, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}	

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	return string(responseData)
}


func main() {

	refreshToken := getStravaRefreshToken()

	fmt.Print(refreshToken)

	// godotenv package
	//dotenv := getEnvVal("STRAVA_ID")

	//fmt.Printf("godotenv : %s = %s \n", "STRONGEST_AVENGER", dotenv)
}
