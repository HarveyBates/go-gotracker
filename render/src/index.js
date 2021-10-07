import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import 'mapbox-gl/dist/mapbox-gl.css';
import Activity from './views/activities';

ReactDOM.render(
  <React.StrictMode>
	<Activity />
  </React.StrictMode>,
  document.getElementById('root')
);

