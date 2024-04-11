import React from 'react';
import { Routes, Route } from 'react-router-dom';
import NavBar from './components/navigation/NavBar';
import HomePage from './components/pages/HomePage';
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
	return (
		<div id='App'>
			<NavBar />
			<Routes>
				<Route
					path='/'
					element={<HomePage />}
				/>
			</Routes>
		</div>
	);
}

export default App;
