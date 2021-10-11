package main

import (
	"fmt"
	"context"
	"github.com/influxdata/influxdb-client-go/v2"
	"time"
)

func InfluxConnection() {
	const token = GetEnvVariable("INFLUXDB_TOKEN")
	const bucket = "activities"
	const org = "user"

	client := influxdb2.NewClient("http://localhost:8086", token)

	query: = fmt.Sprintf("a

	defer client.Close()
}
