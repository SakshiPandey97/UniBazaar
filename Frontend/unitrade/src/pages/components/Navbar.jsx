import React from "react";
import { Link } from "react-router-dom";
import useNavbar from "../../hooks/useNavBar";

import logo from "../../assets/imgs/logo.svg";
import loginIcon from "../../assets/imgs/login.svg";
import menuToggleIcon from "../../assets/imgs/menu-toggle.svg";
import closeIcon from "../../assets/imgs/close.svg";

const Navbar = ({ toggleModal }) => {
  const {
    isMenuOpen,
    isDropdownOpen,
    toggleDropdown,
    toggleMenu,
    handleNavigation,
    handleAuthAction,
    userAuth,
  } = useNavbar({ toggleModal });

  return (
    <nav className="w-full h-[80px] bg-[#F8F8F8] border-b py-2">
      <div className="md:max-w-[1480px] max-w-[600px] mx-auto flex justify-between items-center">
        <Link to="/" className="relative inline-block">
          <img
            src={logo}
            className="h-[50px] md:h-[60px] w-auto transition hover:opacity-75"
            alt="Logo"
          />
          <div className="absolute inset-0 bg-blue-500 opacity-0 hover:opacity-30 transition"></div>
        </Link>

        <ul className="hidden md:flex gap-20 font-[serif] text-[32px] font-medium">
          <li>
            <Link to="/">Home</Link>
          </li>
          <li
            className="cursor-pointer "
            onClick={() => handleNavigation("/buy")}
          >
            Buy
          </li>
          <li
            className="cursor-pointer "
            onClick={() => handleNavigation("/sell")}
          >
            Sell
          </li>
          <li
            className="cursor-pointer "
            onClick={() => handleNavigation("/product")}
          >
            Products
          </li>
          <li>
            <Link to="/about">About Us</Link>
          </li>
        </ul>
        <div className="flex flex-row hidden md:flex items-center gap-4 relative">
          <button
            className="flex flex-row items-center gap-2 px-4 py-3 rounded-md bg-[#008080] hover:bg-[#006666] transition duration-200 text-white text-[24px] font-medium relative"
            onClick={userAuth.userState ? toggleDropdown : toggleModal}
          >
            <img src={loginIcon} className="h-[24px]" alt="Login" />
            <span>{userAuth.userState ? "Profile" : "Login"}</span>
          </button>

          {userAuth.userState && isDropdownOpen && (
            <div className="absolute top-full mt-2 right-0 w-48 bg-white rounded-md shadow-lg border border-gray-300 z-50">
              <ul className="py-2">
                <li
                  className="px-4 py-2 hover:bg-gray-200 cursor-pointer"
                  onClick={() => console.log("View Profile Clicked")}
                >
                  View My Profile
                </li>
                <li
                  className="px-4 py-2 hover:bg-gray-200 cursor-pointer text-red-500"
                  onClick={() => {
                    handleAuthAction();
                  }}
                >
                  Log Out
                </li>
              </ul>
            </div>
          )}
        </div>
        <button className="md:hidden" onClick={toggleMenu}>
          <img src={isMenuOpen ? closeIcon : menuToggleIcon} alt="Menu" />
        </button>
      </div>

      {isMenuOpen && (
        <div className="absolute z-10 p-4 bg-white w-full px-8 md:hidden">
          <ul>
            <li
              className="p-4 hover:bg-gray-100"
              onClick={() => handleNavigation("/buy")}
            >
              Buy
            </li>
            <li
              className="p-4 hover:bg-gray-100"
              onClick={() => handleNavigation("/sell")}
            >
              Sell
            </li>
            <li
              className="p-4 hover:bg-gray-100"
              onClick={() => handleNavigation("/product")}
            >
              Products
            </li>
            <li className="p-4 hover:bg-gray-100">
              <Link to="/about">About Us</Link>
            </li>
          </ul>
          <div className="flex flex-col my-4 gap-4">
            <button
              className="border border-[#008080] flex items-center justify-center bg-transparent px-6 gap-2 py-3"
              onClick={handleAuthAction}
            >
              <img src={loginIcon} className="h-[20px]" alt="Login" />
              <span>{userAuth.userState ? "Log-Out" : "Login"}</span>
            </button>
            <button className="px-4 py-5 rounded-md bg-[#008080] text-white">
              Transactions
            </button>
          </div>
        </div>
      )}
    </nav>
  );
};

export default Navbar;
