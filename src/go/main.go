package main

import (
)


func main() {

	client := InfluxConnection()

	db := PostgresConnection()

	//// API Rotuer
	HandleRequests(db, client)

	defer client.Close()
	defer db.Close()

}
