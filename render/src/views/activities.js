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
			records: [],
			activity_name: "", 
			activity_id: "",
			total_distance: 0,
			avg_heart_rate: 0,
			max_heart_rate: 0,
			avg_running_cadence: 0,
			max_running_cadence: 0,
			avg_speed: 0,
			max_speed: 0,
			start_time: "",
			end_time: "",
			sport: "",
			smoothing: "5s"
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

			const influxToken = process.env.REACT_APP_INFLUX_DB;

			// Get records (time-series of an activity)
			const client = new InfluxDB({url: "http://localhost:8086", token: influxToken}).getQueryApi("user");
			const query = `from(bucket: "records") |> range(start: time(v: ${data.start_time}), stop: time(v: ${data.end_time})) |> filter(fn: (r) => r["_measurement"] == "${data.activity_name}") |> aggregateWindow(every: ${this.state.smoothing}, fn: mean)`

			const handleState = () => {
				this.setState({	activity: data, 
								records: records,
								activity_name: data.activity_name, 
								activity_id: data.activity_id,
								total_distance: data.total_distance,
								avg_heart_rate: data.avg_heart_rate,
								max_heart_rate: data.max_heart_rate,
								avg_running_cadence: data.avg_running_cadence,
								max_running_cadence: data.max_running_cadence,
								avg_speed: data.avg_speed,
								max_speed: data.max_speed,
								start_time: data.start_time,
								end_time: data.end_time,
								sport: data.sport });
			}

			var records = [];
			client.queryRows(query, {
				next(row, tableMeta) {
					var record = tableMeta.toObject(row)
					records.push(record);
				},
				error(error) {
					console.error(error)
				},
				complete() {
					handleState();
				},
			})

			// Get laps 
			const lapsQuery = `from(bucket: "laps") |> range(start: time(v: ${data.start_time}), stop: time(v: ${data.end_time})) |> filter(fn: (r) => r["_measurement"] == "${data.activity_name}")`

			var laps = [];
			client.queryRows(lapsQuery, {
				next(row, tableMeta) {
					var lap = tableMeta.toObject(row)
					laps.push(lap);
				},
				error(error) {
					console.error(error)
				},
				complete() {
					console.log(laps);
				},
			})


		} catch (error) {
			console.log(error);
		}
	}


	render() {
		if (this.state.records.length === 0) {
			return (
			<div>
				<p> Loading chart... </p>
			</div>);
		}
		else{

			function getPace(dec) {
				var minutes = 16.66666666667 / (dec);
				var sign = minutes < 0 ? "-" : "";
				var min = Math.floor(Math.abs(minutes))
				var sec = Math.floor((Math.abs(minutes) * 60) % 60);
				return sign + (min < 10 ? "0" : "") + min + ":" + (sec < 10 ? "0" : "") + sec;
			}

			if (this.state.sport === "running"){
				var mainFields = ["altitude", "cadence", "heart_rate", "speed", 
									"Cadence", "Power", "Ground Time", "Leg Spring Stiffness",
										"Vertical Oscillation", "Form Power", "Elevation"]
				var currentField = this.state.records[0]._field;
				var series = [];
				var arr = [];
				
				// Altitude limits
				var minAlt = 0;
				var maxAlt = 2000;
				var altitudeSet = false;

				// X-axis date fix
				var startDate = new Date(this.state.start_time);

				function msConversion(millis) {
					let sec = Math.floor(millis / 1000);
					let hrs = Math.floor(sec / 3600);
					sec -= hrs * 3600;
					let min = Math.floor(sec / 60);
					sec -= min * 60;

					sec = '' + sec;
					sec = ('00' + sec).substring(sec.length);

					if (hrs > 0) {
						min = '' + min;
						min = ('00' + min).substring(min.length);
						return hrs + ":" + min + ":" + sec;
					}
					else {
						return min + ":" + sec;
					}
				}

				function getAverage(arr) {
					var sum = 0;
					for (let i = 0; i < arr.length; i++){
						sum += arr[i][1];
					}
					return sum / arr.length;
				}


				var minPace = 0;
				var paceSet = false;
				var avgPace = getPace(this.state.avg_speed);
				var maxPace = getPace(this.state.max_speed);

				// Handle Records
				var totalRows = 0;
				for (let i in this.state.records) {
					var row = this.state.records[i];
					if (currentField !== row._field || totalRows === (this.state.records.length - 1)) {
						// Append series if average is not zero
						if (mainFields.indexOf(currentField) !== -1 && 
								getAverage(arr) !== 0) {
							if (currentField === "altitude"){
								series.push({
									name: currentField,
									color: 'rgba(190, 190, 190, 0.5)',
									areaStyle: {},
									z: 0,
									type: "line",
									symbol: "none",
									yAxisIndex: 2,
									data: arr
								});
							}
							else if (currentField === "heart_rate"){
								series.push({
									name: currentField,
									color: 'rgba(240, 52, 52, 1)',
									type: "line",
									symbol: "none",
									data: arr
								});
							}
							else if (currentField === "cadence"){
								series.push({
									name: currentField,
									color: 'rgba(44, 130, 201, 1)',
									type: "line",
									symbol: "none",
									data: arr
								});
							}
							else if (currentField === "speed"){
								series.push({
									name: currentField,
									color: 'rgba(1, 152, 117, 1)',
									type: "line",
									symbol: "none",
									yAxisIndex: 1,
									data: arr
								});
							} else {
								series.push({
									name: currentField,
									type: "line",
									symbol: "none",
									data: arr
								});
							}
						}
						currentField = row._field; // Assign new field to currentField
						arr = []; // Reset array
					} else {
						// Set altitude min max limits on chart
						if (currentField === "altitude") {
							if (!altitudeSet && row._value !== null){
								minAlt = row._value;
								maxAlt = row._value;
								altitudeSet = true;
							}
							if (row._value < minAlt && row._value !== null) {
								minAlt = row._value;
							}
							if (row._value > maxAlt) {
								maxAlt = row._value;
							}
						}
						else if (currentField === "speed") {
							if (!paceSet && row._value !== null && row._value !== 0){
								minPace = row._value;
								paceSet = true;
							}
							if (row._value < minPace && row._value !== null) {
								minPace = row._value;
							}
						}
						var date = new Date(row._time);
						var diff = date.getTime() - startDate.getTime();
						diff = msConversion(diff)
						if (currentField === "speed"){
							var pace = getPace(row._value)
							arr.push([diff, pace]);
						} 
						else if (currentField === "altitude" || currentField === "heart_rate"){
							arr.push([diff, Math.round(row._value)]);

						} else{
							arr.push([diff, row._value]);
						}
					}
					totalRows++;
				}
			}

			const recordsOptions = {
				xAxis: {
					name: "Time",
					nameLocation: 'center',
					nameGap: -15,
					type: 'category',
				},
				yAxis: [
					{
						type: 'value',
						position: 'left'
					},
					{
						type: 'category',
						position: 'left',
						offset: 50,
						name: "Pace (min/km)",
						min: getPace(minPace),
						max: maxPace
					},
					{
						type: 'value',
						position: 'right',
						name: "Altitude (m)",
						splitLine: {
							show: false,
						},
						min: minAlt,
						max: maxAlt
					},
				],
				tooltip: {
					show: true,
					trigger: 'axis',
					axisPointer: {
						type: 'cross',
						label: {
							precision: '0'
						}
					},
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
				series: series
			};

			return (
				<div className="activity-page">
					<div className="main-chart-summary">
						<div className="activity-summary">
							<div className="summary-box">
								<h4>{this.state.activity_name}</h4>
								<h5>{this.state.start_time}</h5>
								<h5>{this.state.sport} - {(this.state.total_distance / 1000).toFixed(2)} km</h5>
							</div>
							<div className="summary-box">
								<h4 style={{color: "rgba(1, 152, 117, 1)"}}>Pace</h4>
								<h4>Avg: {avgPace}</h4>
								<h5>Max: {maxPace}</h5>
							</div>
							<div className="summary-box">
								<h4 style={{color: "rgba(240, 52, 52, 1)"}}>Heart Rate</h4>
								<h4>Avg: {this.state.avg_heart_rate}</h4>
								<h5>Max: {this.state.max_heart_rate}</h5>
							</div>
							<div className="summary-box">
								<h4 style={{color: "rgba(44, 130, 201, 1)"}}>Cadence</h4>
								<h4>Avg: {this.state.avg_running_cadence * 2}</h4>
								<h5>Max: {this.state.max_running_cadence * 2}</h5>
							</div>
						</div>
						<div className="chart">
							<ReactECharts option={recordsOptions} 
								theme={'macarons'} 
								style={{height: 350, width: '100%'}}/>
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

