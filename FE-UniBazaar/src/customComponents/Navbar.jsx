import React, { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import useNavbar from "../hooks/useNavBar";
import logo from "../assets/imgs/logo.svg";
import loginIcon from "../assets/imgs/login.svg";
import menuToggleIcon from "../assets/imgs/menu-toggle.svg";
import closeIcon from "../assets/imgs/close.svg";
import useLoginModal from "@/hooks/useModal";

const Navbar = ({ toggleLoginModal}) => {
  const location = useLocation();
  const {
    isMenuOpen,
    isDropdownOpen,
    toggleDropdown,
    toggleMenu,
    handleNavigation,
    handleAuthAction,
    userAuth,
  } = useNavbar({ toggleLoginModal });

  const isActive = (path) => location.pathname === path;

  return (
    <>
      <nav className="fixed top-0 left-0 w-full h-[80px] bg-[#032B54] border-b py-2 z-50 shadow-lg">
        <div className="flex flex-row justify-between md:max-w-[1480px] max-w-[600px] mx-auto items-center">
          <div className="flex w-1/10">
            <Link to="/" className="relative inline-block">
              <img
                src={logo}
                className="h-[50px] md:h-[60px] w-auto transition hover:opacity-75"
                alt="Logo"
              />
            </Link>
          </div>

          <div className="flex w-7/10 justify-evenly">
            <ul className="hidden md:flex font-[serif] text-[32px] font-medium gap-6">
              <li>
                <Link
                  to="/"
                  className={`text-[#E5E5E5] ${
                    isActive("/") ? "text-black font-bold bg-[#FFC67D] " : ""
                  } hover:bg-[#FFC67D] hover:text-black`}
                >
                  Home
                </Link>
              </li>
              <li
                className={`cursor-pointer text-[#E5E5E5] ${
                  isActive("/buy") ? "text-black font-bold bg-[#FFC67D]" : ""
                } hover:bg-[#FFC67D] hover:text-black`}
                onClick={() => handleNavigation("/buy")}
              >
                Buy
              </li>
              <li
                className={`cursor-pointer text-[#E5E5E5] ${
                  isActive("/sell") ? "text-black font-bold bg-[#FFC67D] " : ""
                } hover:bg-[#FFC67D] hover:text-black`}
                onClick={() => handleNavigation("/sell")}
              >
                Sell
              </li>
              <li
                className={`cursor-pointer text-[#E5E5E5] ${
                  isActive("/product") ? "text-black font-bold bg-[#FFC67D]" : ""
                } hover:bg-[#FFC67D] hover:text-black`}
                onClick={() => handleNavigation("/product")}
              >
                Products
              </li>
              <li>
                <Link
                  to="/about"
                  className={`text-[#E5E5E5] ${
                    isActive("/about") ? "text-black font-bold bg-[#FFC67D]" : ""
                  } hover:bg-[#FFC67D] hover:text-black`}
                >
                  About Us
                </Link>
              </li>
            </ul>
          </div>

          <div className="flex flex-row hidden md:flex items-center gap-4 relative">
            <button
              className="flex flex-row items-center gap-2 px-4 py-3 rounded-md bg-[#E5E5E5] hover:bg-[#D6D2D2] transition duration-200 text-white text-[24px] font-medium relative"
              onClick={userAuth.userState ? toggleDropdown : toggleLoginModal}
            >
              <img src={loginIcon} className="h-[24px]" alt="Login" />
              <span className="text-black">{userAuth.userState ? "Profile" : "Login"}</span>
            </button>

            {userAuth.userState && isDropdownOpen && (
              <div className="absolute top-full mt-2 right-0 w-48 bg-white rounded-md shadow-lg border border-gray-300 z-50">
                <ul className="py-2">
                  <li
                    className="px-4 py-2 hover:bg-gray-200 cursor-pointer"
                    onClick={toggleLoginModal}
                  >
                    View My Profile
                  </li>
                  <li
                    className="px-4 py-2 hover:bg-gray-200 cursor-pointer text-red-500"
                    onClick={handleAuthAction}
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
          <div className="absolute z-500 p-4 bg-white w-full px-8 md:hidden shadow-lg">
            <ul>
              <li
                className={`p-4 hover:bg-gray-100 ${
                  isActive("/buy") ? "text-teal-600 font-bold" : ""
                }`}
                onClick={() => handleNavigation("/buy")}
              >
                Buy
              </li>
              <li
                className={`p-4 hover:bg-gray-100 ${
                  isActive("/sell") ? "text-teal-600 font-bold" : ""
                }`}
                onClick={() => handleNavigation("/sell")}
              >
                Sell
              </li>
              <li
                className={`p-4 hover:bg-gray-100 ${
                  isActive("/product") ? "text-teal-600 font-bold" : ""
                }`}
                onClick={() => handleNavigation("/product")}
              >
                Products
              </li>
              <li className="p-4 hover:bg-gray-100">
                <Link
                  to="/about"
                  className={`${isActive("/about") ? "text-teal-600 font-bold" : ""}`}
                >
                  About Us
                </Link>
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

      {/* Spacer to prevent content from being hidden behind fixed navbar */}
      <div className="pt-[80px]"></div>
    </>
  );
};

export default Navbar;
