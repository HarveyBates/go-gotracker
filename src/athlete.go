package main

import (
//	"fmt"
	"time"
	"encoding/json"
	"log"
//	"io/ioutil"
	"database/sql"
	_ "github.com/lib/pq"
)

func CreateAthlete(db *sql.DB) {

	createAthlete, err := db.Query("CREATE TABLE IF NOT EXISTS athlete (date timestamp, date_local timestamp, weight real, ramp_rate integer, resting_heartrate integer, max_heartrate integer, reserve_heartrate real, bike_FTP integer, bike_power_zones jsonb, bike_heartrate_zones jsonb, bike_threshold_heartrate integer, run_FTP integer, run_power_zones jsonb, run_heartrate_zones jsonb, run_threshold_heartrate integer, run_threshold_pace integer, swim_threshold_pace integer, swim_threshold_heartrate integer)")

	if err != nil {
		log.Fatal(err)
	}

	defer createAthlete.Close()


	insertDefaults, err := db.Prepare("INSERT INTO athlete (date, date_local, weight, ramp_rate, resting_heartrate, max_heartrate, reserve_heartrate, bike_FTP, bike_power_zones, bike_heartrate_zones, bike_threshold_heartrate, run_FTP, run_power_zones, run_heartrate_zones, run_threshold_heartrate, run_threshold_pace, swim_threshold_pace, swim_threshold_heartrate) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)")

	if err != nil {
		log.Fatal(err)
	}

	restingHr := 44
	maxHr := 161
	reserveHr := maxHr - restingHr

	weight := 65.5
	rampRate := 5

	bikeftp := 261
	bikePower := CalculatePowerZones(bikeftp)
	bikePowerJSON, err := json.Marshal(bikePower)

	if err != nil {
		log.Fatal(err)
	}

	bikeThresholdHr := 143
	bikeHr := CalculateHeartRateZones(bikeThresholdHr)
	bikeHrJSON, err := json.Marshal(bikeHr)

	if err != nil {
		log.Fatal(err)
	}

	runftp := 258
	runPower := CalculatePowerZones(runftp)
	runPowerJSON, err := json.Marshal(runPower)

	if err != nil {
		log.Fatal(err)
	}

	runThresholdHr := 143
	runHr := CalculateHeartRateZones(runThresholdHr)
	runHrJSON, err := json.Marshal(runHr)
	
	if err != nil {
		log.Fatal(err)
	}

	runThresholdPace := 233 // Seconds
	swimThresholdPace := 92 // Seconds
	swimThresholdHr := 143

	timeZone, err := time.LoadLocation("Australia/Sydney")

	if err != nil {
		log.Fatal(err)
	}

	localTime := time.Now().In(timeZone)

	_, err = insertDefaults.Exec(
		time.Now(),
		localTime,
		weight,
		rampRate,
		restingHr,
		maxHr,
		reserveHr,
		bikeftp,
		bikePowerJSON,
		bikeHrJSON,
		bikeThresholdHr,
		runftp,
		runPowerJSON,
		runHrJSON,
		runThresholdHr,
		runThresholdPace,
		swimThresholdPace,
		swimThresholdHr)

	if err != nil{
		log.Fatal(err)
	}

}


type PowerZones struct {
	OneLower int
	OneUpper  int 
	TwoLower int 
	TwoUpper int 
	ThreeLower  int 
	ThreeUpper  int 
	FourLower  int 
	FourUpper  int 
	FiveLower int 
	FiveUpper  int 
	SixLower int 
	SixUpper int 
}
func CalculatePowerZones(ftp int) PowerZones {

	var zones PowerZones

	zones.OneLower = 0
	zones.OneUpper = int(float64(0.55) * float64(ftp))
	zones.TwoLower = zones.OneUpper + 1
	zones.TwoUpper = int(float64(0.75) * float64(ftp))
	zones.ThreeLower = zones.TwoUpper + 1
	zones.ThreeUpper = int(float64(0.9) * float64(ftp))
	zones.FourLower = zones.ThreeUpper + 1
	zones.FourUpper = int(float64(1.05) * float64(ftp))
	zones.FiveLower = zones.FourUpper + 1
	zones.FiveUpper = int(float64(1.2) * float64(ftp))
	zones.SixLower = zones.FiveUpper + 1
	zones.SixUpper = int(float64(1.5) * float64(ftp))

	return zones
}


type HeartRateZones struct {
	OneLower int 
	OneUpper  int 
	TwoLower int 
	TwoUpper int 
	ThreeLower  int 
	ThreeUpper  int 
	FourLower  int 
	FourUpper  int 
	FiveLower int 
}
func CalculateHeartRateZones(thresholdHr int) HeartRateZones {

	var zones HeartRateZones 

	zones.OneLower = 0
	zones.OneUpper = int(float64(0.68) * float64(thresholdHr))
	zones.TwoLower = zones.OneUpper + 1
	zones.TwoUpper = int(float64(0.83) * float64(thresholdHr))
	zones.ThreeLower = zones.TwoUpper + 1
	zones.ThreeUpper = int(float64(0.94) * float64(thresholdHr))
	zones.FourLower = zones.ThreeUpper + 1
	zones.FourUpper = int(float64(1.05) * float64(thresholdHr))
	zones.FiveLower = zones.FourUpper + 1

	return zones
}
