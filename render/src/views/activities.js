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
			activityID: 0,
			type: ""
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
			this.setState({name: data.Name, distance: data.Distance, date: data.StartDate, activityID: data.ID, type: data.Type});

			// Get the stream of this activity
			response = await fetch('/activity/' + this.state.activityID + "/stream", {headers:{
				"Accept": "application/json",
				"Content-Type": "application/json"}});
			data = await response.json();
			this.setState({data: data});
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
			var swimStream = [];
			if(this.state.type == "Swim"){
				var timeArr = this.state.data.Time.data;
				var distArr = this.state.data.Distance.data;
				var strokeArr = this.state.data.Cadence.data;
				var paceArr = this.state.data.VelocitySmooth.data;

				var stroke = [];
				var pace = [];
				var pacePer100m = [];
				if (timeArr != null && distArr != null) {
					var prevTSplit = 0;
					var prevDVal = 0;
					for (let i in timeArr) {
						var distSplit = Math.floor(distArr[i]);	
						if(distSplit % 100 == 0 && prevDVal != distSplit){
							var timeSplit = Math.floor(timeArr[i]);
							var adjTSplit = timeSplit - prevTSplit;
							prevTSplit = timeSplit;
							prevDVal = distSplit;
							var mins = Math.floor(adjTSplit / 60);
							var sec = adjTSplit - mins * 60;
							if(distSplit != 0){
								console.log(distSplit, mins + ":" + sec);
							}
						}
					}
				}


					//
					//for (let i in distArr){
					//	if(strokeArr != null) {
					//		stroke.push([distArr[i], strokeArr[i]]);
					//	}
					//	if(paceArr != null) {
					//		pace.push([distArr[i], paceArr[i]]);
					//	}
					//}
				//}
				swimStream = [{
					name: "Stroke Rate (rpm)",
					color: 'rgba(0, 0, 0, 1)',
					type: 'line',
					step: 'start',
					symbol: 'none',
					data: stroke,
				}, 
					{
					name: "Pace",
					color: 'rgba(50, 150, 255, 1)',
					type: 'line',
					step: 'start',
					yAxisIndex: 1,
					symbol: 'none',
					data: pace
				}];
				//console.log(swimStream)
			}
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
				if (heartrate != null) {
					hr.push([dMeters, heartrate[i]]);
				}
				rpm.push([dMeters, cadence[i]]);
				if (watts != null) {
					power.push([dMeters, watts[i]]);
				}
				if (altitude != null) {
					alt.push([dMeters, altitude[i]]);
				}

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
					swimStream[0],
					swimStream[1],
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

