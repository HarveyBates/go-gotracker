import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import 'mapbox-gl/dist/mapbox-gl.css';
import { Navbar } from './components/navbar.js';

ReactDOM.render(
  <React.StrictMode>
	<Navbar />
  </React.StrictMode>,
  document.getElementById('root')
);

