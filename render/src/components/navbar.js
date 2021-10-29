// React imports 
import React from 'react';
import { BrowserRouter as Router, Route, Switch, Link, withRouter } from "react-router-dom";

// User imports 
import "./navbar.css";

// Logo / icon imports
import { library } from '@fortawesome/fontawesome-svg-core';
import { fas } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
library.add(fas);

export class Navbar extends React.Component {
	render() {
		return (
			<Router>
			<nav className="navbar">
				<ul className="navlist"> 
					<Link to="#" className="navitem">
						<FontAwesomeIcon className="navimg" icon="fa-solid fa-home"/>
					</Link>
					<Link to="#" className="navitem">
						<FontAwesomeIcon className="navimg" icon="fa-solid fa-calendar-days"/>
					</Link>

					<Link className="navitem">
						<FontAwesomeIcon className="navimg" icon="fa-solid fa-heart-pulse"/>
					</Link>

					<Link className="navitem">
						<FontAwesomeIcon className="navimg" icon={["fas", "chart-line"]}/>
					</Link>

					<Link className="navitem">
						<FontAwesomeIcon className="navimg" icon="fa-solid fa-battery-three-quarters"/>
					</Link>
				</ul>
			</nav>
			</Router>
		);
	}
}

