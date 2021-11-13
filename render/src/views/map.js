/* eslint import/no-webpack-loader-syntax: off */
import React from 'react';
import queryString from 'query-string';
import mapboxgl from '!mapbox-gl';
import './map.css';
import {InfluxDB} from '@influxdata/influxdb-client'

mapboxgl.accessToken = process.env.REACT_APP_MAPBOX;

export default class Map extends React.Component {
	constructor() {
		super();
		this.state = {
			records: [],
			cntrCoords: [],
			zoom: 13,
		};
		this.mapContainer = React.createRef();
	}

	async componentDidMount() {
		// Get params from url
		const params = queryString.parse(window.location.search);
		// Get get activity summary
		const reqUrl = `/activity/${params.sport}/${params.activity_id}`;
		const response = await fetch(reqUrl, {headers:{
			"Accept": "application/json",
			"Content-Type": "application/json"}});
		const data = await response.json();

		// Get latitude and longitude of activity
		const influxToken = process.env.REACT_APP_INFLUX_DB;
		const client = new InfluxDB({url: "http://localhost:8086", token: influxToken}).getQueryApi("user");
		const recordsQuery = `from(bucket: "records") |> range(start: time(v: ${data.start_time}), stop: time(v: ${data.end_time})) |> filter(fn: (r) => r["_measurement"] == "${data.activity_name}") |> filter(fn: (r) => r["_field"] == "position_lat" or r["_field"] == "position_long")`

		const handleState = () => {
			this.setState({
				records: records,
				cntrCoords: [data.swc_long, data.swc_lat]
			});

			const map = new mapboxgl.Map({
				container: this.mapContainer.current,
				style: 'mapbox://styles/mapbox/outdoors-v11',
				center: this.state.cntrCoords,
				zoom: this.state.zoom
			});

			var latitude = [];
			var longitude = [];
			for (let i in this.state.records) {
				var row = this.state.records[i];
				var currentField = row._field;
				if (currentField === "position_lat") {
					latitude.push(row._value);
				}
				else if(currentField === "position_long"){
					longitude.push(row._value);
				}
			}

			var coords = [];
			var pointCoords = []; 
			for (let i in latitude) {
				coords.push([longitude[i], latitude[i]]);
				var geoJson = {
					"type": "Feature",
					"geometry": {
						"type": "Point",
						"coordinates": [longitude[i], latitude[i]]
					},
					"id": parseFloat(i)
				}
				pointCoords.push(geoJson);
			}

			let hoverPointId = null;
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
				map.addSource('point-route', {
					'type': 'geojson',
					'data': {
						"type": "FeatureCollection",
						"features": pointCoords 
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
				map.addLayer({
					'id': 'route-points',
					'type': 'circle',
					'source': 'point-route',
					'paint': {
						'circle-color': '#ee6611',
						'circle-radius': {
							'base': 10,
							'stops': [
								[12, 8],
								[22, 6]
							]
						},
						'circle-opacity': [
							'case',
							['boolean', ['feature-state', 'hover'], false],
							1, 
							0	
						]
					}
				});
				map.on('mousemove', 'route-points', (e) => {
					if (e.features.length > 0) {
						if (hoverPointId !== null) {
							map.setFeatureState({
								source: 'point-route',
								id: hoverPointId
							}, 
							{
								hover: false
							});
						}
						hoverPointId = e.features[0].id;
						map.setFeatureState(
							{source: 'point-route', id: hoverPointId},
							{hover: true}
						);
					}
				});
				map.on('mouseleave', 'route-points', () => {
					if (hoverPointId !== null) {
						map.setFeatureState(
							{ source: 'point-route', id: hoverPointId },
							{ hover: false }
						);
					}
					hoverPointId = null;
				});


				const bounds = new mapboxgl.LngLatBounds(
					coords[0],
					coords[0]
				);
				for (const coord of coords) {
					bounds.extend(coord);
				}
				map.fitBounds(bounds, {
					padding: 20
				});

			});
		}

		var records = [];
		client.queryRows(recordsQuery, {
			next(row, tableMeta){
				var record = tableMeta.toObject(row);
				records.push(record);
			},
			error(error) {
				console.log(error);
			},
			complete() {
				handleState();
			},
		});
	}


	
		
	render() {
		if (this.state.records.length !== 0) {

			return (
				<div ref={this.mapContainer} className="map-container" />
			);
		} 
		else {
			return null
		}
	}

}
