/* eslint import/no-webpack-loader-syntax: off */
import React from 'react';
import mapboxgl from '!mapbox-gl';
import './map.css';
import polyline from '@mapbox/polyline';

mapboxgl.accessToken = process.env.REACT_APP_MAPBOX;

export default class Map extends React.Component {
	constructor() {
		super();
		this.state = {
			lat: -33.2835,
			lng: 149.101273,
			zoom: 13,
		};
		this.mapContainer = React.createRef();
		this.polyline = "";
	}

	async componentDidMount() {
		const response = await fetch('/activity/stream/6010261290', {headers:{
			"Accept": "application/json",
			"Content-Type": "application/json"}});
		const data = await response.json();

		this.polyline = data.Attributes.map.summary_polyline;
		const decodedPoly = polyline.decode(this.polyline);

		function flip(coords) {
			var flipped = [];
			for (var i = 0; i < coords.length; i++) {
				var coord = coords[i].slice();
				flipped.push([coord[1], coord[0]]);
			}
			return flipped
		}

		const coords = flip(decodedPoly);

		const { lng, lat, zoom } = this.state;
		const map = new mapboxgl.Map({
			container: this.mapContainer.current,
			style: 'mapbox://styles/mapbox/outdoors-v11',
			center: coords[0],
			zoom: zoom
		});
		
		map.on('load', () => {
			map.addSource('route', {
				'type': 'geojson',
				'data': {
					'type': 'Feature',
					'properties': {},
					'geometry': {
						'type': 'LineString',
						'coordinates': coords
					}
				}
			});
			map.addLayer({
				'id': 'route',
				'type': 'line',
				'source': 'route',
				'layout': {
					'line-join': 'round',
					'line-cap': 'round'
				},
				'paint': {
					'line-color': '#001f3f',
					'line-width': 3
				}
			});
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
