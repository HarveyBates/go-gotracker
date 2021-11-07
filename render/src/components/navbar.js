// React imports 
import React from 'react';
import { BrowserRouter, Routes, Route, Switch, Link, withRouter } from "react-router-dom";
// User imports 
import "./navbar.css";
import Activities from "../views/activities";
import Calendar from "../views/calendar";

// Logo / icon imports
import { library } from '@fortawesome/fontawesome-svg-core';
import { fas } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
library.add(fas);


export class Navbar extends React.Component {
	render() {
		return (
			<BrowserRouter>
			<nav className="navbar">
				<ul className="navlist"> 
					<Link to="#" className="navitem">
						<FontAwesomeIcon className="navimg" icon="fa-solid fa-home"/>
					</Link>
					<Link to="/calendar" className="navitem">
						<FontAwesomeIcon className="navimg" icon="fa-solid fa-calendar-days"/>
					</Link>

					<Link to="/activities" className="navitem">
						<FontAwesomeIcon className="navimg" icon="fa-solid fa-heart-pulse"/>
					</Link>

					<Link to="#" className="navitem">
						<FontAwesomeIcon className="navimg" icon={["fas", "chart-line"]}/>
					</Link>

					<Link to="#" className="navitem">
						<FontAwesomeIcon className="navimg" icon="fa-solid fa-battery-three-quarters"/>
					</Link>
				</ul>
			</nav>
				<Routes>
					<Route path="/activities" element={<Activities activity_id={4892341954} sport="running"/>}/>
					<Route path="/calendar" element={<Calendar />}/>
				</Routes>
			</BrowserRouter>
		);
	}
}

