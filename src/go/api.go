package main


import (
	"fmt"
	"log"
	"os"
	"time"
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


type ListActivities struct {
	ActivityName string `json:"activity_name"`
	ActivityId int64 `json:"activity_id"`
	Sport string `json:"sport"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
	TotalDistance float64 `json:"total_distance"`
}
func ServeActivities(w http.ResponseWriter, r *http.Request, db *sql.DB){

	actList := make([]ListActivities, 0)

	query := fmt.Sprintf("SELECT activity_name, activity_id, sport, start_time, end_time, total_distance  FROM cycling_session UNION SELECT activity_name, activity_id, sport, start_time, end_time, total_distance  FROM running_session UNION SELECT activity_name, activity_id, sport, start_time, end_time, total_distance  FROM swimming_session")

	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	var act ListActivities
	for rows.Next() {
		err := rows.Scan(&act.ActivityName, &act.ActivityId, &act.Sport,
			&act.StartTime, &act.EndTime, &act.TotalDistance)
		actList = append(actList, act)

		if err != nil {
			log.Fatal(err)
		}
	}

	json.NewEncoder(w).Encode(actList)

	fmt.Println("[GET] Activities summary")
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

	router.HandleFunc("/activities/", func(w http.ResponseWriter, r *http.Request) {
		ServeActivities(w, r, db)
	})

	http.ListenAndServe(":8080", router)
}
