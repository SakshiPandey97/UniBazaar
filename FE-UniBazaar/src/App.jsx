import React from "react";
import {
  BrowserRouter as Router,
} from "react-router-dom";
import Navbar from "./customComponents/Navbar";
import { AuthProvider } from "./hooks/useUserAuth";
import useModal from "./hooks/useModal";
import { SearchProvider } from "./context/SearchContext";
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
      <SearchProvider>
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
      </SearchProvider>
    </AuthProvider>
  );
}

export default App;