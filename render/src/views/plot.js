import React from 'react';
import ReactECharts from 'echarts-for-react';
import './plot.css';

class Page extends React.Component {
	constructor() {
		super();
		this.state = { data: []};
	}

	async componentDidMount() {
		try {
		const response = await fetch('/activity/stream/5963598195', {headers:{
			"Accept": "application/json",
			"Content-Type": "application/json"}});
		const data = await response.json();
		this.setState({data: data});
		} catch (error) {
			console.log(error);
		}
	}

	render() {
		if (this.state.data.length === 0) {
			console.log("Null data");
			return (<div>
				<p> Loading chart... </p>
			</div>);
		}
		else{
			console.log("Data available...");
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
				<div className="chart">
					<ReactECharts option={options} 
						theme={'macarons'} 
						style={{height: 'inherit', width: 'inherit'}}/>
				</div>
			);
		}
	}
};

export default Page;

