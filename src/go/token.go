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

func GetEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}


type RefreshStravaAccess struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
func GetRefreshToken() string {

	var refresh RefreshStravaAccess 

	data := map[string]string{
		"client_id": GetEnvVariable("STRAVA_ID"),
		"client_secret": GetEnvVariable("STRAVA_SECRET"),
		"refresh_token": GetEnvVariable("STRAVA_REFRESH"),
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

	err = json.Unmarshal(responseData, &refresh)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return refresh.AccessToken
}
