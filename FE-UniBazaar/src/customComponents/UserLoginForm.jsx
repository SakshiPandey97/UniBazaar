import React, { useState } from "react";
import { Formik, Form } from "formik";
import InputField from "./InputField";
import { validationSchema } from "@/utils/validSchema";

const initialValues = { email: "", password: "" };

function UserLoginForm({ handleLoginFormSubmission }) {
  const { handleSubmit, isSubmitting } = handleLoginFormSubmission;

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={handleSubmit}
      role="form"
    >
      {({ isSubmitting, handleChange }) => (
        <Form className="w-full">
          <InputField
            data_testid="loginEmail"
            name="email"
            type="email"
            placeholder="Email"
            isSubmitting={isSubmitting}
            onChange={handleChange}
          />

          <InputField
            data_testid="loginPassowrd"
            name="password"
            type="password"
            placeholder="Password"
            isSubmitting={isSubmitting}
            isPassword={true}
          />

          <div className="flex justify-center">
            <button
              data-testid="submitLoginBtn"
              type="submit"
              disabled={isSubmitting}
              className="w-1/3 hover:border-[#F58B00] border-2 p-2 bg-[#F58B00] hover:bg-[#FFC67D] text-balck font-bold py-2 px-4 rounded-md transition disabled:bg-gray-400"
            >
              {isSubmitting ? "Submitting..." : "Login"}
            </button>
          </div>
        </Form>
      )}
    </Formik>
  );
}

export default UserLoginForm;
