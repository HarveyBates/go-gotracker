package main 

import (
	"os"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func MakeChart(xValues []float64, yValues []float64, title string, subtitle string){

	values := make([]opts.LineData, 0)
	for _, v := range yValues{
		values = append(values, opts.LineData{Value: v})
	}
	
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title: title,
			Subtitle: subtitle,
	}))

	line.SetXAxis(xValues).
		AddSeries("Watts", values).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	
	f, _ := os.Create("plot.html")
	line.Render(f)
}
