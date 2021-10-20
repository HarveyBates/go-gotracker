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


func InfluxConnection() influxdb2.Client {
	var token = GetEnvVariable("INFLUXDB_TOKEN")

	fmt.Print("Attempting to connect to influxdb... ")

	client := influxdb2.NewClient("http://localhost:8086", token)

	fmt.Println("connected to influxdb.")

	return client
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
func ServeLatestActivity(w http.ResponseWriter, r *http.Request, db *sql.DB) {

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


func GetActivity(activity_id string, db *sql.DB) JSON{

	var activity JSON

	query := fmt.Sprintf("SELECT row_to_json(activity) FROM (SELECT * FROM running_session WHERE activity_id = %s LIMIT 1) activity", activity_id)

	err := db.QueryRow(query).Scan(&activity)

	if err != nil {
		log.Fatal(err)
	}

	return activity
}


func ServeRecord(w http.ResponseWriter, r *http.Request, db *sql.DB, client influxdb2.Client) {

	vars := mux.Vars(r)

	activity_id := vars["id"]

	activity := GetActivity(activity_id, db)

	q := fmt.Sprintf(`
        from(bucket: "%s") 
            |> range(start: time(v: %s), stop: time(v: %s)) 
            |> filter(fn: (r) => r["_measurement"] == "%s") 
	`, "records", activity["start_time"].(string), activity["end_time"].(string), activity["activity_name"].(string))

	fmt.Println(q)

	const org = "user"
	queryAPI := client.QueryAPI(org)

	result, err := queryAPI.Query(context.Background(), q)

	// Get record field
	jsonResponse := make(map[string][]map[time.Time]interface{})
	if err == nil {
		var fields []string
		values := make(map[time.Time]interface{})
		for result.Next() {
			fieldChange := false
			currentField := result.Record().Field()
			if !Contains(fields, currentField) {
				fields = append(fields, currentField)
				fieldChange = true
			}
			currentValue := result.Record().Value()
			currentTs := result.Record().Time()
			if (!fieldChange) {
				// Add value with timestamp pair
				values[currentTs] = currentValue
			} else {
				jsonResponse[fields[len(fields)-1]] = append(jsonResponse[fields[len(fields)-1]], values)
				values[currentTs] = currentValue
				fmt.Println("Field changed", currentField)	
			}
		}	
		//fmt.Println(jsonResponse["Cadence"].(interface{}))
//	if err == nil {
//		for result.Next() {
//			if result.TableChanged() {
//				fmt.Printf("table: %s\n", result.TableMetadata().String())
//			}
//
//			// TODO need to have something like:
//
//				// "speed": [
//					//	"timestamp": "value",
//					//	"timestamp": "value"
//				//	]
//
//			fmt.Printf("ts: %v field: %v value: %v\n", result.Record().Time(), result.Record().Field(), result.Record().Value())
//		}
//		if result.Err() != nil {
//			fmt.Printf("query parsing error: %s\n", result.Err().Error())
//		}
	fmt.Println(fields)
	} else {
		log.Fatal(err)
	}
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}


func HandleRequests(db *sql.DB, client influxdb2.Client) {
	fmt.Println("Starting gotracker router")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/activity/latest/", func(w http.ResponseWriter, r *http.Request) {
		ServeLatestActivity(w, r, db)
	})

	router.HandleFunc("/activity/{type}/{id}/", func(w http.ResponseWriter, r *http.Request) {
		ServeRecord(w, r, db, client)
	})

	http.ListenAndServe(":8080", router)
}
