package main 

import (
	"time"
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"database/sql"
	_ "github.com/lib/pq"
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


func PopulateRunStats(db *sql.DB, stats AthleteStats){

	createRecent, err := db.Query("CREATE TABLE IF NOT EXISTS recent_run_stats (created_at datetime DEFAULT CURRENT_TIMESTAMP, n_activities int, distance float, moving_time int, elapsed_time int, elevation_gain float, achievement_count int)")

	if err != nil {
		log.Fatal(err)
	}

	defer createRecent.Close()

	statementRecent, err := db.Prepare("INSERT INTO recent_run_stats(created_at, n_activities, distance, moving_time, elapsed_time, elevation_gain, achievement_count) VALUES (?, ?, ?, ?, ?, ?, ?)")

	if err != nil{
		log.Fatal(err)
	}

	_, err = statementRecent.Exec(
		time.Now(), 
		stats.RecentRun.Count, 
		stats.RecentRun.Distance, 
		stats.RecentRun.MovingTime, 
		stats.RecentRun.ElapsedTime,
		stats.RecentRun.ElevationGain,
		stats.RecentRun.AcheivementCount)

	if err != nil{
		log.Fatal(err)
	}

	createYtd, err := db.Query("CREATE TABLE IF NOT EXISTS ytd_run_stats (created_at datetime DEFAULT CURRENT_TIMESTAMP, n_activities int, distance float, moving_time int, elapsed_time int, elevation_gain float)")

	if err != nil {
		log.Fatal(err)
	}

	defer createYtd.Close()

	statementYtd, err := db.Prepare("INSERT INTO ytd_run_stats(created_at, n_activities, distance, moving_time, elapsed_time, elevation_gain) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil{
		log.Fatal(err)
	}

	_, err = statementYtd.Exec(
		time.Now(), 
		stats.AllRun.Count, 
		stats.AllRun.Distance, 
		stats.AllRun.MovingTime, 
		stats.AllRun.ElapsedTime,
		stats.AllRun.ElevationGain)

	if err != nil{
		log.Fatal(err)
	}

	createAll, err := db.Query("CREATE TABLE IF NOT EXISTS all_run_stats (created_at datetime DEFAULT CURRENT_TIMESTAMP, n_activities int, distance float, moving_time int, elapsed_time int, elevation_gain float)")

	if err != nil {
		log.Fatal(err)
	}

	defer createAll.Close()

	statementAll, err := db.Prepare("INSERT INTO all_run_stats(created_at, n_activities, distance, moving_time, elapsed_time, elevation_gain) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil{
		log.Fatal(err)
	}

	_, err = statementAll.Exec(
		time.Now(), 
		stats.AllRun.Count, 
		stats.AllRun.Distance, 
		stats.AllRun.MovingTime, 
		stats.AllRun.ElapsedTime,
		stats.AllRun.ElevationGain)

	if err != nil{
		log.Fatal(err)
	}
}

func PopulateRideStats(db *sql.DB, stats AthleteStats){

	createRecent, err := db.Query("CREATE TABLE IF NOT EXISTS recent_ride_stats (created_at datetime DEFAULT CURRENT_TIMESTAMP, n_activities int, distance float, moving_time int, elapsed_time int, elevation_gain float, achievement_count int)")

	if err != nil {
		log.Fatal(err)
	}

	defer createRecent.Close()

	statementRecent, err := db.Prepare("INSERT INTO recent_ride_stats(created_at, n_activities, distance, moving_time, elapsed_time, elevation_gain, achievement_count) VALUES (?, ?, ?, ?, ?, ?, ?)")

	if err != nil{
		log.Fatal(err)
	}

	_, err = statementRecent.Exec(
		time.Now(), 
		stats.RecentRide.Count, 
		stats.RecentRide.Distance, 
		stats.RecentRide.MovingTime, 
		stats.RecentRide.ElapsedTime,
		stats.RecentRide.ElevationGain,
		stats.RecentRide.AcheivementCount)

	if err != nil{
		log.Fatal(err)
	}

	createYtd, err := db.Query("CREATE TABLE IF NOT EXISTS ytd_ride_stats (created_at datetime DEFAULT CURRENT_TIMESTAMP, n_activities int, distance float, moving_time int, elapsed_time int, elevation_gain float)")

	if err != nil {
		log.Fatal(err)
	}

	defer createYtd.Close()

	statementYtd, err := db.Prepare("INSERT INTO ytd_ride_stats(created_at, n_activities, distance, moving_time, elapsed_time, elevation_gain) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil{
		log.Fatal(err)
	}

	_, err = statementYtd.Exec(
		time.Now(), 
		stats.YTDRide.Count, 
		stats.YTDRide.Distance, 
		stats.YTDRide.MovingTime, 
		stats.YTDRide.ElapsedTime,
		stats.YTDRide.ElevationGain)

	if err != nil{
		log.Fatal(err)
	}

	createAll, err := db.Query("CREATE TABLE IF NOT EXISTS all_ride_stats (created_at datetime DEFAULT CURRENT_TIMESTAMP, n_activities int, distance float, moving_time int, elapsed_time int, elevation_gain float)")

	if err != nil {
		log.Fatal(err)
	}

	defer createAll.Close()

	statementAll, err := db.Prepare("INSERT INTO all_ride_stats(created_at, n_activities, distance, moving_time, elapsed_time, elevation_gain) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil{
		log.Fatal(err)
	}

	_, err = statementAll.Exec(
		time.Now(), 
		stats.AllRide.Count, 
		stats.AllRide.Distance, 
		stats.AllRide.MovingTime, 
		stats.AllRide.ElapsedTime,
		stats.AllRide.ElevationGain)

	if err != nil{
		log.Fatal(err)
	}
}


func PopulateSwimStats(db *sql.DB, stats AthleteStats){

	createRecent, err := db.Query("CREATE TABLE IF NOT EXISTS recent_swim_stats (created_at datetime DEFAULT CURRENT_TIMESTAMP, n_activities int, distance float, moving_time int, elapsed_time int, elevation_gain float, achievement_count int)")

	if err != nil {
		log.Fatal(err)
	}

	defer createRecent.Close()

	statementRecent, err := db.Prepare("INSERT INTO recent_swim_stats(created_at, n_activities, distance, moving_time, elapsed_time, elevation_gain, achievement_count) VALUES (?, ?, ?, ?, ?, ?, ?)")

	if err != nil{
		log.Fatal(err)
	}

	_, err = statementRecent.Exec(
		time.Now(), 
		stats.RecentSwim.Count, 
		stats.RecentSwim.Distance, 
		stats.RecentSwim.MovingTime, 
		stats.RecentSwim.ElapsedTime,
		stats.RecentSwim.ElevationGain,
		stats.RecentSwim.AcheivementCount)

	if err != nil{
		log.Fatal(err)
	}

	createYtd, err := db.Query("CREATE TABLE IF NOT EXISTS ytd_swim_stats (created_at datetime DEFAULT CURRENT_TIMESTAMP, n_activities int, distance float, moving_time int, elapsed_time int, elevation_gain float)")

	if err != nil {
		log.Fatal(err)
	}

	defer createYtd.Close()

	statementYtd, err := db.Prepare("INSERT INTO ytd_swim_stats(created_at, n_activities, distance, moving_time, elapsed_time, elevation_gain) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil{
		log.Fatal(err)
	}

	_, err = statementYtd.Exec(
		time.Now(), 
		stats.YTDSwim.Count, 
		stats.YTDSwim.Distance, 
		stats.YTDSwim.MovingTime, 
		stats.YTDSwim.ElapsedTime,
		stats.YTDSwim.ElevationGain)

	if err != nil{
		log.Fatal(err)
	}

	createAll, err := db.Query("CREATE TABLE IF NOT EXISTS all_swim_stats (created_at datetime DEFAULT CURRENT_TIMESTAMP, n_activities int, distance float, moving_time int, elapsed_time int, elevation_gain float)")

	if err != nil {
		log.Fatal(err)
	}

	defer createAll.Close()

	statementAll, err := db.Prepare("INSERT INTO all_swim_stats(created_at, n_activities, distance, moving_time, elapsed_time, elevation_gain) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil{
		log.Fatal(err)
	}

	_, err = statementAll.Exec(
		time.Now(), 
		stats.AllSwim.Count, 
		stats.AllSwim.Distance, 
		stats.AllSwim.MovingTime, 
		stats.AllSwim.ElapsedTime,
		stats.AllSwim.ElevationGain)

	if err != nil{
		log.Fatal(err)
	}
}
