import React from 'react';
import ReactECharts from 'echarts-for-react';
import {InfluxDB, FluxTableMetaData} from '@influxdata/influxdb-client'
import './activities.css';
import moment from 'moment';
import Map from './map';

class Activity extends React.Component {
	constructor() {
		super();
		this.state = {
			activity: [],
			stream: "",
			activity_name: "", 
			activity_id: "",
			start_time: "",
			end_time: "",
			sport: ""
		};
	}

	async componentDidMount() {
		// Initial state
		try {
			// Get the most recent activity 
			var response = await fetch('/activity/latest', {headers:{
				"Accept": "application/json",
				"Content-Type": "application/json"}}, {mode: 'no-cors'});
			var data = await response.json();

			// Influx db query
			const inTok = process.env.REACT_APP_INFLUX_DB;
			const client = new InfluxDB({url: "http://localhost:8086", token: inTok}).getQueryApi("user");
			const query = `from(bucket: "records") |> range(start: time(v: ${data.start_time}), stop: time(v: ${data.end_time})) |> filter(fn: (r) => r["_measurement"] == "${data.activity_name}") |> aggregateWindow(every: 20s, fn: mean)`

			let influxResponse = [];
			client.queryRows(query, {
				next(row, tableMeta) {
					let record = tableMeta.toObject(row)
					influxResponse.push(record);
				},
				error(error) {
					console.error(error)
				},
				complete() {
					console.log('\nFinished SUCCESS')
				},
			})
			this.setState({	activity: data, 
							stream: influxResponse,
							activity_name: data.activity_name, 
							activity_id: data.activity_id,
							start_time: data.start_time,
							end_time: data.end_time,
							sport: data.sport });
		} catch (error) {
			console.log(error);
		}
	}


	render() {
		if (this.state.activity.length === 0) {
			console.log("Empty arr")
			return (
			<div>
				<p> Loading chart... </p>
			</div>);
		}
		else{
			console.log(this.state.activity)
			console.log(this.state.stream[0])
			
			const options = {
				title: {
					text: this.state.activity_name,
					subtext: moment(this.state.start_time).format("DD-MM-YYYY hh:mm A"),
					left: 100,
				},
				xAxis: {
					name: "Distance (km)",
					nameLocation: 'center',
					nameGap: -15,
					//max: Math.floor(distance.slice(-1)[0]) / 1000,
					type: 'value',
				},
				yAxis: [
					{
						type: 'value',
						position: 'left'
					},
					{
						type: 'value',
						position: 'right',
						name: "Altitude (m)"
					},
				],
				tooltip: {
					trigger: 'axis',
					axisPointer: {
						type: 'cross'
					}
				},
				toolbox: {
					show: true,
					right: 100,
					feature: {
						saveAsImage: {},
						dataZoom: {},
					}
				},
				dataZoom: [
					{
						show: true,
						realtime: true,
						start: 0,
						end: 100,
					}
				],
				series: [
					{
						name: "Heartrate (bpm)",
						color: 'rgba(255, 80, 80, 1)',
						type: 'line',
						symbol: 'none',
						data: []
					},
					{
						name: "Power Zones",
						type: "line",
						data: [],
						markArea: {
							label: {
								show: true,
								position: 'inside',
								color: "#000",
							},
							data: [
								[{
									name: "Zone 1",
									itemStyle: {
										color: 'rgba(85, 155, 255, 0.2)'
									},
										yAxis: 0
									},
									{
										yAxis: 145
									}
								],
								[{
									name: "Zone 2",
									itemStyle: {
										color: 'rgba(125, 220, 80, 0.2)'
									},
										yAxis: 145
									},
									{
										yAxis: 200
									}
								],
								[{
									name: "Zone 3",
									itemStyle: {
										color: 'rgba(200, 165, 80, 0.2)'
									},
										yAxis: 200
									},
									{
										yAxis: 261
									}
								],
								[{
									name: "Zone 4",
									itemStyle: {
										color: 'rgba(255, 55, 80, 0.2)'
									},
										yAxis: 261
									},
									{
										yAxis: 2000
									}
								]
							]
						}
					}
				]
			};
			return (
				<div className="main-container">
					<div className="chart-map">
						<div className="chart">
							<ReactECharts option={options} 
								theme={'macarons'} 
								style={{height: '100%', width: '100%'}}/>
						</div>
					</div>
					<table className="activity-table">
						<thead>
							<tr>
								<th>Key</th>
								<th>Value</th>
							</tr>
						</thead>
						<tbody>{Object.keys(this.state.activity).map((key, i) => (
								<tr key = {i}>
									<td>{key}</td>
									<td>{this.state.activity[key]}</td>
								</tr>
							)
							)}
						</tbody>
					</table>
				</div>
			);
		}
	}
};

export default Activity;

