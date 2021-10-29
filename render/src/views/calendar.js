import React from 'react'
import FullCalendar from '@fullcalendar/react' 
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from "@fullcalendar/interaction"
import './calendar.css';

export default class Calendar extends React.Component {
  render() {
	  return (
		  <div className="calendar">
			  <FullCalendar
				  plugins={[ dayGridPlugin, interactionPlugin ]}
				  initialView="dayGridMonth"
				  events={[
					{ title: 'Morning Run', 
				  		date: '2021-10-30' },
				  ]}  
				  eventContent={renderEventContent}
			  />
		  </div>
    )
  }
}
function renderEventContent(eventInfo) {
  return (
	<>
		<b> &#x1f3c3; {eventInfo.event.title}</b>
	</>
  )
}
