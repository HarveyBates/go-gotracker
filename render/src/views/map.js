/* eslint import/no-webpack-loader-syntax: off */
import React from 'react';
import 'mapbox-gl/dist/mapbox-gl.css';
import mapboxgl from '!mapbox-gl';

mapboxgl.accessToken = process.env.REACT_APP_MAPBOX;

class Map extends React.PureComponent {
	constructor(props) {
		super(props);
		this.state = {
			lng: -70.9,
			lat: 42.35,
			zoom: 9
		};
		this.mapContainer = React.createRef();
	}

	componentDidMount() {
		const { lng, lat, zoom } = this.state;
		const map = new mapboxgl.Map({
			container: this.mapContainer.current,
			style: 'mapbox://styles/mapbox/streets-v11',
			center: [lng, lat],
			zoom: zoom
		});
	}
		
	render() {
		return (
			<div>
				<div ref={this.mapContainer} className="map-container" />
			</div>
		);
	}

}

export default Map;
