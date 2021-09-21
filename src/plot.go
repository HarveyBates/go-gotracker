package main 

import (
	"log"
	"strings"
	"strconv"
	"math"
	"database/sql"
	"net/http"
	_ "github.com/lib/pq"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func StringToLineData(values string) []opts.LineData {

	values = strings.ReplaceAll(values, "[", "")
	values = strings.ReplaceAll(values, "]", "")
	arrValues := strings.Split(values, ",")

	LDValues := make([]opts.LineData, 0)
	for _, v := range arrValues {
		LDValues = append(LDValues, opts.LineData{Value: v})
	}

	return LDValues
}


func DistanceStringToLineData(values string) []opts.LineData {

	values = strings.ReplaceAll(values, "[", "")
	values = strings.ReplaceAll(values, "]", "")
	arrValues := strings.Split(values, ",")

	LDValues := make([]opts.LineData, 0)
	for _, v := range arrValues {

		val, err := strconv.ParseFloat(strings.TrimSpace(v), 32)

		if err != nil{
			log.Fatal(err)
		}	

		val = math.Round((val / 10)) / 100
		LDValues = append(LDValues, opts.LineData{Value: val})
	}

	return LDValues
}



func ActivitiyChart(w http.ResponseWriter, _ *http.Request, db *sql.DB){

	var (
		name string
		start_date_local string
		heartrate string
		cadence string
		watts string
		distance string
	)

	rows, err := db.Query("SELECT name, start_date_local, heartrate_stream, cadence_stream, watts_stream, distance_stream FROM activities WHERE id='5984320971'")

	if err != nil{
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&name, &start_date_local, &heartrate, &cadence, &watts, &distance)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	} 

	heartrateSeries := StringToLineData(heartrate)
	cadenceSeries := StringToLineData(cadence)
	wattsSeries := StringToLineData(watts)
	distanceSeries := DistanceStringToLineData(distance)

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{
			Title: name,
			Subtitle: start_date_local,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
		charts.WithXAxisOpts(opts.XAxis{Name: "Distance (km)"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
	)

	line.SetXAxis(distanceSeries).
		AddSeries("Heart Rate", heartrateSeries).
		AddSeries("Watts", wattsSeries).
		AddSeries("Cadence", cadenceSeries).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
			charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{Name: "Avg", Type: "average"}),)

	line.Render(w)
}
