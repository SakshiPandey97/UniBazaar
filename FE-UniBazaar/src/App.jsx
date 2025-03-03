import React, { lazy, Suspense } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Navbar from "./customComponents/Navbar";
import Banner from "./customComponents/Banner";
import Spinner from "./customComponents/Spinner";
import { AuthProvider } from "./hooks/useUserAuth";
import useModal from "./hooks/useModal";
import "./App.css";
import AuthPage from "./pages/AuthPage";
import ViewMyProfilePage from "./pages/ViewMyProfilePage";
const SellProductPage = lazy(() => import("./pages/SellProductPage"));
const Products = lazy(() => import("./customComponents/Products"));

function Layout() {
  return (
    <div className="App bg-[#D6D2D2]">
      <Banner />
      <Suspense fallback={<Spinner />}>
        <Products />
      </Suspense>
    </div>
  );
}

function App() {
  const { isModalOpen: isProfileModalOpen, toggleModal: toggleProfileModal } = useModal();
  const { isModalOpen: isLoginModalOpen, toggleModal: toggleLoginModal } = useModal();

  return (
    <AuthProvider>
      {isLoginModalOpen && (
        <div className="modal-overlay" onClick={toggleLoginModal}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <AuthPage toggleModal={toggleLoginModal} />
          </div>
        </div>
      )}

      {isProfileModalOpen && (
        <div className="modal-overlay" onClick={toggleProfileModal}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <ViewMyProfilePage />
          </div>
        </div>
      )}

      <Router>
        <Navbar toggleLoginModal={toggleLoginModal} toggleViewProfile={toggleProfileModal} />
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
