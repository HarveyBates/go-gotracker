package main


import (
	"fmt"
	"log"
	"os"
	"time"
	"context"
	"errors"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/influxdata/influxdb-client-go/v2"
)


func GetEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}


func PostgresConnection() *sql.DB{
	const (
		hostname = "localhost"
		port = 27222
		username = "postgres"
		password = "admin"
		dbname= "gogotracker"
	)

	fmt.Print("Attempting to connect to postgres... ")

	conn := fmt.Sprintf("port=%d host=%s user=%s password=%s dbname=%s sslmode=disable",
		port, hostname, username, password, dbname)

	psqlDB, err := sql.Open("postgres", conn)	

	if err != nil {
		log.Fatal(err)
	}

	psqlDB.SetConnMaxLifetime(time.Minute * 3) // Timeout. Ensures conns close safely.
	psqlDB.SetMaxOpenConns(5)
	psqlDB.SetMaxIdleConns(5)

	fmt.Println("connected to postgres.")

	return psqlDB 
}


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
type JSON map[string]interface{}
func ServeLatestRunActivity(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var activity JSON

	query := fmt.Sprintf("SELECT row_to_json(activity) FROM (SELECT * FROM running_session ORDER BY activity_id DESC LIMIT 1) activity")

	err := db.QueryRow(query).Scan(&activity)

	if err != nil {
		log.Fatal(err)
	}


	json.NewEncoder(w).Encode(activity)

	fmt.Print("[GET] Lastest activity ")

	sport := activity["sport"].(string)
	id := activity["activity_id"].(float64)
	startTime := activity["start_time"].(string)
	endTime := activity["start_time"].(string)
	fmt.Println(sport, id, startTime, endTime)

}
	


func InfluxConnection() influxdb2.Client {
	var token = GetEnvVariable("INFLUXDB_TOKEN")

	fmt.Print("Attempting to connect to influxdb... ")

	client := influxdb2.NewClient("http://localhost:8086", token)

	fmt.Println("connected to influxdb.")

	return client
}


func GetRecord(client influxdb2.Client, bucket string, activityName int64, startTime string, endTime string) {
	const org = "user"

	q := fmt.Sprintf(`
        from(bucket: %s) 
            |> range(start: time(v: %s), stop: time(v: %s)) 
            |> filter(fn: (r) => r["_measurement"] == %s) 
	`, bucket, startTime, endTime, activityName)

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
}


func HandleRequests(db *sql.DB) {
	fmt.Println("Starting gotracker router")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/activity/latest/", func(w http.ResponseWriter, r *http.Request) {
		ServeLatestActivity(w, r, db)
	})

	http.ListenAndServe(":8080", router)
}
