import React from "react";
import { useUserAuth } from "./components/useUserAuth";
import { useNavigate } from "react-router-dom";
import { userLoginAPI } from "../api/axios";
import UserLoginForm from "./components/UserLoginForm";

function LoginPage({ handleGotoRegisterPage }) {
  const userAuth = useUserAuth();
  const navigate = useNavigate();

  const handleSubmit = (values, { setSubmitting }) => {
    setSubmitting(true);
    userLoginAPI({ userLoginObject: values })
      .then((data) => {
        console.log("Login successful");
        userAuth.setUserID(data)
        userAuth.toggleUserLogin();
        navigate("/");
      })
      .catch((err) => {
        console.log("Login failed: ", err);
        setSubmitting(false);
      });
  };

  return (
    <div className="w-full h-full">
      <div className="flex flex-row justify-center">
        <div className="flex flex-col w-[600px] h-[500px] loginDiv justify-center">
          <span className="ml-[50px] font-mono text-5xl font-bold text-[#008080]">
            Login
          </span>
          <span className="ml-[50px] mt-[40px] font-mono text-gray-500">
            Welcome Back!!
          </span>
          <div className="flex justify-center mt-[10px]">
            <UserLoginForm handleSubmit={handleSubmit} />
          </div>
          <span className="font-mono flex flex-row justify-center mt-2">
            I don’t have an account?
            <span
              className="font-mono font-bold text-[#008080] cursor-pointer"
              onClick={handleGotoRegisterPage}
            >
              Sign-Up
            </span>
          </span>
        </div>
      </div>
    </div>
  );
}

export default LoginPage;
