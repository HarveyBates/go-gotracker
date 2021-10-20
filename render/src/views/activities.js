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
			this.setState({ activity: data,
							activity_name: data.activity_name, 
							activity_id: data.activity_id,
							start_time: data.start_time,
							end_time: data.end_time,
							sport: data.sport});

			// Influx db query
			const inTok = process.env.REACT_APP_INFLUX_DB;
			const client = new InfluxDB({url: "http://localhost:8086", token: inTok}).getQueryApi("user");
			const query = `from(bucket: "records") |> range(start: time(v: ${this.state.start_time}), stop: time(v: ${this.state.end_time})) |> filter(fn: (r) => r["_measurement"] == "${this.state.activity_name}") |> filter(fn: (r) => r["_field"] == "speed")`
			console.log(query);
		client.queryRows(query, {
		  next(row, tableMeta) {
			const o = tableMeta.toObject(row)
			  console.log(`${o._time} ${o._field}=${o._value}`);
		  },
		  error(error) {
			console.error(error)
		  },
		  complete() {
			console.log('\\nFinished SUCCESS')
		  },
		})

		} catch (error) {
			console.log(error);
		}
	}


	render() {
		if (this.state.activity.length === 0) {
			return (
			<div>
				<p> Loading chart... </p>
			</div>);
		}
		else{
			return (
				<div className="main-container">
					<div className="chart-map">
						<div className="chart">
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

