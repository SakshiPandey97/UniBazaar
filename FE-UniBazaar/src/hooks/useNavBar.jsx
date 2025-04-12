import { useState } from "react";
import { useUserAuth } from "./useUserAuth";
import { useNavigate } from "react-router-dom";

const useNavbar = ({ toggleLoginModal }) => {
  const [isMenuOpen, setMenuOpen] = useState(false);
  const [isDropdownOpen, setDropdownOpen] = useState(false);

  const userAuth = useUserAuth();
  const navigate = useNavigate();

  const toggleDropdown = () => {
    setDropdownOpen((prev) => !prev);
  };

  const toggleMenu = () => {
    setMenuOpen((prev) => !prev);
  };

  const handleNavigation = (path) => {
    userAuth.userState ? navigate(path) : toggleLoginModal();
  };

  const handleLogout = () => {
    userAuth.logoutUser();
  };

  return {
    isMenuOpen,
    isDropdownOpen,
    toggleDropdown,
    toggleMenu,
    handleNavigation,
    handleLogout,
    userAuth,
  };
};

export default useNavbar;
