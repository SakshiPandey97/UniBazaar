import React from "react";
import UserRegisterForm from "./components/UserRegisterForm";
import { userRegisterAPI } from "../api/axios";

function RegisterPage() {
  const handleSubmit = (values, { setSubmitting }) => {
    console.log("Form values:", values);
    setSubmitting(true);
    userRegisterAPI({ userRegisterObject: values })
      .then((data) => {
        console.log("Registration successful", data);
      })
      .catch((err) => {
        console.log("Registration failed: ", err);
        setSubmitting(false);
      });
  };

  return (
    <div className="w-full h-full">
      <div className="flex flex-row justify-center">
        <div className="flex flex-col w-[600px] h-[500px] loginDiv justify-center">
          <span className="ml-[50px] font-mono text-5xl font-bold text-[#008080]">
            Registration
          </span>
          <span className="ml-[50px] mt-[40px] font-mono text-gray-500">
            Let's Buy and Sell!!
          </span>
          <div className="flex justify-center mt-[10px]">
            <UserRegisterForm handleSubmit={handleSubmit} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default RegisterPage;
