import logo from './logo.svg';
import './App.css';
import Page from "./views/plot";

function App() {
	return (
		<div className="App">
			<div id="chart"></div>
			<header className="App-header">
				<img src={logo} className="App-logo" alt="logo" />
			</header>
		</div>
	);
}

export default App;
