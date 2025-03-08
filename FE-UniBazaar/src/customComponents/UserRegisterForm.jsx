import React from "react";
import { Formik, Form } from "formik";
import { validationSchema } from "@/utils/validSchema";
import InputField from "./InputField";

function UserRegisterForm({ handleSubmit }) {
  const initialValues = { name: "", email: "", password: "", phoneNumber: "" };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={handleSubmit}
    >
      {({ isSubmitting, handleChange }) => (
        <Form className="w-full">
          <InputField
            data_testid="registerName"
            label="Name"
            name="name"
            type="text"
            isSubmitting={isSubmitting}
            onChange={handleChange}
          />

          <InputField
            data_testid="registerEmail"
            label="Email"
            name="email"
            type="email"
            isSubmitting={isSubmitting}
            onChange={handleChange}
          />

          <InputField
            data_testid="registerPassword"
            label="Password"
            name="password"
            type="password"
            isSubmitting={isSubmitting}
            onChange={handleChange}
          />
          <InputField
            data-testid="registerNumber"
            label="Phone Number"
            name="phoneNumber"
            type="tel"
            isSubmitting={isSubmitting}
            onChange={handleChange}
          />

          <div className="flex flex-row justify-center">
            <button
              data_testid="submitRegisterBtn"
              type="submit"
              disabled={isSubmitting}
              className="w-1/3 hover:border-[#F58B00] border-2 p-2 bg-[#F58B00] hover:bg-[#FFC67D] text-balck font-bold py-2 px-4 rounded-md transition disabled:bg-gray-400"
            >
              {isSubmitting ? "Registering..." : "Register"}
            </button>
          </div>
        </Form>
      )}
    </Formik>
  );
}

export default UserRegisterForm;
