import React, { lazy, Suspense } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Navbar from "./pages/components/Navbar";
import Banner from "./pages/components/Banner";
import Spinner from "./pages/components/Spinner";
import { AuthProvider } from "./hooks/useUserAuth";
import useLoginModal from "./hooks/useLoginModal";
import "./App.css";
import AuthPage from "./pages/AuthPage";
const SellProductPage = lazy(() => import("./pages/SellProduct"));
const Products = lazy(() => import("./pages/components/Products"));

function Layout() {
  return (
    <div className="App">
      <Banner />
      <Suspense fallback={<Spinner />}>
        <Products />
      </Suspense>
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
            <AuthPage toggleModal={toggleModal} />
          </div>
        </div>
      )}

      <Router>
        <Navbar toggleModal={toggleModal} />
        <Routes>
          <Route path="/" element={<Layout />} />
          <Route
            path="/sell"
            element={
              <Suspense fallback={<Spinner />}>
                <SellProductPage />
              </Suspense>
            }
          />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
