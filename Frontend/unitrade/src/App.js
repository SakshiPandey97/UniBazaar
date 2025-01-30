import React from 'react';
import LandingPage from './pages/LandingPage';
import Navbar from './pages/components/Navbar';
import Banner from './pages/components/Banner';
import Products from './pages/components/Products';
import LoginPage from './pages/LoginPage';
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import './App.css';

function App() {
  return (
    <Router>
    <div className="App">
      <Navbar />
        <Routes>
          <Route path="/" element={
            <>
            {/* <LandingPage /> */}
            <Banner />
            <Products/>
            </>
            } />

          <Route path="/login" element={<LoginPage/>} />
        </Routes>
    </div>
  </Router>
  );
}

export default App;
