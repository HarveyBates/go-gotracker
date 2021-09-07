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

func SecondDB(db *sql.DB){


	fmt.Println("Inserting values in table...")

	insert, err := db.Query("INSERT INTO test VALUES ( 2, 'TEST' )")

	if err != nil {
		log.Fatal(err)
	}

	defer insert.Close()
}

func PopulateStatsTable(db *sql.DB){
	fmt.Println("Creating \"recent_stats\" table in db...")

	create, err := db.Query("CREATE TABLE IF NOT EXISTS recent_stats (created_at datetime default CURRENT_TIMESTAMP, n_activities int, distance float, moving_time int, elapsed_time int, elevation_gain float, achievement_count int)")

	if err != nil {
		log.Fatal(err)
	}

	defer create.Close()

	statement := "INSERT INTO recent_stats (created_at, n_activities, distance, moving_time, elapsed_time, elevation_gain, achievement_count) VALUES (`$1`, `$2`, `$3`, `$4`, `$5`, `$6`, `$7`)"
	
	_, err = db.Exec(statement, 1, 12.0, 10, 10, 111.1, 1)

	if err != nil{
		log.Fatal(err)
	}
}



func main() {

	db := SqlConnect()

	//SecondDB(db)

	PopulateStatsTable(db)

	db.Close()


//	refreshToken := GetRefreshToken()

//	distanceArr, wattsArr := GetWatts("5901981172", refreshToken)

//	MakeChart(distanceArr, wattsArr, "Workout #1", "2 x 20 min @ 240 Watts")
}
