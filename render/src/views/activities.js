import React from 'react';
import ReactECharts from 'echarts-for-react';
import './activities.css';
import moment from 'moment';


class Activities extends React.Component {
	constructor() {
		super();
		this.state = {
			data: [], 
			name: "", 
			date: "",
			distance: 0
		};
	}

	async componentDidMount() {
		// Initial state
		try {
			const response = await fetch('/activity/stream/5963598195', {headers:{
				"Accept": "application/json",
				"Content-Type": "application/json"}});
			const data = await response.json();
			this.setState({data: data, name: data.Name, date: data.Attributes.start_date, 
				distance: data.Attributes.distance});
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
			return (<div>
				<p> Loading chart... </p>
			</div>);
		}
		else{
			const distance = this.state.data.Distance.distance.data;
			const heartrate = this.state.data.HeartRate.heartrate.data;
			const cadence = this.state.data.Cadence.cadence.data;
			const watts = this.state.data.Watts.watts.data;

			const hr = [];
			const rpm = [];
			const power = []
			for (let i in distance) {
				hr.push([distance[i] / 1000, heartrate[i]]);
				rpm.push([distance[i] / 1000, cadence[i]]);
				power.push([distance[i] / 1000, watts[i]]);
			}
			
			const options = {
				title: {
					text: this.state.name,
					subtext: moment(this.state.date).format("DD-MM-YYYY hh:mm A")
				},
				xAxis: {
					name: "Distance (km)",
					nameLocation: 'center',
					nameGap: -15,
					max: distance.slice(-1)[0] / 1000,
					type: 'value',
				},
				yAxis: {
					type: 'value',
				},
				tooltip: {
					trigger: 'axis',
					axisPointer: {
						type: 'cross'
					}
				},
				toolbox: {
					show: true,
					feature: {
						saveAsImage: {},
						dataZoom: {},
						restore: {}
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
					<div className="chart">
						<button className="btn btn-secondary" onClick={this.updateContent}>
							Click Me
						</button>
						<ReactECharts option={options} 
							theme={'macarons'} 
							style={{height: '100%', width: '100%'}}/>
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
				</div>
			);
		}
	}
};

export default Activities;

