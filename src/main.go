package main

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type TestValues struct {
	Name string `json:"name"`
	DOB string `json:"date_of_birth"`
}

func SqlConnect(){
	fmt.Println("Attempting to connect to MySQl server...")

	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3001)/activities")

	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(1000)

	defer db.Close()

	fmt.Println("Creating table in db...")

	create, err := db.Query("CREATE TABLE [IF NOT EXISTS] test ( ID int, NAME varchar(255) )")

	if err != nil {
		log.Fatal(err)
	}

	defer create.Close()

	fmt.Println("Inserting values in table...")

	insert, err := db.Query("INSERT INTO test VALUES ( 2, 'TEST' )")

	if err != nil {
		log.Fatal(err)
	}

	defer insert.Close()
}


func main() {

	SqlConnect()


//	refreshToken := GetRefreshToken()

//	distanceArr, wattsArr := GetWatts("5901981172", refreshToken)

//	MakeChart(distanceArr, wattsArr, "Workout #1", "2 x 20 min @ 240 Watts")
}
