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

type JSON map[string]interface{}
func (s *JSON) Scan(src interface{}) error {
	source, ok := src.([]byte)

	if !ok {
		return errors.New("Assert type: .([]byte)")
	}

	var i interface{}

	err := json.Unmarshal(source, &i)

	if err != nil {
		return err
	}

	*s, ok = i.(map[string]interface{})

	if !ok {
		return errors.New("Assert type: .(map[string]interface{})")
	}

	return nil
}


type LapsArr []Laps 
func (l *LapsArr) Scan(src interface{}) error {
	// The JSON method doesn't work for json arrays 
	source, ok := src.([]byte)

	if !ok {
		return errors.New("Assert type: .([]byte) failed")
	}

	var i []Laps

	err := json.Unmarshal(source, &i)

	if err != nil {
		return err
	}

	*l = i

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
	Summary JSON `db:"summary"`
	Laps LapsArr `db:"laps"`
	Stats JSON `db:"stats"`
}
func ServeLatestActivity(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get the latest activity

	a := new(ActivityResponse)

	query := fmt.Sprintf("SELECT * FROM activities ORDER BY id DESC LIMIT 1")

	err := db.QueryRow(query).Scan(&a.Name, &a.StartDate, &a.StartDateLocal,
							&a.Type, &a.ID, &a.ElapsedTime, &a.MovingTime,
								&a.Distance, &a.HasHeartRate, &a.Summary, 
								    &a.Laps, &a.Stats)
	
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(a)

	fmt.Println("[GET] Lastest activity - ", a.ID)
}


type Stream	 struct {
	Name string `db:"name"`
	Date string `db:"date"`
	DateLocal string `db:"date_local"`
	Type string `db:"type"`
	ID int64 `db:"id"`
	Time 				  JSON `db:"time"`
	Distance		      JSON `db:"distance"` 
	Latlng                JSON `db:"latlng"`
	Altitude              JSON `db:"altitude"`
	VelocitySmooth        JSON `db:"velocity_smooth"`
	Heartrate             JSON `db:"heartrate"`		
	Cadence               JSON `db:"cadence"`		
	Watts                 JSON `db:"watts"`		
	Temp                  JSON `db:"temperature"`	
	Moving                JSON `db:"moving"`	
	GradeSmooth           JSON `db:"grade_smooth"`
	GradeAdjustedDistance JSON `db:"grade_adjusted_distance"`
}
func ServeStream(w http.ResponseWriter, r *http.Request, db *sql.DB){
	// Serves up all stream data as a json object

	vars := mux.Vars(r)
	id := vars["id"]

	query := fmt.Sprintf("SELECT name, date, date_local, type, id, time, distance, latlng, altitude, velocity_smooth, heartrate, cadence, watts, temperature, moving, grade_smooth, grade_adjusted_distance FROM streams WHERE id='%s'", id)

	rows, err := db.Query(query)

	if err != nil{
		log.Fatal(err)
	}

	s := Stream{}

	for rows.Next() {
		err := rows.Scan(&s.Name, &s.Date, &s.DateLocal, &s.Type, &s.ID, &s.Time,
							&s.Distance, &s.Latlng, &s.Altitude, &s.VelocitySmooth, 
							&s.Heartrate, &s.Cadence, &s.Watts, &s.Temp, &s.Moving, 
							&s.GradeSmooth, &s.GradeAdjustedDistance)
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
		ServeStream(w, r, db)
	})

	router.HandleFunc("/activity/latest", func(w http.ResponseWriter, r *http.Request) {
		ServeLatestActivity(w, r, db)
	})

	http.ListenAndServe(":8080", router)
}
