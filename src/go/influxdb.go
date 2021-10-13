package main

import (
	"fmt"
	"log"
	"context"
	"github.com/influxdata/influxdb-client-go/v2"
//	"time"
)

func InfluxConnection() {
	var token = GetEnvVariable("INFLUXDB_TOKEN")
	const bucket = "activities"
	const org = "user"

	client := influxdb2.NewClient("http://localhost:8086", token)

	q := fmt.Sprintf(`
        from(bucket: "activities") 
            |> range(start: time(v: "2021-09-24T22:45:25Z"), stop: time(v: "2021-09-25T02:17:58Z")) 
            |> filter(fn: (r) => r["_measurement"] == "bike-outdoors") 
	`)

	queryAPI := client.QueryAPI(org)

	result, err := queryAPI.Query(context.Background(), q)

	if err == nil {
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			fmt.Printf("ts: %v field: %v value: %v\n", result.Record().Time(), result.Record().Field(), result.Record().Value())
		}
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		log.Fatal(err)
	}

	defer client.Close()
}

// TODO Transform this into a json response
// We want a connection to the db and then multiple queries 
// Setup postgres for other queries
