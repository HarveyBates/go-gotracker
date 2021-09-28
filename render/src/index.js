import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import Activities from './views/activities';
import Map from './views/map';

ReactDOM.render(
  <React.StrictMode>
	<Map />
	<Activities />
  </React.StrictMode>,
  document.getElementById('root')
);

