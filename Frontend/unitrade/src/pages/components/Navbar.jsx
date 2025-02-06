import React, { useState, useEffect } from "react";
import logo from "../../assets/imgs/logo.svg";
import login from "../../assets/imgs/login.svg";
import menuToggle from "../../assets/imgs/menu-toggle.svg";
import close from "../../assets/imgs/close.svg";
import { Link, useNavigate } from "react-router-dom";
import { useUserAuth } from "./useUserAuth";

const Navbar = () => {
  const [toggle, setToggle] = useState(false);
  const handleClick = () => setToggle(!toggle);
  const navigate = useNavigate();
  const userAuth = useUserAuth();

  const checkForUserLogIn = (path) => {
    if (userAuth.userState) {
      navigate(path);
    } else {
      navigate("/auth/login");
    }
  };

  useEffect(() => {}, [userAuth.userState]);
  return (
    <div className="w-full h-[80px] bg-[#F8F8F8] border-b py-2">
      <div className="md:max-w-[1480px] max-w-[600px] m-auto w-full flex justify-between items-center">
        <Link to="/">
          <div className="relative inline-block">
            <img
              src={logo}
              className="h-[50px] md:h-[60px] w-auto transition duration-300 hover:opacity-75"
            />
            <div className="absolute inset-0 bg-blue-500 opacity-0 hover:opacity-30 transition duration-300"></div>
          </div>
        </Link>
        <div className="hidden md:flex items-center">
          <ul className="flex gap-20 font-[serif] font-[playfair display] text-[32px] font-[medium]">
            <li>
              <Link to="/">Home</Link>
            </li>
            <li onClick={() => checkForUserLogIn("/buy")}>Buy</li>
            <li onClick={() => checkForUserLogIn("/sell")}>Sell</li>
            <li onClick={() => checkForUserLogIn("/product")}>Products</li>
            <li>
              <Link to="/about">About Us</Link>
            </li>
          </ul>
        </div>
        <div className="hidden md:flex">
          <button
            className="flex justify-between items-center bg-transparent px-6 gap-2 font-[playfair display] text-[24px] font-[medium]"
            onClick={() => navigate("/auth/")}
          >
            <img src={login} className="h-[24px]" />
            {!userAuth.userState ? "Login" : "Log-Out"}
          </button>
          <button className="px-4 py-3 rounded-md bg-[#008080] text-white font-[playfair display] text-[24px] font-[medium]">
            Transactions
          </button>
        </div>
        <div className="md:hidden" onClick={handleClick}>
          <img src={toggle ? close : menuToggle} />
        </div>
      </div>

      <div
        className={
          toggle ? "absolute z-10 p-4 bg-white w-full px-8 md:hidden" : "hidden"
        }
      >
        <ul>
          <li className="p-4 hover:bg-gray-100">Buy</li>
          <li className="p-4 hover:bg-gray-100">Sell</li>
          <li className="p-4 hover:bg-gray-100">Products</li>
          <li className="p-4 hover:bg-gray-100">About Us</li>
          <div className="flex flex-col my-4 gap-4">
            <button className="border border-[#008080] flex justify-center items-center bg-transparent px-6 gap-2 py-3">
              <img
                src={login}
                className="h-[20px] font-[serif] font-[playfair display] text-[24px] font-[medium]"
              />
              <Link to="/login">Login</Link>
            </button>
            <button className="px-4 py-5 rounded-md bg-[#008080] text-white">
              Transactions
            </button>
          </div>
        </ul>
      </div>
    </div>
  );
};

export default Navbar;
