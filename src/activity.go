package main 

import (
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
)


type WattsStream struct {
	Watts struct {
		Data []float64 `json:"data"`
		SeriesType string `json:"distance"`
		OriginalSize int `json:"original_size"`
		Resolution string `json:"resolution"`
	} `json:"watts"`
	Distance struct {
		Data []float64 `json:"data"`
		SeriesType string `json:"distance"`
		OriginalSize int `json:"original_size"`
		Resolution string `json:"high"`
	}
}
func GetWatts(activity string, accessToken string) ([]float64, []float64) {

	var watts WattsStream

	var bearer = "Bearer " + accessToken
	url := "https://www.strava.com/api/v3/activities/" + activity + "/streams?keys=watts&key_by_type=true"
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", bearer)

	client := &http.Client{}
	response, err := client.Do(request)
	
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(responseData, &watts)
	
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return watts.Distance.Data, watts.Watts.Data
}
