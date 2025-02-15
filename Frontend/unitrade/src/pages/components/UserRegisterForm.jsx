import React from "react";
import { Formik, Form, Field, ErrorMessage } from "formik";
import * as Yup from "yup";

const validationSchema = Yup.object({
  name: Yup.string().required("Name is required"),
  email: Yup.string()
    .email("Invalid email address")
    .required("Email is required"),
  password: Yup.string()
    .min(6, "Password must be at least 6 characters")
    .required("Password is required"),
});

function UserRegisterForm({ handleSubmit }) {
  const initialValues = { name: "", email: "", password: "" };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={handleSubmit}
    >
      {({ isSubmitting }) => (
        <Form className="w-full">
          {/* Name Field */}
          <div className="flex flex-col mb-4">
            <label
              htmlFor="name"
              className="ml-[50px] block text-sm font-medium text-gray-700"
            >
              Name
            </label>
            <div className="flex flex-row justify-center">
              <Field
                type="text"
                id="name"
                name="name"
                className="mt-1 p-2 w-4/5 border rounded-md focus:outline-none focus:ring focus:ring-[#6D9886]"
              />
            </div>
            <ErrorMessage
              name="name"
              component="div"
              className="text-red-500 text-sm mt-1 ml-[50px]"
            />
          </div>

          {/* Email Field */}
          <div className="flex flex-col mb-4">
            <label
              htmlFor="email"
              className="ml-[50px] block text-sm font-medium text-gray-700"
            >
              Email
            </label>
            <div className="flex flex-row justify-center">
              <Field
                type="email"
                id="email"
                name="email"
                className="mt-1 p-2 w-4/5 border rounded-md focus:outline-none focus:ring focus:ring-[#6D9886]"
              />
            </div>
            <ErrorMessage
              name="email"
              component="div"
              className="text-red-500 text-sm mt-1 ml-[50px]"
            />
          </div>

          {/* Password Field */}
          <div className="flex flex-col mb-4">
            <label
              htmlFor="password"
              className="ml-[50px] block text-sm font-medium text-gray-700"
            >
              Password
            </label>
            <div className="flex flex-row justify-center">
              <Field
                type="password"
                id="password"
                name="password"
                className="mt-1 p-2 w-4/5 border rounded-md focus:outline-none focus:ring focus:ring-[#6D9886]"
              />
            </div>
            <ErrorMessage
              name="password"
              component="div"
              className="text-red-500 text-sm mt-1 ml-[50px]"
            />
          </div>

          {/* Submit Button */}
          <div className="flex flex-row justify-center">
            <button
              type="submit"
              disabled={isSubmitting}
              className="w-1/3 bg-[#6D9886] text-[#FFFFFF] text-white py-2 px-4 rounded-md hover:bg-[#008080] transition"
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
