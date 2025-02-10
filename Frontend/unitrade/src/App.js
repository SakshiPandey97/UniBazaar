import React, { useEffect } from "react";
import Navbar from "./pages/components/Navbar";
import Banner from "./pages/components/Banner";
import Products from "./pages/components/Products";
import SellProduct from "./pages/SellProduct";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  useLocation,
} from "react-router-dom";
import "./App.css";
import Auth from "./pages/Auth";
import { AuthProvider } from "./pages/components/useUserAuth";

function Layout() {
  return (
    <div className="App">
      {/* <Navbar /> */}
      <Banner />
      <Products />
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      
      <Router>
        <ScrollLockOnLogin />
        <Navbar /> {/* Keep Navbar outside to persist across pages */}
        <Routes>
          <Route path="/" element={<Layout />} />
          <Route path="/sell" element={<SellProduct />} />
          <Route path="/auth/*" element={<Auth />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

function ScrollLockOnLogin() {
  const location = useLocation();

  useEffect(() => {
    if (location.pathname.includes("/auth")) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "auto";
    }
  }, [location.pathname]);

  return null;
}

export default App;
