package main 

import (
	"fmt"
	"log"
	"strings"
	"database/sql"
	"net/http"
	_ "github.com/lib/pq"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func ActivitiyChart(w http.ResponseWriter, _ *http.Request, db *sql.DB){

	var (
		watts string
		distance string
	)

	rows, err := db.Query("SELECT distance_stream, watts_stream FROM activities WHERE id='5984320971'")

	if err != nil{
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&distance, &watts)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	} 

	distance = strings.ReplaceAll(distance, "[", "")
	distance = strings.ReplaceAll(distance, "]", "")
	arrDistance := strings.Split(distance, ",")

	watts = strings.ReplaceAll(watts, "[", "")
	watts = strings.ReplaceAll(watts, "]", "")
	arrWatts := strings.Split(watts, ",")

	yValues := make([]opts.LineData, 0)
	for _, v := range arrWatts{
		yValues = append(yValues, opts.LineData{Value: v})
	}

	fmt.Println(yValues)

	xValues := make([]opts.LineData, 0)
	for _, v := range arrDistance{
		xValues = append(xValues, opts.LineData{Value: v})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{
			Title: "Watts",
			Subtitle: "Test plot",
	}))

	line.SetXAxis(xValues).
		AddSeries("Watts", yValues).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	line.Render(w)
}
