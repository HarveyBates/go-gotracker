import React from 'react';
import ReactECharts from 'echarts-for-react';
import './activities.css';
import moment from 'moment';
import Map from './map';


class Activity extends React.Component {
	constructor() {
		super();
		this.state = {
			data: [], 
			name: "", 
			date: "",
			distance: 0,
			activityID: 0
		};
	}

	async componentDidMount() {
		// Initial state
		try {
			// Get the most recent activity 
			var response = await fetch('/activity/latest', {headers:{
				"Accept": "application/json",
				"Content-Type": "application/json"}});
			var data = await response.json();
			this.setState({name: data.Name, distance: data.Distance, date: data.StartDateLocal, activityID: data.ID});

			// Get the stream of this activity
			response = await fetch('/activity/' + this.state.activityID + "/stream", {headers:{
				"Accept": "application/json",
				"Content-Type": "application/json"}});
			data = await response.json();
			this.setState({data: data});
			console.log(data)
		} catch (error) {
			console.log(error);
		}
	}

	updateContent = async () => {
		// Updated state
		try {
			const response = await fetch('/activity/stream/5999892579', {headers:{
				"Accept": "application/json",
				"Content-Type": "application/json"}});
			const data = await response.json();
			this.setState({data: data, name: data.Name, date: data.Attributes.start_date});
		} catch (error) {
			console.log(error);
		}
	}

	render() {
		
		if (this.state.data.length === 0) {
			return (
			<div>
				<p> Loading chart... </p>
			</div>);
		}
		else{
			const distance = this.state.data.Distance.data;
			const heartrate = this.state.data.Heartrate.data;
			const cadence = this.state.data.Cadence.data;
			const watts = this.state.data.Watts.data;
			const altitude = this.state.data.Altitude.data;

			const hr = [];
			const rpm = [];
			const power = []
			const alt = []
			
			for (let i in distance) {
				var dMeters = (distance[i] / 1000).toFixed(2)
				hr.push([dMeters, heartrate[i]]);
				rpm.push([dMeters, cadence[i]]);
				power.push([dMeters, watts[i]]);
				alt.push([dMeters, altitude[i]]);
			}
			
			const options = {
				title: {
					text: this.state.name,
					subtext: moment(this.state.date).format("DD-MM-YYYY hh:mm A"),
					left: 100,
				},
				xAxis: {
					name: "Distance (km)",
					nameLocation: 'center',
					nameGap: -15,
					max: Math.floor(distance.slice(-1)[0]) / 1000,
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
						data: hr
					},
					{
						name: "Cadence (rpm)",
						color: 'rgba(0, 0, 0, 1)',
						type: 'line',
						symbol: 'none',
						data: rpm 
					},
					{
						name: "Power (W)",
						color: 'rgba(50, 150, 255, 1)',
						type: 'line',
						symbol: 'none',
						data: power
					},
					{
						name: "Altitude (m)",
						color: 'rgba(190, 190, 190, 0.5)',
						type: 'line',
						areaStyle: {},
						symbol: 'none',
						yAxisIndex: 1,
						z: 0,
						data: alt 
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
							<button className="btn btn-secondary" onClick={this.updateContent}>
								Previous workout
							</button>
							<ReactECharts option={options} 
								theme={'macarons'} 
								style={{height: '100%', width: '100%'}}/>
						</div>
					</div>
					<table className="activity-table">
						<thead>
							<tr>
								<th>Name</th>
								<th>Date</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td>{this.state.name}</td>
								<td>{this.state.date}</td>
							</tr>
						</tbody>
					</table>
					<Map />
				</div>
			);
		}
	}
};

export default Activity;

