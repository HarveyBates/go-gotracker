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
	"github.com/wcharczuk/go-chart"
)


// return the value of the key
func get_env_var(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}


type RefreshStravaAccess struct {
	TOKEN_TYPE string
	ACCESS_TOKEN string
	EXPIRES_AT int
	EXPIRES_IN int
	REFRESH_TOKEN string
}
func get_strava_refresh_token() RefreshStravaAccess {

	var refresh RefreshStravaAccess 

	data := map[string]string{
		"client_id": get_env_var("STRAVA_ID"),
		"client_secret": get_env_var("STRAVA_SECRET"),
		"refresh_token": get_env_var("STRAVA_REFRESH"),
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

	return refresh
}


type RecentTotals struct {
	Count int `json:"count"`
	Distance float64 `json:"distance"`
	MovingTime int `json:"moving_time"`
	ElapsedTime int `json:"elapsed_time"`
	ElevationGain float64 `json:"elevation_gain"`
	AcheivementCount int `json:"achievement_count"`
}

type AllTotals struct {
	Count int `json:"count"`
	Distance float64 `json:"distance"`
	MovingTime int `json:"moving_time"`
	ElapsedTime int `json:"elapsed_time"`
	ElevationGain float64 `json:"elevation_gain"`
}

type YTDTotals struct {
	Count int `json:"count"`
	Distance float64 `json:"distance"`
	MovingTime int `json:"moving_time"`
	ElapsedTime int `json:"elapsed_time"`
	ElevationGain float64 `json:"elevation_gain"`
}

type AthleteStats struct {
	BiggestRideDistance float64 `json:"biggest_ride_distance"`
	BiggestClimbElevationGain float64 `json:"biggest_climb_elevation_gain"`
	RecentRide RecentTotals `json:"recent_ride_totals"`
	AllRide AllTotals `json:"all_ride_totals"`
	YTDRide YTDTotals `json:"ytd_ride_totals"`
	RecentRun RecentTotals `json:"recent_run_totals"`
	YTDRun YTDTotals `json:"ytd_run_totals"`
	AllRun AllTotals `json:"all_run_totals"`
	RecentSwim RecentTotals `json:"recent_swim_totals"`
	YTDSwim YTDTotals `json:"ytd_swim_totals"`
	AllSwim AllTotals `json:"all_swim_totals"`
}
func get_athlete_stats(accessToken string) AthleteStats {

	var stats AthleteStats

	var bearer = "Bearer " + accessToken
	url := "https://www.strava.com/api/v3/athletes/" + get_env_var("STRAVA_ATHLETE_ID") + "/stats"
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

	err = json.Unmarshal(responseData, &stats)
	
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return stats
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
	}
}
func get_activity_watts(activity string, accessToken string) ([]float64, []float64) {

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



func main() {

	refreshToken := get_strava_refresh_token()
	//fmt.Println(refreshToken.REFRESH_TOKEN)
	
	//stats := get_athlete_stats(refreshToken.ACCESS_TOKEN)
//	fmt.Println(stats.AllRide.Count)

	GetBikeZones()

	distanceArr, wattsArr := get_activity_watts("5901981172", refreshToken.ACCESS_TOKEN)

	mainSeries := chart.ContinuousSeries{
		Name: "Workout Power",
		Style: chart.Style{
			Show: true,
			StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
			FillColor: chart.GetDefaultColor(0).WithAlpha(64),
		},
		XValues: distanceArr, 
		YValues: wattsArr, 
	}

	maSeries := &chart.SMASeries {
		InnerSeries: mainSeries,
	}

	graph := chart.Chart{
		Series: []chart.Series{
			mainSeries, 
			maSeries,
		},
	}

	pngGraph, err := os.Create("plot.png")
	if err != nil {
		log.Fatal(err)
	}

	err = graph.Render(chart.PNG, pngGraph)

	if err != nil {
		log.Fatal(err)
	}

	err = pngGraph.Close()

	if err != nil {
		log.Fatal(err)
	}

}
