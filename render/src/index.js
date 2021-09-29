import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import 'mapbox-gl/dist/mapbox-gl.css';
import Activities from './views/activities';
import Map from './views/map';

ReactDOM.render(
  <React.StrictMode>
	<Activities />
	<Map />
  </React.StrictMode>,
  document.getElementById('root')
);

