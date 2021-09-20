package main

import (
	"fmt"
	"log"
	"time"
	"database/sql"
	"net/http"
	_ "github.com/lib/pq"
)

type TestValues struct {
	Name string `json:"name"`
	DOB string `json:"date_of_birth"`
}


const (
	hostname = "127.0.0.1"
	port = 27222
	username = "postgres"
	password = "admin"
	dbname= "go-gotracker"
)

func SqlConnect() *sql.DB{
	fmt.Print("Attempting to connect to db... ")

	conn := fmt.Sprintf("port=%d host=%s user=%s password=%s dbname=%s sslmode=disable",
		port, hostname, username, password, dbname)

	db, err := sql.Open("postgres", conn)	

	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3) // Timeout. Ensures conns close safely.
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	fmt.Println("Connected")

	return db
}


func main() {

	db := SqlConnect()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ActivitiyChart(w, r, db)
	})
	http.ListenAndServe(":8081", nil)


//	token := GetRefreshToken()

//	activity := GetActivity(token, 10)

//	PopulateActivites(db, activity, token)
	//stats := GetStats(token)

	// Totals - Recent, YTD, All-time
	//PopulateRunStats(db, stats)
	//PopulateRideStats(db, stats)
	//PopulateSwimStats(db, stats)

//	PopulateStatsTable(db)

	db.Close()

//	distanceArr, wattsArr := GetWatts("5901981172", refreshToken)

//	MakeChart(distanceArr, wattsArr, "Workout #1", "2 x 20 min @ 240 Watts")
}
