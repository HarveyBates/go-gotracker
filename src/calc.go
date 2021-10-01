package main 

import (
	"fmt"
	//"strconv"
	//"reflect"
	//"encoding/json"
	"log"
	//"io/ioutil"
	"database/sql"
	_ "github.com/lib/pq"
)

func Training(db *sql.DB) {

	cTraining, err := db.Query("CREATE TABLE IF NOT EXISTS training (date timestamptz NOT NULL DEFAULT NOW(), ftp integer, atl integer, ctl integer, tsb integer)")

	if err != nil {
		log.Fatal(err)
	}

	var populated bool
	query := fmt.Sprintf("SELECT COUNT(1) WHERE EXISTS (SELECT from training)")

	err = db.QueryRow(query).Scan(&populated)

	if err != nil {
		log.Fatal(err)
	}

	if(!populated){
		q, err := db.Query("INSERT INTO training VALUES (NOW(), 261, 0, 0, 0)")	

		if err != nil {
			log.Fatal(err)
		}
		defer q.Close()
	}
	CalculateTrainingParameters(db)
	defer cTraining.Close()
}



func CalculateTrainingParameters(db *sql.DB) {
	
//	ftp := 261
//
//	varIndex := activity.NormWatts / activity.AvWatts
//	intensityFactor := activity.NormWatts / FTP
//	tss := ((activity.ElapsedTime * activity.NormWatts * intensityFactor) / (ftp * 3600)) * 100 

	tss := 42

	var date string;
	var ftp, atlp, ctlp, tsbp int;
	err := db.QueryRow("SELECT date, ftp, atl, ctl, tsb FROM training ORDER BY date DESC").Scan(&date, &ftp, &atlp, &ctlp, &tsbp)

	if err != nil {
		log.Fatal(err)
	}

	atl := atlp + ((tss - atlp) / 10)
	ctl := ctlp + ((tss - ctlp) / 42)
	tsb := ctl - atl

	query := fmt.Sprintf("INSERT INTO training(date, ftp, atl, ctl, tsb) VALUES (NOW(), %d, %d, %d, %d)", ftp, atl, ctl, tsb)	

	q, err := db.Query(query)	

	if err != nil {
		log.Fatal(err)
	}

	defer q.Close()
}

//Acute Training Load (Fatigue):
//- Previous ATL - ATLp
//- Time constant - TC = ~5-10 days
//ATL = ATLp + ((TSS - ATLp)/TC)
//
//Chronic Training Load (Fitness):
//- Previous CTL - CTLp
//- TC = ~42 days
//CTL = CTLp + ((TSS - CTLp)/TC)
//
//Training Stress Balance (Form):
//TSB = CTL - ATL
