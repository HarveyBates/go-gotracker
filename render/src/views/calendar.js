import React from 'react'
import FullCalendar from '@fullcalendar/react' 
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from "@fullcalendar/interaction"
import './calendar.css';

export default class Calendar extends React.Component {

	constructor() {
		super();
		this.state = {
			activities: []
		};
	}


	async componentDidMount() {
		try {
			var response = await fetch("/activities", {headers:{
				"Accept": "application/json",
				"Content-Type": "application/json"}}, {mode: "no-cors"});
			var data = await response.json();

			this.setState({
				activities: data
			});

		} catch (error) {
			console.log(error);
		}
	}


	render() {

		if (this.state.activities.length === 0) {
			return (
				<div>
					<p> Loading chart... </p>
				</div>
			);
		} 

		else {

			function formatName(str) {
				const name = str
					.toLowerCase()
					.replaceAll("_"," ")
					.replace("activity", "")
					.replace("generic", "outdoors")
					.replace("lap swimming", "laps")
					.replace(/[0-9]/g, "")
					.split(" ")
					.map(word => {
						return word.charAt(0).toUpperCase() + word.slice(1);
					})
					.join(" ");

				return name;
			}
			var events = [];
			for (let index in this.state.activities){
				var activity = this.state.activities[index];
				var name = formatName(activity.activity_name);
				var distance = (activity.total_distance / 1000).toFixed(2);
				var sport = activity.sport;
				var startDt = new Date(activity.start_time);
				events.push({
					title: name,
					id: activity.activity_id,
					date: startDt,
					sport: sport,
					distance: distance,
					elapsed_time: activity.total_elapsed_time
				});
			}

			return (
				<div className="calendar">
					<FullCalendar
						plugins={[ dayGridPlugin, interactionPlugin ]}
						initialView="dayGridMonth"
						initialDate="2021-01-01"
						selectable={true}
						events={events}  
						weekNumbers='true'
						eventContent={renderEventContent}
						eventClick={handleClick}
					/>
				</div>
			)
		}
	}
}

function handleClick(eventInfo) {
	console.log(eventInfo.event._def.publicId);
}



function renderEventContent(eventInfo) {

	function elapsedToTime(sec) {
		var date = new Date(null);
		date.setSeconds(sec);
		return date.toISOString().substr(11, 8);
	}	

	var props = eventInfo.event._def.extendedProps;
	var startTime = eventInfo.timeText + "m";

	if (props.sport === "cycling") {
		return (
			<div className="bike-cal-entry">
				<div className="bike-entry-head">
					<b>{eventInfo.event.title}</b>
				</div>
				<div className="bike-entry-body">
					<span><b>Distance: </b>{props.distance} km</span>
					<span><b>Start Time: </b>{startTime}</span>
					<span><b>Elapsed Time: </b>{elapsedToTime(props.elapsed_time)}</span>
				</div>
			</div>
		)
	} 
	else if (props.sport === "running") {
		return (
			<div className="run-cal-entry">
				<div className="run-entry-head">
					<b>{eventInfo.event.title}</b>
				</div>
				<div className="run-entry-body">
					<span><b>Distance: </b>{props.distance} km</span>
					<span><b>Start Time: </b>{startTime}</span>
					<span><b>Elapsed Time: </b>{elapsedToTime(props.elapsed_time)}</span>
				</div>
			</div>
		)
	} 
	else if (props.sport === "swimming") {
		return (
			<div className="swim-cal-entry">
				<div className="swim-entry-head">
					<b>{eventInfo.event.title}</b>
				</div>
				<div className="swim-entry-body">
					<span><b>Distance: </b>{props.distance} km</span>
					<span><b>Start Time: </b>{startTime}</span>
					<span><b>Elapsed Time: </b>{elapsedToTime(props.elapsed_time)}</span>
				</div>
			</div>
		)
	} 
	else {
		return (
			<>
				<b className="cal-entry">{eventInfo.event.title}</b>
			</>
		)
	}
}

