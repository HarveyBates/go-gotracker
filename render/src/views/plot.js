import React from 'react';
import ReactECharts from 'echarts-for-react';
import './plot.css';

class Chart extends React.Component {
	constructor() {
		super();
		this.state = {
			data: [], 
			activity: "", 
			date: ""
		};
	}

	async componentDidMount() {
		// Initial state
		try {
			const response = await fetch('/activity/stream/5963598195', {headers:{
				"Accept": "application/json",
				"Content-Type": "application/json"}});
			const data = await response.json();
			this.setState({data: data, activity: data.Name, date: data.Date});
		} catch (error) {
			console.log(error);
		}
	}

	updateContent = async () => {
		try {
			const response = await fetch('/activity/stream/5999892579', {headers:{
				"Accept": "application/json",
				"Content-Type": "application/json"}});
			const data = await response.json();
			this.setState({data: data, activity: data.Name, date: data.Date});
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
					text: this.state.activity,
					subtext: this.state.date
				},
				xAxis: {
					name: "Distance (km)",
					type: 'value',
				},
				yAxis: {
					type: 'value',
				},
				series: [
					{
						name: "Heart Rate",
						type: 'line',
						symbol: 'none',
						data: hr
					},
					{
						name: "Cadence",
						type: 'line',
						symbol: 'none',
						data: rpm 
					},
					{
						name: "Watts",
						type: 'line',
						symbol: 'none',
						data: power 
					},
				],
				tooltip: {
					trigger: 'axis',
					axisPointer: {
						type: 'cross'
					}
				},
				dataZoom: [
					{
						show: true,
						type: "inside",
						filterMode: "none"
					}
				],
			};
			return (
				<div className="main-container">
					<button className="btn btn-secondary" onClick={this.updateContent}>
						Click Me
					</button>
					<div className="chart">
						<ReactECharts option={options} 
							theme={'macarons'} 
							style={{height: 'inherit', width: 'inherit'}}/>
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
								<td>{this.state.activity}</td>
								<td>{this.state.date}</td>
							</tr>
						</tbody>
					</table>
				</div>
			);
		}
	}
};

export default Chart;

