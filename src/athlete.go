package main 

import (
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
)

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
func GetStats(accessToken string) AthleteStats {

	var stats AthleteStats

	var bearer = "Bearer " + accessToken
	url := "https://www.strava.com/api/v3/athletes/" + GetEnvVariable("STRAVA_ATHLETE_ID") + "/stats"
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




