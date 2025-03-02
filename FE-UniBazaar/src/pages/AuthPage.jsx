import { useAuthHandler } from "../hooks/useAuthHandler";
import React from "react";
import Confetti from "react-confetti";
import { useWindowSize } from "react-use"; 
import UserRegisterForm from "../customComponents/UserRegisterForm";
import UserLoginForm from "../customComponents/UserLoginForm";

function AuthPage({ toggleModal }) {
  const { width, height } = useWindowSize(); 
  const {
    isRegistering,
    isSubmitting,
    successMessage,
    showConfetti,
    handleSubmit,
    toggleAuthMode,
  } = useAuthHandler({ toggleModal });

  return (
    <>
      {showConfetti && <Confetti width={width} height={height} />}
      <div className="w-full h-full flex justify-center items-center relative">
        <div className="flex flex-col w-[600px] h-[500px] loginDiv justify-center p-6 bg-white rounded-lg shadow-lg">
          <h1 className="ml-[50px] font-mono text-5xl font-bold text-[#032B54]">
            {isRegistering ? "Sign Up" : "Login"}
          </h1>
          <p className="ml-[50px] mt-[40px] font-mono text-gray-500">
            {isRegistering ? "Create your account" : "Welcome Back!!"}
          </p>

          {/* Success Message */}
          {successMessage && (
            <p className="text-[#F58B00] font-mono text-center mt-2 animate-fadeIn">
              {successMessage}
            </p>
          )}

          <div className="flex justify-center mt-[10px]">
            {isRegistering ? (
              <UserRegisterForm
                handleSubmit={handleSubmit}
                isSubmitting={isSubmitting}
              />
            ) : (
              <UserLoginForm
                handleSubmit={handleSubmit}
                isSubmitting={isSubmitting}
              />
            )}
          </div>

          <p className="font-mono flex flex-row justify-center mt-2">
            {isRegistering
              ? "Already have an account?"
              : "I don’t have an account?"}
            <span
              className="font-bold text-[#032B54] cursor-pointer ml-1 hover:underline"
              onClick={toggleAuthMode}
            >
              {isRegistering ? "Login" : "Sign-Up"}
            </span>
          </p>
        </div>
      </div>
    </>
  );
}

export default AuthPage;
