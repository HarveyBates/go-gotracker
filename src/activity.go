package main 

import (
	"fmt"
	"strconv"
	"math"
	"strings"
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
	MinElevation 	float64 `json:"elev_low"`
	MaxElevation 	float64 `json:"elev_high"`
}
func GetActivities(accessToken string, nResults int) []Activity {
	/* Gets an array of activities from Strava.
	 *
	 * @param accessToken Token from Strava to access API.
	 * @param nResults Number of results to return (size of array)
	 *
	 * @return activity An array of activities.
	 */

	var activities []Activity

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

	err = json.Unmarshal(responseData, &activities)
	
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return activities
}

type Laps struct {
	ID            int64  `json:"id"`
	ResourceState int    `json:"resource_state"`
	Name          string `json:"name"`
	Activity      struct {
		ID            int64 `json:"id"`
		ResourceState int   `json:"resource_state"`
	} `json:"activity"`
	ElapsedTime        int       `json:"elapsed_time"`
	MovingTime         int       `json:"moving_time"`
	StartDate          string `json:"start_date"`
	StartDateLocal     string `json:"start_date_local"`
	Distance           float64   `json:"distance"`
	StartIndex         int       `json:"start_index"`
	EndIndex           int       `json:"end_index"`
	TotalElevationGain float64   `json:"total_elevation_gain"`
	AverageSpeed       float64   `json:"average_speed"`
	MaxSpeed           float64   `json:"max_speed"`
	AverageCadence     float64   `json:"average_cadence"`
	DeviceWatts        bool      `json:"device_watts"`
	AverageWatts       float64   `json:"average_watts"`
	AverageHeartrate   float64   `json:"average_heartrate"`
	MaxHeartrate       float64   `json:"max_heartrate"`
	LapIndex           int       `json:"lap_index"`
	Split              int       `json:"split"`
	PaceZone           int       `json:"pace_zone"`
}
func GetLaps(accessToken string, activityId int64) []Laps {

	var laps []Laps

	var bearer = "Bearer " + accessToken
	url := fmt.Sprintf("https://www.strava.com/api/v3/activities/%d/laps", activityId) 
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

	err = json.Unmarshal(responseData, &laps)
	
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return laps
}


func PopulateActivites(db *sql.DB, activities []Activity, accessToken string){
	/*
	 * Populate new table (activities) with indexing values (i.e. name, id etc.) and 
	 * JSONB responses from strava.
	 */

	createActivities, err := db.Query("CREATE TABLE IF NOT EXISTS activities (name text, date timestamp, date_local timestamp, type text, id bigint, elapsed_time bigint, moving_time bigint, distance real, has_heart_rate boolean, summary jsonb, laps jsonb, stats jsonb)")	

	if err != nil {
		log.Fatal(err)
	}

	defer createActivities.Close()

	for _, activity := range activities {

		// Check if activity already exists
		var exists bool
		query := fmt.Sprintf("SELECT EXISTS(SELECT id FROM activities WHERE id = %s)", strconv.FormatInt(activity.ID, 10))

		err := db.QueryRow(query).Scan(&exists)
		
		if err != nil{
			log.Fatal(err)
		}

		if(!exists) {
			
			fmt.Println("Adding activity:\t", activity.Name, "\t", activity.ID)

			// Convert to json objects
			activityJSON, err := json.Marshal(activity)

			if err != nil {
				log.Fatal(err)
			}

			laps := GetLaps(accessToken, activity.ID)
			lapsJSON, err := json.Marshal(laps)

			if err != nil {
				log.Fatal(err)
			}

			// Get Stream
			streams := GetStreams(db, activity, accessToken)

			// Esitmate watts if the activity is a run 
			var normWattsRun, avWattsRun int
			//var normWattsRunArr []int64
			if(strings.Contains(activity.Type, "Run")) {
				_, normWattsRun, avWattsRun = EstimateRunWatts(db, streams)
			}

			stats := CalcActivityStats(db, activity, normWattsRun, avWattsRun)
			statsJSON, err := json.Marshal(stats)

			if err != nil {
				log.Fatal(err)
			}

			statement, err := db.Prepare("INSERT INTO activities (name, date, date_local, type, id, elapsed_time, moving_time, distance, has_heart_rate, summary, laps, stats) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)")

			if err != nil {
				log.Fatal(err)
			}

			_, err = statement.Exec(
				activity.Name, 
				activity.StartDate, 
				activity.StartDateLocal, 
				activity.Type,
				activity.ID, 
				activity.ElapsedTime,
				activity.MovingTime,
				activity.Distance,
				activity.HasHeartRate,
				activityJSON,
				lapsJSON,
				statsJSON)

			if err != nil{
				log.Fatal(err)
			}


		} else {
			fmt.Println(activity.Name, "\talready exists... skipping")	
		}
	}
}


