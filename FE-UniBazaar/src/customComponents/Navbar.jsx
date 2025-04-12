import React from "react";
import { Link, useLocation } from "react-router-dom";
import useNavbar from "../hooks/useNavBar";
import loginIcon from "../assets/imgs/login.svg";
import menuToggleIcon from "../assets/imgs/menu-toggle.svg";
import closeIcon from "../assets/imgs/close.svg";
import Logo from "./Logo";

const Navbar = ({ toggleLoginModal }) => {
  const location = useLocation();
  const {
    isMenuOpen,
    isDropdownOpen,
    toggleDropdown,
    toggleMenu,
    handleNavigation,
    handleLogout,
    userAuth,
    } = useNavbar({ toggleLoginModal });

  const isActive = (path) => location.pathname === path;

  return (
    <>
      <nav className="fixed top-0 left-0 w-full h-[80px] bg-gradient-to-r from-[#002855] via-[#014F86] to-[#032B54] border-b py-2 z-50 shadow-md md:px-8 px-4 flex justify-between items-center">
        <div className="flex items-center">
          <Link to="/" className="relative inline-block">
            <Logo />
          </Link>
        </div>

        <div className="hidden md:flex gap-8 text-lg font-semibold">
          {[
            { path: "/", label: "Home" },
            { path: "#", label: "Messaging", onClick: () => handleNavigation("/messaging") },
            { path: "#", label: "Sell", onClick: () => handleNavigation("/sell") },
            { path: "/products", label: "Products" },
            {
              path: "/userproducts",
              label: "My  Products",
              onClick: () => handleNavigation("/userproducts"),
            },
            { path: "/about", label: "About Us" },
          ].map(({ path, label, onClick }) =>
            onClick ? (
              <button
                key={label}
                onClick={onClick}
                className="px-4 py-2 rounded-lg text-[#E5E5E5] hover:bg-[#FFC67D] hover:text-black transition-all"
              >
                {label}
              </button>
            ) : (
              <Link
                key={path}
                to={path}
                className={`px-4 py-2 rounded-lg transition-all duration-300 ${isActive(path)
                    ? "bg-[#FFC67D] text-black font-bold"
                    : "text-[#E5E5E5] hover:bg-[#FFC67D] hover:text-black"
                  }`}
              >
                {label}
              </Link>
            )
          )}
        </div>

        <div className="hidden md:flex items-center gap-4 relative">
          <button
            data-testid="loginBtn"
            className="flex items-center gap-2 px-6 py-3 rounded-lg bg-[#E5E5E5] hover:bg-[#D6D2D2] transition duration-200 text-black font-medium"
            onClick={userAuth.userState ? toggleDropdown : toggleLoginModal}
          >
            <img src={loginIcon} className="h-[24px]" alt="Login" />
            {userAuth.userState ? "Profile" : "Login"}
          </button>

          {userAuth.userState && isDropdownOpen && (
            <div className="absolute top-full mt-2 right-0 w-48 bg-white rounded-lg shadow-lg border border-gray-300 z-50">
              <ul className="py-2">
                <li
                  className="px-4 py-2 hover:bg-gray-200 cursor-pointer"
                  onClick={toggleLoginModal}
                >
                  View My Profile
                </li>
                <li
                  className="px-4 py-2 hover:bg-gray-200 cursor-pointer text-red-500"
                  onClick={handleLogout}
                >
                  Log Out
                </li>
              </ul>
            </div>
          )}
        </div>

        <button className="md:hidden" onClick={toggleMenu}>
          <img
            src={isMenuOpen ? closeIcon : menuToggleIcon}
            alt="Menu"
            className="h-8 w-8"
          />
        </button>
      </nav>

      {isMenuOpen && (
        <div className="absolute top-[80px] left-0 w-full bg-white px-6 py-4 md:hidden shadow-lg z-50">
          <ul className="space-y-4">
            {[

              { path: "#", label: "Messaging", onClick: () => handleNavigation("/messaging") },
              { path: "#", label: "Sell", onClick: () => handleNavigation("/sell") },
              { path: "/product", label: "Products" },
              {
                path: "/userproducts",
                label: "My Products",
                onClick: () => handleNavigation("/userproducts"),
              },
              { path: "/about", label: "About Us" },
            ].map(({ path, label, onClick }) => (
              <li
                key={path}
                className={`p-4 rounded-lg transition ${isActive(path)
                    ? "bg-[#FFC67D] text-black font-bold"
                    : "hover:bg-gray-100"
                }`}
                onClick={onClick || (() => handleNavigation(path))}
              >
                {label}
              </li>
            ))}
          </ul>

          <div className="mt-4 flex flex-col gap-4">
            <button
              className="border border-[#008080] flex items-center justify-center bg-transparent px-6 gap-2 py-3 rounded-lg"
              onClick={userAuth.userState ? handleLogout : toggleLoginModal}
            >
              <img src={loginIcon} className="h-[20px]" alt="Login" />
              {userAuth.userState ? "Log-Out" : "Login"}
            </button>
          </div>
        </div>
      )}

      <div className="pt-[80px]"></div>
    </>
  );
};

export default Navbar;
