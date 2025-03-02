import { useState } from "react";
import { useUserAuth } from "./useUserAuth";
import { useNavigate } from "react-router-dom";

const useNavbar = ({ toggleModal }) => {
  const [isMenuOpen, setMenuOpen] = useState(false);
  const [isDropdownOpen, setDropdownOpen] = useState(false);

  const userAuth = useUserAuth();
  const navigate = useNavigate();

  const toggleDropdown = () => {
    setDropdownOpen(!isDropdownOpen);
  };
  const toggleMenu = () => setMenuOpen((prev) => !prev);

  const handleNavigation = (path) => {
    userAuth.userState ? navigate(path) : toggleModal();
  };

  const handleAuthAction = () => {
    if (userAuth.userState) {
      userAuth.toggleUserLogin();
      userAuth.setUserID("");
    } else {
      toggleModal();
    }
  };

  return {
    isMenuOpen,
    isDropdownOpen,
    toggleDropdown,
    toggleMenu,
    handleNavigation,
    handleAuthAction,
    userAuth,
  };
};

export default useNavbar;
