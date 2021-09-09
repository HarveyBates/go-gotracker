package main 

import (
	"fmt"
	"reflect"
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
)


type Activity struct {
	Name 			string 	`json:"name"`
	Distance 		float64 `json:"distance"`
	MovingTime 		int64 	`json:"moving_time"`
	ElapsedTime 	int64 	`json:"elapsed_time"`
	Type 			string 	`json:"type"`
	ID 				int64 	`json:"id"`
	ExternalID 		string 	`json:"external_id"`
	StartDate 		string 	`json:"start_date"`
	StartDateLocal 	string 	`json:"start_date_local"`
	Map struct {
		SummaryPolyine	string 	`json:"summary_polyline, omitempty"`
	} `json:"map, omitempty"`
	AvSpeed			float64 `json:"average_speed"`
	MaxSpeed		float64 `json:"max_speed"`
	AvCadence		float64 `json:"average_cadence"`
	AvWatts			float64 `json:"average_watts"`
	NormWatts		float64 `json:"weighted_average_watts"`
	MaxWatts		int64	`json:"max_watts"` 
	Kilojoules		float64 `json:"kilojoules"`
	HasHeartRate	bool	`json:"has_heartrate"`
	AvHeartRate		float64 `json:"average_heartrate, omitempty"`
	MaxHeartRate	float64	`json:"max_heartrate, omitempty"`
	ElevationGain 	float64 `json:"total_elevation_gain"`
	MaxElevation 	float64 `json:"elev_high"`
	MinElevation 	float64 `json:"elev_low"`
}
func GetActivity(accessToken string) []Activity {

	var activity []Activity

	var bearer = "Bearer " + accessToken
	url := "https://www.strava.com/api/v3/athlete/activities/?per_page=2" 
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

	err = json.Unmarshal(responseData, &activity)
	
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	for _, items := range activity {
		value := reflect.ValueOf(items)
		for i := 0; i < value.NumField(); i++ {
			fmt.Println(value.Field(i))
		}
	}

	return activity
}


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
	} `json:"distance"`
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


type HeartrateStream struct {
	Heartrate struct {
		Data []float64 `json:"data"`
		SeriesType string `json:"distance"`
		OriginalSize int `json:"original_size"`
		Resolution string `json:"resolution"`
	} `json:"heartrate"`
	Distance struct {
		Data []float64 `json:"data"`
		SeriesType string `json:"distance"`
		OriginalSize int `json:"original_size"`
		Resolution string `json:"high"`
	} `json:"distance"`
}
func GetHeartRate(activity string, accessToken string) ([]float64, []float64) {

	var hr HeartrateStream

	var bearer = "Bearer " + accessToken
	url := "https://www.strava.com/api/v3/activities/" + activity + "/streams?keys=heartrate&key_by_type=true"
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

	err = json.Unmarshal(responseData, &hr)
	
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return hr.Distance.Data, hr.Heartrate.Data
}


type CadenceStream struct {
	Cadence struct {
		Data []float64 `json:"data"`
		SeriesType string `json:"distance"`
		OriginalSize int `json:"original_size"`
		Resolution string `json:"resolution"`
	} `json:"cadence"`
	Distance struct {
		Data []float64 `json:"data"`
		SeriesType string `json:"distance"`
		OriginalSize int `json:"original_size"`
		Resolution string `json:"high"`
	} `json:"distance"`
}
func GetCadence(activity string, accessToken string) ([]float64, []float64) {

	var cadence CadenceStream 

	var bearer = "Bearer " + accessToken
	url := "https://www.strava.com/api/v3/activities/" + activity + "/streams?keys=cadence&key_by_type=true"
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

	err = json.Unmarshal(responseData, &cadence)
	
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return cadence.Distance.Data, cadence.Cadence.Data
}

