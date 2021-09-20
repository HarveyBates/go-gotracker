package main 

import (
	"fmt"
	"strconv"
	//"reflect"
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
	/* Gets an array of activities from Strava.
	 *
	 * @param accessToken Token from Strava to access API.
	 * @param nResults Number of results to return (size of array)
	 *
	 * @return activity An array of activities.
	 */

	var activity []Activity

	var bearer = "Bearer " + accessToken
	url := fmt.Sprintf("https://www.strava.com/api/v3/athlete/activities/?per_page=%d", nResults) 
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

//	for _, items := range activity {
//		value := reflect.ValueOf(items)
//		for i := 0; i < value.NumField(); i++ {
//			fmt.Println(value.Field(i))
//		}
//	}

	return activity
}


type DistanceStream struct {
	Distance struct {
		Data []float64 `json:"data"`
		SeriesType string `json:"distance"`
		OriginalSize int `json:"original_size"`
		Resolution string `json:"high"`
	} `json:"distance"`
}
func GetDistance(activity string, accessToken string) []float64 {

	var distance DistanceStream 

	var bearer = "Bearer " + accessToken
	url := "https://www.strava.com/api/v3/activities/" + activity + "/streams?keys=distance&key_by_type=true"
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

	err = json.Unmarshal(responseData, &distance)
	
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return distance.Distance.Data
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
func GetCadence(activity string, accessToken string) []float64 {

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

	return cadence.Cadence.Data
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
func GetHeartRate(activity string, accessToken string) []float64 {

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

	return hr.Heartrate.Data
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
func GetWatts(activity string, accessToken string) []float64 {

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

	return watts.Watts.Data
}


func PopulateActivites(db *sql.DB, activities []Activity, accessToken string){
	/*
	 * Populate new table (activities) with indexing values (i.e. name, id etc.) and 
	 * JSONB responses from strava.
	 */

	createRide, err := db.Query("CREATE TABLE IF NOT EXISTS activities (name text, type text, id bigint, start_date_local date, attributes jsonb, distance_stream jsonb, cadence_stream jsonb, heartrate_stream jsonb, watts_stream jsonb)")	

	if err != nil {
		log.Fatal(err)
	}

	defer createRide.Close()

	for _, activity := range activities {

		var exists bool
		query := fmt.Sprintf("SELECT EXISTS(SELECT id FROM activities WHERE id = %s)", strconv.FormatInt(activity.ID, 10))

		err := db.QueryRow(query).Scan(&exists)
		
		if err != nil{
			log.Fatal(err)
		}

		if(!exists) {
			
			fmt.Println("Adding activity:\t", activity.Name)

			// Convert to json objects
			activityStruct, err := json.Marshal(activity)

			// Get distance
			distanceStruct := GetDistance(strconv.FormatInt(activity.ID, 10), accessToken)
			distance, err := json.Marshal(distanceStruct)

			// Get Cadence
			cadenceStruct := GetCadence(strconv.FormatInt(activity.ID, 10), accessToken)
			cadence, err := json.Marshal(cadenceStruct)

			// Get heart rate
			hrStruct := GetHeartRate(strconv.FormatInt(activity.ID, 10), accessToken)
			heartRate, err := json.Marshal(hrStruct)

			// Get watts
			wattsStruct := GetWatts(strconv.FormatInt(activity.ID, 10), accessToken)
			watts, err := json.Marshal(wattsStruct)

			statementRide, err := db.Prepare("INSERT INTO activities(name, type, id, start_date_local, attributes, distance_stream, cadence_stream, heartrate_stream, watts_stream) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")

			if err != nil {
				log.Fatal(err)
			}

			_, err = statementRide.Exec(
				activity.Name, 
				activity.Type,
				activity.ID, 
				activity.StartDateLocal,
				string(activityStruct),
				string(distance), 
				string(cadence), 
				string(heartRate), 
				string(watts))

			if err != nil{
				log.Fatal(err)
			}
		} else {
			fmt.Println(activity.Name, "\talready exists... skipping")	
		}
	}
}
