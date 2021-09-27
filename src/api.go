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

	if err != nil {
		return err
	}

	*s, ok = i.(map[string]interface{})

	if !ok {
		return errors.New("Assert type: .(map[string]interface{})")
	}

	return nil
}


type ActivityStream struct {
	Name string `db:"name"`
	Attributes Stream `db:"attributes"`	
	HeartRate Stream `db:"heartrate_stream"`
	Cadence Stream `db:"cadence_stream"`
	Watts Stream `db:"watts_stream"`
	Distance Stream `db:"distance_stream"`
}
func ServeActivity(w http.ResponseWriter, r *http.Request, db *sql.DB){

	vars := mux.Vars(r)
	id := vars["id"]

	query := fmt.Sprintf("SELECT name, attributes, heartrate_stream, cadence_stream, watts_stream, distance_stream FROM activities WHERE id='%s'", id)

	rows, err := db.Query(query)

	if err != nil{
		log.Fatal(err)
	}

	s := ActivityStream{}

	for rows.Next() {
		err := rows.Scan(&s.Name, &s.Attributes, &s.HeartRate, &s.Cadence, &s.Watts, &s.Distance)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	} 

	json.NewEncoder(w).Encode(s)

	fmt.Println("GET: Activity Stream")

}



func HandleRequests(db *sql.DB) {
	fmt.Println("Starting gotracker router")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/activity/stream/{id}", func(w http.ResponseWriter, r *http.Request) {
		ServeActivity(w, r, db)
	})

	http.ListenAndServe(":8080", router)
}