type ActivityStats struct {
	Exertion int
	RunExertion int
	SwimExertion int
	HRExertion int

	Variability float64
	RunVariability float64

	Intensity float64
	HRIntensity float64
	RunIntensity float64
	SwimIntensity float64

	Efficiency float64
	EstimatedEfficiency float64
}
func CalcActivityStats(db *sql.DB, activity Activity, normWattsRun int, avWattsRun int) ActivityStats {

	var stats ActivityStats

	// Get the most recent athlete stats and config
	var bFTP, rFTP, bThHr, rThHr, sThHr, restingHr, reserveHr, sThP int
	query := fmt.Sprintf("SELECT bike_ftp, run_ftp, bike_threshold_heartrate, run_threshold_heartrate, swim_threshold_heartrate, resting_heartrate, reserve_heartrate, swim_threshold_pace FROM athlete ORDER BY date DESC LIMIT 1")
	err := db.QueryRow(query).Scan(&bFTP, &rFTP, &bThHr, &rThHr, &sThHr, &restingHr, &reserveHr, &sThP)

	if err != nil {
		log.Fatal(err)
	}	

	if(strings.Contains(activity.Type, "Run")) {
		// Use real watts to calculate Exertion
		if(activity.NormWatts != 0) {
			stats.Intensity = float64(activity.NormWatts) / float64(rFTP)
			stats.Exertion = int(((float64(activity.MovingTime) * float64(activity.NormWatts) * 
						float64(stats.Intensity)) / (float64(rFTP) * float64(3600))) * float64(100)) 
			stats.Variability = float64(activity.NormWatts) / float64(activity.AvWatts)
		}
		// Use estimated watts to calculate RunExertion if real watts are not avalibale 
		if(normWattsRun != 0 && avWattsRun != 0) {
			stats.RunIntensity = float64(normWattsRun) / float64(rFTP)
			stats.RunExertion = int(((float64(activity.MovingTime) * float64(normWattsRun) * 
						float64(stats.RunIntensity)) / (float64(rFTP) * float64(3600))) * float64(100)) 
			stats.RunVariability = float64(normWattsRun) / float64(avWattsRun)
		}
		if(activity.HasHeartRate) {
			stats.HRIntensity = float64(activity.AvHeartRate) / float64(rThHr)
			hrr := (activity.AvHeartRate - float64(restingHr)) / float64(reserveHr)
			// Activity TRIMP
			trimp := ((float64(activity.MovingTime) / 60)) * hrr * 0.64 * math.Pow(math.E, (1.92 * hrr))
			// ltTrimp - TRIMP at latate threshold for 1 hour
			ltTrimp := 60 * ((float64(rThHr) - float64(restingHr)) / float64(reserveHr)) * 0.64 * math.Pow(math.E, (1.92 * ((float64(rThHr) - float64(restingHr)) / float64(reserveHr))))
			stats.HRExertion = int((float64(trimp) / float64(ltTrimp)) * 100)
			fmt.Println("Run HRExertion: ", stats.HRExertion)
		}
		if(activity.HasHeartRate && activity.NormWatts != 0) {
			stats.Efficiency = float64(activity.NormWatts) / float64(activity.AvHeartRate)
		} else if(activity.HasHeartRate && normWattsRun != 0) {
			stats.EstimatedEfficiency = float64(normWattsRun) / float64(activity.AvHeartRate)
		}
	}

	if(strings.Contains(activity.Type, "Ride")) {
		if(activity.NormWatts != 0) {
			stats.Intensity = float64(activity.NormWatts) / float64(bFTP)
			stats.Exertion = int(((float64(activity.MovingTime) * float64(activity.NormWatts) * 
						float64(stats.Intensity)) / (float64(bFTP) * float64(3600))) * float64(100)) 
			stats.Variability = float64(activity.NormWatts) / float64(activity.AvWatts)
		}
		if(activity.HasHeartRate) {
			stats.HRIntensity = float64(activity.AvHeartRate) / float64(bThHr)
			hrr := (activity.AvHeartRate - float64(restingHr)) / float64(reserveHr)
			trimp := (float64(activity.MovingTime) / 60) * hrr * 0.64 * (math.Pow(math.E, (1.92 * hrr)))
			ltTrimp := 60 * ((143 - float64(restingHr)) / float64(reserveHr)) * 0.64 * math.Pow(math.E, (1.92 * ((143 - float64(restingHr)) / float64(reserveHr))))
			stats.HRExertion = int((float64(trimp) / float64(ltTrimp)) * 100)
		}
		if(activity.HasHeartRate && activity.NormWatts != 0) {
			stats.Efficiency = float64(activity.NormWatts) / float64(activity.AvHeartRate)
		}
	}

	if(strings.Contains(activity.Type, "Swim")) {
		// Convert to meters per minute
		dt := activity.Distance / (float64(activity.MovingTime) / 60)
		// Express as a percentage of threshold pace to calculate intensity
		stats.SwimIntensity = float64(dt / ((100 / float64(sThP)) * 60))
		// SE = intensity^3 x movingTime (hours) x 100
		stats.SwimExertion = int(math.Pow(stats.SwimIntensity, 3) * ((float64(activity.MovingTime) / 60) / 60) * 100)
		// If using heartrate monitor 
		if(activity.HasHeartRate) {
			stats.HRIntensity = float64(activity.AvHeartRate) / float64(sThHr)
			hrr := (activity.AvHeartRate - float64(restingHr)) / float64(reserveHr)
			trimp := (float64(activity.MovingTime) / 60) * hrr * 0.64 * (math.Pow(math.E, (1.92 * hrr)))
			ltTrimp := 60 * ((143 - float64(restingHr)) / float64(reserveHr)) * 0.64 * math.Pow(math.E, (1.92 * ((143 - float64(restingHr)) / float64(reserveHr))))
			stats.HRExertion = int((float64(trimp) / float64(ltTrimp)) * 100)
		}
	}

	return stats
}


