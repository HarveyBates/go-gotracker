import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import 'mapbox-gl/dist/mapbox-gl.css';
import Activities from './views/activities';

ReactDOM.render(
  <React.StrictMode>
	<Activities />
  </React.StrictMode>,
  document.getElementById('root')
);

