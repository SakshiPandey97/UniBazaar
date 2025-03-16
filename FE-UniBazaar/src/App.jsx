import React, { lazy, Suspense } from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  useLocation,
} from "react-router-dom";
import Navbar from "./customComponents/Navbar";
import { AuthProvider } from "./hooks/useUserAuth";
import useModal from "./hooks/useModal";
import "./App.css";
import AuthPage from "./pages/AuthPage";
import ViewMyProfilePage from "./pages/ViewMyProfilePage";
import Modal from "./customComponents/Modal";
import AnimatedRoutes from "./customComponents/AnimatedRoutes";

function App() {
  const { isModalOpen: isProfileModalOpen, toggleModal: toggleProfileModal } =
    useModal();
  const { isModalOpen: isLoginModalOpen, toggleModal: toggleLoginModal } =
    useModal();

  return (
    <AuthProvider>
      <Modal isOpen={isLoginModalOpen} toggleModal={toggleLoginModal}>
        <AuthPage toggleModal={toggleLoginModal} />
      </Modal>

      <Modal isOpen={isProfileModalOpen} toggleModal={toggleProfileModal}>
        <ViewMyProfilePage />
      </Modal>

      <Router>
        <Navbar
          toggleLoginModal={toggleLoginModal}
          toggleViewProfile={toggleProfileModal}
        />
        <AnimatedRoutes />
      </Router>
    </AuthProvider>
  );
}

export default App;