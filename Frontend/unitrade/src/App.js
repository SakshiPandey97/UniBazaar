import React, { useState, useEffect } from "react";
import Navbar from "./pages/components/Navbar";
import Banner from "./pages/components/Banner";
import Products from "./pages/components/Products";
import LoginPage from "./pages/LoginPage";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import "./App.css";

function App() {
  const [isLoginClicked, setLoginClicked] = useState(false);
  useEffect(() => {
    if (isLoginClicked) {
      document.body.style.overflow = "hidden"; // Disable scrolling
    } else {
      document.body.style.overflow = "auto"; // Enable scrolling
    }
  }, [isLoginClicked]);

  const handleLoginButtonClick = () => {
    console.log(isLoginClicked);
    setLoginClicked(!isLoginClicked);
  };
  return (
    <Router>
      <div className="App">
        <Navbar handleLoginButtonClick={handleLoginButtonClick} />
        <Routes>
          <Route
            path="/"
            element={
              <>
                {/* <LandingPage /> */}
                <Banner />
                <Products />
              </>
            }
          />
        </Routes>
        {isLoginClicked && (
          <div className="modal-overlay" onClick={handleLoginButtonClick}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
              <LoginPage/>
            </div>
          </div>
        )}
      </div>
    </Router>
  );
}

export default App;