func EstimateRunWatts(db *sql.DB, streams Streams) ([]int64, int, int) {
	// Caluculted from https://github.com/SauceLLC/sauce4strava

	interval := streams.Time.Data
	GAD := streams.GradeAdjustedDistance.Data

	var weight float64
	query := fmt.Sprintf("SELECT weight FROM athlete ORDER BY date DESC LIMIT 1")
	err := db.QueryRow(query).Scan(&weight)

	if err != nil {
		log.Fatal(err)
	}

	var estimatedWatts []int64
	var rollingAvWatts []float64
	var sumWatts, rollingSum int64
	for i := 1; i < len(GAD); i++ {
		distance := GAD[i] - GAD[i - 1]
		step := interval[i] - interval[i - 1]
		j := 4.35 / ((1 / float64(weight) * (1 / float64(distance))))
		kj := j * 0.00024
		eWatts := int64(float64(kj) * 1000 / float64(step))
		sumWatts += eWatts
		estimatedWatts = append(estimatedWatts, eWatts)
		rollingSum += eWatts
		if i % 30 == 0 {
			av := math.Pow(float64((rollingSum) / 30), 4)
			rollingAvWatts = append(rollingAvWatts, float64(av))
			rollingSum = 0
		}
	}

	var avWatts int
	if(sumWatts != 0 && len(estimatedWatts) != 0) {
		avWatts = int(sumWatts / int64(len(estimatedWatts)))
	} else {
		fmt.Println("Zero Division Error")
	}

	var normSum float64
	for i := 0; i < len(rollingAvWatts); i++ {
		normSum += rollingAvWatts[i]
	}

	var normWatts int
	if(normSum != 0 && len(rollingAvWatts) != 0) {
		normWatts = int(math.Sqrt(math.Sqrt(float64(normSum / float64(len(rollingAvWatts))))))
	} else {
		fmt.Println("Zero Division Error")
	}

	return estimatedWatts, avWatts, normWatts
}


//func AddititionalParameters(db *sql.DB) {
//	
////	ftp := 261
////
////	varIndex := activity.NormWatts / activity.AvWatts
////	intensityFactor := activity.NormWatts / FTP
////	tss := ((activity.ElapsedTime * activity.NormWatts * intensityFactor) / (ftp * 3600)) * 100 
//
//	tss := 42
//
//	var date string;
//	var ftp, atlp, ctlp, tsbp int;
//	err := db.QueryRow("SELECT date, ftp, atl, ctl, tsb FROM training ORDER BY date DESC").Scan(&date, &ftp, &atlp, &ctlp, &tsbp)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	atl := atlp + ((tss - atlp) / 10)
//	ctl := ctlp + ((tss - ctlp) / 42)
//	tsb := ctl - atl
//
//	query := fmt.Sprintf("INSERT INTO training(date, ftp, atl, ctl, tsb) VALUES (NOW(), %d, %d, %d, %d)", ftp, atl, ctl, tsb)	
//
//	q, err := db.Query(query)	
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	defer q.Close()
//}

