package main

import (
	"fmt"
	"log"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type TestValues struct {
	Name string `json:"name"`
	DOB string `json:"date_of_birth"`
}

func SqlConnect() *sql.DB{
	fmt.Println("Attempting to connect to MySQl server...")

	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:33060)/gotracker")

	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3) // Timeout. Ensures conns close safely.
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	return db
}


func main() {

	db := SqlConnect()

	token := GetRefreshToken()

	GetActivity(token)

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
