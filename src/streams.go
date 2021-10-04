package main 

import (
	"encoding/json"
	"strconv"
	"log"
	"net/http"
	"io/ioutil"
	"database/sql"
	_ "github.com/lib/pq"
)


type Streams struct {
	Temp struct {
		Data         []int  `json:"data"`
		SeriesType   string `json:"series_type"`
		OriginalSize int    `json:"original_size"`
		Resolution   string `json:"resolution"`
	} `json:"temp"`
	Moving struct {
		Data         []bool `json:"data"`
		SeriesType   string `json:"series_type"`
		OriginalSize int    `json:"original_size"`
		Resolution   string `json:"resolution"`
	} `json:"moving"`
	Latlng struct {
		Data         [][]float64 `json:"data"`
		SeriesType   string      `json:"series_type"`
		OriginalSize int         `json:"original_size"`
		Resolution   string      `json:"resolution"`
	} `json:"latlng"`
	VelocitySmooth struct {
		Data         []float64 `json:"data"`
		SeriesType   string    `json:"series_type"`
		OriginalSize int       `json:"original_size"`
		Resolution   string    `json:"resolution"`
	} `json:"velocity_smooth"`
	GradeSmooth struct {
		Data         []float64 `json:"data"`
		SeriesType   string    `json:"series_type"`
		OriginalSize int       `json:"original_size"`
		Resolution   string    `json:"resolution"`
	} `json:"grade_smooth"`
	Distance struct {
		Data         []float64 `json:"data"`
		SeriesType   string    `json:"series_type"`
		OriginalSize int       `json:"original_size"`
		Resolution   string    `json:"resolution"`
	} `json:"distance"`
	Cadence struct {
		Data         []float64 `json:"data"`
		SeriesType   string    `json:"series_type"`
		OriginalSize int       `json:"original_size"`
		Resolution   string    `json:"resolution"`
	} `json:"cadence"`
	Watts struct {
		Data         []float64 `json:"data"`
		SeriesType   string    `json:"series_type"`
		OriginalSize int       `json:"original_size"`
		Resolution   string    `json:"resolution"`
	} `json:"watts"`
	Altitude struct {
		Data         []float64 `json:"data"`
		SeriesType   string    `json:"series_type"`
		OriginalSize int       `json:"original_size"`
		Resolution   string    `json:"resolution"`
	} `json:"altitude"`
	Heartrate struct {
		Data         []int  `json:"data"`
		SeriesType   string `json:"series_type"`
		OriginalSize int    `json:"original_size"`
		Resolution   string `json:"resolution"`
	} `json:"heartrate"`
	Time struct {
		Data         []int  `json:"data"`
		SeriesType   string `json:"series_type"`
		OriginalSize int    `json:"original_size"`
		Resolution   string `json:"resolution"`
	} `json:"time"`
	GradeAdjustedDistance struct {
		Data         []float64 `json:"data"`
		SeriesType   string    `json:"series_type"`
		OriginalSize int       `json:"original_size"`
		Resolution   string    `json:"resolution"`
	} `json:"grade_adjusted_distance"`
}
func GetStreams(db *sql.DB, activity Activity, accessToken string) Streams {

	var stream Streams

	var bearer = "Bearer " + accessToken
	url := "https://www.strava.com/api/v3/activities/" + strconv.FormatInt(activity.ID, 10) + "/streams?keys=time,distance,latlng,altitude,velocity_smooth,heartrate,cadence,watts,temp,moving,grade_smooth,grade_adjusted_distance&key_by_type=true"
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

	err = json.Unmarshal(responseData, &stream)
	
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	createStreams, err := db.Query("CREATE TABLE IF NOT EXISTS streams (name text, date timestamp, date_local timestamp, type text, id bigint, time jsonb, distance jsonb, latlng jsonb, altitude jsonb, velocity_smooth jsonb, heartrate jsonb, cadence jsonb, watts jsonb, temperature jsonb, moving jsonb, grade_smooth jsonb, grade_adjusted_distance jsonb)")	

	if err != nil {
		log.Fatal(err)
	}

	defer createStreams.Close()

	statement, err := db.Prepare("INSERT INTO streams (name, date, date_local, type, id, time, distance, latlng, altitude, velocity_smooth, heartrate, cadence, watts, temperature, moving, grade_smooth, grade_adjusted_distance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)")

	if err != nil {
		log.Fatal(err)
	}

	timeJSON, err := json.Marshal(stream.Time)

	if err != nil{
		log.Fatal(err)
	}

	distanceJSON, err := json.Marshal(stream.Distance)

	if err != nil{
		log.Fatal(err)
	}

	latlngJSON, err := json.Marshal(stream.Latlng)

	if err != nil{
		log.Fatal(err)
	}

	altitudeJSON, err := json.Marshal(stream.Altitude)

	if err != nil{
		log.Fatal(err)
	}

	velocityJSON, err := json.Marshal(stream.VelocitySmooth)

	if err != nil{
		log.Fatal(err)
	}

	heartrateJSON, err := json.Marshal(stream.Heartrate)

	if err != nil{
		log.Fatal(err)
	}

	cadenceJSON, err := json.Marshal(stream.Cadence)

	if err != nil{
		log.Fatal(err)
	}

	wattsJSON, err := json.Marshal(stream.Watts)

	if err != nil{
		log.Fatal(err)
	}

	tempJSON, err := json.Marshal(stream.Temp)

	if err != nil{
		log.Fatal(err)
	}

	movingJSON, err := json.Marshal(stream.Moving)

	if err != nil{
		log.Fatal(err)
	}

	gradeJSON, err := json.Marshal(stream.GradeSmooth)

	if err != nil{
		log.Fatal(err)
	}

	gadJSON, err := json.Marshal(stream.GradeAdjustedDistance)

	if err != nil{
		log.Fatal(err)
	}

	_, err = statement.Exec(
		activity.Name, 
		activity.StartDate, 
		activity.StartDateLocal, 
		activity.Type,
		activity.ID, 
		timeJSON, 
		distanceJSON, 
		latlngJSON, 
		altitudeJSON, 
		velocityJSON,
		heartrateJSON,
		cadenceJSON,
		wattsJSON,
		tempJSON,
		movingJSON,
		gradeJSON,
		gadJSON)

	if err != nil{
		log.Fatal(err)
	}
	return stream
}
