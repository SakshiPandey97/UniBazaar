import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Navbar from "./pages/components/Navbar";
import Banner from "./pages/components/Banner";
import Products from "./pages/components/Products";
import SellProduct from "./pages/SellProduct";
import { AuthProvider } from "./hooks/useUserAuth";
import useLoginModal from "./hooks/useLoginModal"; 
import "./App.css";
import AuthPage from "./pages/AuthPage";

function Layout() {
  return (
    <div className="App">
      <Banner />
      <Products />
    </div>
  );
}

function App() {
  const { isModalOpen, toggleModal } = useLoginModal();

  return (
    <AuthProvider>
      {isModalOpen && (
        <div className="modal-overlay" onClick={toggleModal}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <AuthPage toggleModal={toggleModal}/>
          </div>
        </div>
      )}

      <Router>
        <Navbar toggleModal={toggleModal} />
        <Routes>
          <Route path="/" element={<Layout />} />
          <Route path="/sell" element={<SellProduct />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
