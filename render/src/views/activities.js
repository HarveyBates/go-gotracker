import React from 'react';
import ReactECharts from 'echarts-for-react';
import { graphic } from 'echarts';
import {InfluxDB} from '@influxdata/influxdb-client'
import './activities.css';
import Map from './map';


export default class Activity extends React.Component {
	constructor() {
		super();
		this.state = {
			activity: [],
			records: [],
			laps: [],
			activity_name: "", 
			activity_id: "",
			total_distance: 0,
			avg_heart_rate: 0,
			max_heart_rate: 0,
			avg_running_cadence: 0,
			max_running_cadence: 0,
			avg_speed: 0,
			max_speed: 0,
			num_laps: 0,
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
			const recordsQuery = `from(bucket: "records") |> range(start: time(v: ${data.start_time}), stop: time(v: ${data.end_time})) |> filter(fn: (r) => r["_measurement"] == "${data.activity_name}") |> aggregateWindow(every: ${this.state.smoothing}, fn: mean)`

			const handleState = () => {
				this.setState({	activity: data, 
								records: records,
								laps: laps,
								activity_name: data.activity_name, 
								activity_id: data.activity_id,
								total_distance: data.total_distance,
								avg_heart_rate: data.avg_heart_rate,
								max_heart_rate: data.max_heart_rate,
								avg_running_cadence: data.avg_running_cadence,
								max_running_cadence: data.max_running_cadence,
								avg_speed: data.avg_speed,
								max_speed: data.max_speed,
								num_laps: data.num_laps,
								start_time: data.start_time,
								end_time: data.end_time,
								sport: data.sport });
			}

			var records = [];
			client.queryRows(recordsQuery, {
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

			// Add a bit of time to get the last data point
			// By default end time is excluded from the query
			var endDate = new Date(data.end_time);
			var endDate = (endDate.getTime() + 1000000) / 1000;
			// Get laps 
			const lapsQuery = `from(bucket: "laps") |> range(start: time(v: ${data.start_time}), stop: ${endDate}) |> filter(fn: (r) => r["_measurement"] == "${data.activity_name}")`

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
					handleState();
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
									smooth: true,
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
									smooth: true,
									symbol: "none",
									data: arr
								});
							}
							else if (currentField === "cadence"){
								series.push({
									name: currentField,
									color: 'rgba(44, 130, 201, 1)',
									type: "line",
									smooth: true,
									symbol: "none",
									data: arr
								});
							}
							else if (currentField === "speed"){
								series.push({
									name: currentField,
									color: 'rgba(1, 152, 117, 1)',
									type: "line",
									smooth: true,
									symbol: "none",
									yAxisIndex: 1,
									data: arr
								});
							} else {
								series.push({
									name: currentField,
									type: "line",
									smooth: true,
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
							var mins = pace.slice(0, 2);
							var secs = pace.slice(3);
							var datePace = new Date(0, 0, 0, 0, mins, secs, 0);
							arr.push([diff, datePace]);
						} 
						else if (currentField === "altitude" || currentField === "heart_rate"){
							arr.push([diff, Math.round(row._value)]);

						} else{
							arr.push([diff, row._value]);
						}
					}
					totalRows++;
				}


				// Handle laps
				var minPace = new Date(0, 0, 0, 0, 6, 0, 0); // Paces are handled as dates
				var excludeLapFields = ["sport", "sub_sport", "lap_trigger", "event_type", "event",
										"start_position_lat", "start_position_long", 
										"end_position_lat", "end_position_long", "message_index",
										"avg_fractional_cadence", "max_fractional_cadence",
										"enhanced_avg_speed", "enhanced_max_speed", "start_time"]
				var lapFields = [];
				totalRows = 0;
				currentField = this.state.laps[0]._field;
				var lapArr = [];
				var lapSeries = [];
				var lapIndex = 1;
				for (let i in this.state.laps){
					var row = this.state.laps[i];
					if (currentField !== row._field || totalRows === (this.state.laps.length - 1)) {
						if(excludeLapFields.indexOf(currentField) === -1 && 
							getAverage(lapArr) !== 0 &&
							lapFields.indexOf(row._field) === -1){
							lapFields.push(row._field);
							if (currentField === "avg_speed"){
								lapArr.push([lapIndex, minPace])
								lapSeries.push({
									name: currentField,
									type: "bar",
									barWidth: "30%",
									label: {
										show: true,
									},
									itemStyle: {
										color: new graphic.LinearGradient(0, 0, 0, 1, [
											{ offset: 0, color: '#0bab64' },
											{ offset: 1, color: '#3bb78f' }
										])
									},
									data: lapArr
								});
							} 
							else if (currentField === "avg_heart_rate") {
								lapArr.push([lapIndex, 0]);
								lapSeries.push({
									name: currentField,
									type: "bar",
									barWidth: "30%",
									yAxisIndex: 1,
									itemStyle: {
										color: new graphic.LinearGradient(0, 0, 0, 1, [
											{ offset: 0, color: '#dd1a1a' },
											{ offset: 1, color: '#d15959' }
										])
									},
									data: lapArr
								});
							}
							else {
								lapSeries.push({
									name: currentField,
									type: "bar",
									barWidth: "95%",
									label: {
										show: true,
									},
									data: lapArr
								});
							}
						}
						lapIndex = 1;
						currentField = row._field;
						lapArr = [];
					} 
					if (row._field === "avg_speed") {
						var pace = getPace(row._value)
						var mins = pace.slice(0, 2);
						var secs = pace.slice(3);
						var datePace = new Date(0, 0, 0, 0, mins, secs, 0);
						lapArr.push([lapIndex, datePace]);
					} 
					else if(row._field === "max_speed"){
						var pace = getPace(row._value);
						lapArr.push([lapIndex, pace]);
					} else{
						lapArr.push([lapIndex, row._value]);
					}
					if (excludeLapFields.indexOf(row._field) === -1) {
						lapIndex++;
					} else {
						lapIndex = 1;
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
						type: 'time',
						position: 'left',
						inverse: true,
						offset: 25,
						axisLabel: {
							color: 'rgba(1, 152, 117, 1)',
							formatter: '{mm}:{ss}'
						}
					},
					{
						type: 'value',
						position: 'right',
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

			const lapsOptions = {
				xAxis: {
					name: "Lap",
					nameLocation: 'center',
					nameGap: 25,
					min: 0,
					max: this.state.num_laps,
					type: 'category',
				},
				yAxis: [
				{
					type: 'time',
					inverse: true,
					axisLabel: {
						formatter: '{mm}:{ss}',
						color: 'rgba(1, 152, 117, 1)'
					},
					position: 'left',
				},
				{
					type: 'value',
					position: 'right',
					axisLabel: {
						color: 'rgba(240, 52, 52, 1)',
					},
				}
				],
				tooltip: {
					show: true,
					formatter: function (params) {
						var date = new Date(params[0].data[1]);
						var pace = date.getMinutes() + ":" + date.getSeconds();
						try {
							var hr = params[1].data[1];
						} catch {
							var hr = 0;
						}
						return (
							"Pace: <b>" + pace + "</b><br>Heart Rate: <b>" + hr + "</b>"
						);
					},
					trigger: 'axis',
				},
				toolbox: {
					show: true,
					right: 100,
					feature: {
						dataZoom: {},
					}
				},
				series: [lapSeries[2],
						lapSeries[0]
				]
			};

			console.log(lapSeries);


			function formatWord(str) {
				const titleCase = str
					.toLowerCase()
					.replaceAll("_", " ")
					.replaceAll("running", "")
					.replaceAll("total", "")
					.replaceAll("speed", "pace")
					.split(" ")
					.map(word => {
						return word.charAt(0).toUpperCase() + word.slice(1);
					})
					.join(" ");

				return titleCase;
			}

			const lapSummary = () => {
				var laps = []
				var addLaps = false;
				for (let i = 0; i < this.state.num_laps; i++){
					var lap = []
					lap.push({name: "Lap Num.", value: i + 1});
					for (let field in lapSeries){
						if (lapSeries[field].name === "avg_speed") {
							var name = formatWord(lapSeries[field].name);
							var pace = lapSeries[field].data[i][1];
							pace = pace.getMinutes() + ':' + pace.getSeconds()
							lap.push({name: name, value: pace});
						} else {
							var name = formatWord(lapSeries[field].name);
							lap.push({name: name, 
										value: lapSeries[field].data[i][1]});
						}
					}
					laps.push({number: i + 1, data: lap});
				}

				return (
					<table className="activity-table">
						<thead>
							<tr key="lap-table-head">
								{laps[0].data.map((info, _) => {
									return <th key={info.name}>{info.name}</th>
								})}
							</tr>
						</thead>
						<tbody>
							{laps.map((lap) => {
								return (<tr key={lap.name}>
									{lap.data.map((info, _) => {
										return <td key={lap.name + info.name}>{info.value}</td>
									})}
								</tr>)
							})}
						</tbody>
					</table>
				);
			};

			function formatTitle(str) {
				const titleCase = str
					.toLowerCase()
					.replaceAll("_", " ")
					.replace(/[0-9]/g, "")
					.split(" ")
					.map(word => {
						return word.charAt(0).toUpperCase() + word.slice(1);
					})
					.join(" ");

				return titleCase;
			}


			return (
				<div className="activity-page">
					<div className="main-chart-summary">
						<div className="section-head">
							<h3>Record</h3>
						</div>
						<div className="activity-summary">
							<div className="summary-box">
								<h4>{formatTitle(this.state.activity_name)}</h4>
								<h5>{new Date(this.state.start_time).toDateString()}</h5>
								<h5>{formatTitle(this.state.sport)} - {(this.state.total_distance / 1000).toFixed(2)} km</h5>
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
								<h4>Avg: {this.state.avg_running_cadence}</h4>
								<h5>Max: {this.state.max_running_cadence}</h5>
							</div>
						</div>
						<div className="chart">
							<ReactECharts option={recordsOptions} 
								theme={'macarons'} 
								style={{height: 500, width: '100%'}}/>
						</div>
					</div>
					<div className="lap-chart-summary">
						<div className="section-head">
							<h3>Laps</h3>
						</div>
						<div className="lap-chart-section">
							<div className="laps-chart">
								<ReactECharts option={lapsOptions} 
									theme={'macarons'} 
									style={{height: 400, width: '100%'}}/>
							</div>
							<div className="laps-summary">
								{lapSummary()}
							</div>
						</div>
					</div>
				</div>
			);
		}
	}
};

