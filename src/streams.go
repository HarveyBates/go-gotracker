package main 

import (
	"encoding/json"
	"strconv"
	"log"
	"net/http"
	"io/ioutil"
//	"database/sql"
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
func GetStreams(activity int64, accessToken string) Streams {

	var stream Streams

	var bearer = "Bearer " + accessToken
	url := "https://www.strava.com/api/v3/activities/" + strconv.FormatInt(activity, 10) + "/streams?keys=time,distance,latlng,altitude,velocity_smooth,heartrate,cadence,watts,temp,moving,grade_smooth,grade_adjusted_distance&key_by_type=true"
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

	return stream
}
