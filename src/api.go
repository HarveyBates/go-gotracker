package main

import (
	"fmt"
	"log"
	"errors"
	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

type Stream map[string]interface{}
func (s *Stream) Scan(src interface{}) error {
	source, ok := src.([]byte)

	if !ok {
		return errors.New("Assert type: .([]byte)")
	}

	var i interface{}

	err := json.Unmarshal(source, &i)

	fmt.Println(i)

	if err != nil {
		return err
	}

	*s, ok = i.(map[string]interface{})

	if !ok {
		return errors.New("Assert type: .(map[string]interface{})")
	}


	return nil
}


type ActivityResponse struct {
	Name string`db:"name"`
	StartDate string `db:"date"`
	StartDateLocal string `db:"date_local"`
	Type string `db:"type"`
	ID int64 `db:"id"`
	ElapsedTime int64 `db:"elapsed_time"`
	MovingTime int64 `db:"moving_time"`
	Distance float64 `db:"distance"`
	HasHeartRate bool `db:"has_heart_rate"`
	Summary Stream `db:"summary"`
	Laps Stream `db:"laps"`
	Stats Stream `db:"stats"`
}
func ServeLatestActivity(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	a := new(ActivityResponse)

	query := fmt.Sprintf("SELECT name, date, date_local, type, id, elapsed_time, moving_time, distance, has_heart_rate, summary, laps, stats FROM activities ORDER BY id DESC LIMIT 1")

	err := db.QueryRow(query).Scan(&a.Name, &a.StartDate, &a.StartDateLocal,
							&a.Type, &a.ID, &a.ElapsedTime, &a.MovingTime,
								&a.Distance, &a.HasHeartRate, &a.Summary, 
								    &a.Laps, &a.Stats)
	
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(a)

	fmt.Println("[GET] Last activity - ", a.ID)
}


type ActivityStream struct {
	Name string `db:"name"`
	Attributes Stream `db:"attributes"`	
	HeartRate Stream `db:"heartrate_stream"`
	Cadence Stream `db:"cadence_stream"`
	Watts Stream `db:"watts_stream"`
	Distance Stream `db:"distance_stream"`
	Altitude Stream `db:"altitude_stream"`
	LatLng Stream `db:"latlng_stream"`
}
func ServeActivity(w http.ResponseWriter, r *http.Request, db *sql.DB){

	vars := mux.Vars(r)
	id := vars["id"]

	query := fmt.Sprintf("SELECT name, attributes, heartrate_stream, cadence_stream, watts_stream, distance_stream, altitude_stream, latlng_stream FROM streams WHERE id='%s'", id)

	rows, err := db.Query(query)

	if err != nil{
		log.Fatal(err)
	}

	s := ActivityStream{}

	for rows.Next() {
		err := rows.Scan(&s.Name, &s.Attributes, &s.HeartRate, &s.Cadence, 
			&s.Watts, &s.Distance, &s.Altitude, &s.LatLng)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	} 

	json.NewEncoder(w).Encode(s)

	fmt.Println("[GET] Activity Stream - ", id)

}



func HandleRequests(db *sql.DB) {
	fmt.Println("Starting gotracker router")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/activity/{id}/stream", func(w http.ResponseWriter, r *http.Request) {
		ServeActivity(w, r, db)
	})

	router.HandleFunc("/activity/latest", func(w http.ResponseWriter, r *http.Request) {
		ServeLatestActivity(w, r, db)
	})

	http.ListenAndServe(":8080", router)
}
