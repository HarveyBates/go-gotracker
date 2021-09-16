package main 

import (
	"fmt"
	"strconv"
	"reflect"
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"database/sql"
	_ "github.com/lib/pq"
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
func GetActivity(accessToken string, nResults int) []Activity {

	var activity []Activity

	var bearer = "Bearer " + accessToken
	url := "https://www.strava.com/api/v3/athlete/activities/?per_page=1"
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


func PopulateRide(db *sql.DB, activity []Activity, accessToken string){
	createRide, err := db.Query("CREATE TABLE IF NOT EXISTS ride_activities (name text, total_distance float, moving_time int, elapsed_time int, type text, activity_id int, external_id text, start_date date, start_date_local date, map_polyline text, av_speed float, max_speed float, av_cadence float, normalised_watts float, max_watts int, kilojoules float, has_heartrate bool, av_heartrate float, max_heartrate float, elevation_gain float, max_elevation float, min_elevation float, cadence text, watts text, distance text, heartrate text)")	

	if err != nil {
		log.Fatal(err)
	}

	defer createRide.Close()

	// Get Cadence and Distance
	distance, cadence := GetCadence(strconv.FormatInt(activity[0].ID, 10), accessToken)

	strDistance := "" 
	for index, value := range distance {
		if index != len(distance) - 1{
			strDistance += fmt.Sprint(value) + ", "
		} else {
			strDistance += fmt.Sprint(value)
		}
	}

	strCadence := "" 
	for index, value := range cadence {
		if index != len(cadence) - 1{
			strCadence += fmt.Sprint(value) + ", "
		} else {
			strCadence += fmt.Sprint(value)
		}
	}

	// Get heart rate
	_, heartRate = GetHeartRate(activity.ID, accessToken)

	strHeartrate := "" 
	for index, value := range heartRate {
		if index != len(heartRate) - 1{
			strHeartrate += fmt.Sprint(value) + ", "
		} else {
			strHeartrate += fmt.Sprint(value)
		}
	}


	// Get watts
	_, watts = GetWatts(activity.ID, accessToken)

	strWatts := "" 
	for index, value := range watts {
		if index != len(watts) - 1{
			strWatts += fmt.Sprint(value) + ", "
		} else {
			strWatts += fmt.Sprint(value)
		}
	}
	
	// TODO needs some error handling for no-existant values

	
//	statement, err := db.Prepare("INSERT INTO ride_activites (name, distance, moving_time, elapsed_time, type, activity_id, external_id, start_date, start_date_local, map_polyline, av_speed, max_speed, av_cadence, normalised_watts, max_watts, kilojoules, has_heartrate, av_heartrate, max_heartrate, elevation_gain, max_elevation, min_elevation, cadence, watts, distance, heartrate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	_, err = statement.Recent.Exec(
//		activity.Name, 
//		activity.Distance,
//		acitivty.MovingTime,
//		activity.ElapsedTime, 
//		activity.Type,
//		activity.ID, 
//		activity.ExternalID, 
//		activity.StartDate, 
//		activity.StartDateLocal, 
//		activity.Map.SummaryPolyline,
//		activity.AvSpeed, 
//		activity.MaxSpeed, 
//		activity.AvCadence, 
//		activity.AvWatts,
//		activity.NormWatts, 
//		activity.MaxWatts,
//		activity.Kilojoules,
//		activity.HasHeartRate,
//		activity.AvHeartRate,
//		activity.MaxHeartRate,
//		activity.ElevationGain,
//		activity.MaxElevation,
//		activity.MinElevation
//	)


}
