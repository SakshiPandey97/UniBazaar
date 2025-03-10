import React from "react";
import UserRegisterForm from "./UserRegisterForm";

function UserRegisterView({ authHandlers }) {
  const { isSubmitting,successMessage, toggleAuthMode, handleSubmit } = authHandlers;

  return (
    <div
      className={
        "flex flex-col w-full sm:w-[400px] md:w-[500px] lg:w-[550px] xl:w-[600px] h-[600px] loginDiv justify-center p-6 bg-white rounded-lg shadow-lg"
      }
    >
      <h1 className="ml-12 font-mono text-5xl font-bold text-[#032B54]">
        Register
      </h1>
      <p className="ml-12 mt-10 font-mono text-gray-500">
        Lets start selling and buying together
      </p>

      {successMessage && (
        <p className="text-[#F58B00] font-mono text-center mt-2 animate-fadeIn">
          {successMessage}
        </p>
      )}

      <div className="flex justify-center mt-2">
        <UserRegisterForm
          handleRegisterFormSubmission={{ handleSubmit, isSubmitting }}
        />
      </div>

      <p className="font-mono flex justify-center mt-2">
        Already have a account ?
        <span
          data-testid="toggleLoginRegister"
          className="font-bold text-[#032B54] cursor-pointer ml-1 hover:underline"
          onClick={toggleAuthMode}
        >
          Login
        </span>
      </p>
    </div>
  );
}

export default UserRegisterView;
