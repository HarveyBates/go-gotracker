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
			lat: -33.335,
			lng: 150.521273,
			zoom: 9,
		};
		this.mapContainer = React.createRef();
		this.polyline = "";
	}

	async componentDidMount() {
		const response = await fetch('/activity/latest', {headers:{
			"Accept": "application/json",
			"Content-Type": "application/json"}});
		const data = await response.json();

		this.polyline = data.Summary.map.summary_polyline;
		console.log(this.polyline);
		if (this.polyline != "") {
			const decodedPoly = polyline.decode(this.polyline);

			function flip(coords) {
				var flipped = [];
				var sumLat = 0;
				var sumLng = 0;
				for (var i = 0; i < coords.length; i++) {
					var coord = coords[i].slice();
					flipped.push([coord[1], coord[0]]);
					sumLat += coord[1];
					sumLng += coord[0]
				}
				var avLat = sumLat / coords.length;
				var avLong = sumLng / coords.length;

				var latlng = [avLat, avLong];

				return {
					flipped: flipped,
					av: latlng
				};
			}

			const coords = flip(decodedPoly);

			const { lng, lat, zoom } = this.state;
			const map = new mapboxgl.Map({
				container: this.mapContainer.current,
				style: 'mapbox://styles/mapbox/outdoors-v11',
				center: coords.av,
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
							'coordinates': coords.flipped
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
	}
	
		
	render() {
		if (this.polyline != null) {
			return (
					<div ref={this.mapContainer} className="map-container" />
			);
		}
	}

}
