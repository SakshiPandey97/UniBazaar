import React, { useState } from "react";
import UserLoginForm from "./UserLoginForm";
import UserForgetPasswordForm from "./UserForgetPasswordForm";

function UserLoginView({ authHandlers }) {
  const { isSubmitting, successMessage, toggleAuthMode, handleSubmit } =
    authHandlers;

  const [isForgotPassword, setIsForgotPassword] = useState(false);

  return (
    <div className="flex flex-col w-full sm:w-[400px] md:w-[500px] lg:w-[550px] xl:w-[600px] h-[500px] loginDiv justify-center p-6 bg-white rounded-lg shadow-lg">
      <h1 className="ml-12 font-mono text-5xl font-bold text-[#032B54]">
        {isForgotPassword ? "Reset Password" : "Login"}
      </h1>
      <p className="ml-12 mt-10 font-mono text-gray-500">
        {isForgotPassword ? "Enter your email to reset password" : "Welcome Back!!"}
      </p>

      {successMessage && (
        <p className="text-[#F58B00] font-mono text-center mt-2 animate-fadeIn">
          {successMessage}
        </p>
      )}

      <div className="flex justify-center mt-2">
        {isForgotPassword ? (
          <UserForgetPasswordForm
            handleResetFormSubmission={{ handleSubmit, isSubmitting }}
          />
        ) : (
          <UserLoginForm
            handleLoginFormSubmission={{ handleSubmit, isSubmitting }}
          />
        )}
      </div>

      <div className="mt-4">
        {!isForgotPassword && (
          <p className="font-mono flex justify-center">
            I don’t have an account?
            <span
              data-testid="toggleLoginRegister"
              className="font-bold text-[#032B54] cursor-pointer ml-1 hover:underline"
              onClick={toggleAuthMode}
            >
              Sign-Up
            </span>
          </p>
        )}

        <span
          className="text-sm font-medium text-gray-500 cursor-pointer flex justify-center mt-2"
          onClick={() => setIsForgotPassword(!isForgotPassword)}
        >
          {isForgotPassword ? "Back to Login" : "Forgot password?"}
        </span>
      </div>
    </div>
  );
}

export default UserLoginView;
