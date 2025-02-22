import React, { lazy, Suspense } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Navbar from "./customComponents/Navbar";
import Banner from "./customComponents/Banner";
import Spinner from "./customComponents/Spinner";
import { AuthProvider } from "./hooks/useUserAuth";
import useLoginModal from "./hooks/useLoginModal";
import "./App.css";
import AuthPage from "./pages/AuthPage";
import ViewMyProfilePage from "./pages/ViewMyProfilePage";
const SellProductPage = lazy(() => import("./pages/SellProductPage"));
const Products = lazy(() => import("./customComponents/Products"));

function Layout() {
  return (
    <div className="App bg-gray-200">
      <Banner />
      <Suspense fallback={<Spinner />}>
        <Products />
      </Suspense>
    </div>
  );
}

function App() {
  const { isModalOpen, isViewProfileOpen, openViewMyProfile, toggleModal } =
    useLoginModal();

  return (
    <AuthProvider>
      {isModalOpen && (
        <div className="modal-overlay" onClick={toggleModal}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <AuthPage toggleModal={toggleModal} />
          </div>
        </div>
      )}

      {isViewProfileOpen && (
        <div className="modal-overlay" onClick={openViewMyProfile}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <ViewMyProfilePage />
          </div>
        </div>
      )}

      <Router>
        <Navbar toggleModal={toggleModal} toggleViewProfile={openViewMyProfile} />
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
