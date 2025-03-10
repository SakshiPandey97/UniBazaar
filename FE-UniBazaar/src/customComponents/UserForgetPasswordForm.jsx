import React, { useState } from "react";
import { Formik, Form } from "formik";
import PasswordReqBox from "./PasswordReqBox";
import InputField from "./InputField";
import { validationSchema } from "@/utils/validSchema";

const initialValues = { email: "", password: "" };

function UserForgetPasswordForm({ handleResetFormSubmission }) {
  const {handleSubmit, isSubmitting} = handleResetFormSubmission;
  const [password, setPassword] = useState("");
  const [isFocused, setIsFocused] = useState(false);

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={handleSubmit}
    >
      {({ isSubmitting, handleChange }) => (
        <Form className="w-full">
          <InputField
            label="Email"
            name="email"
            type="email"
            isSubmitting={isSubmitting}
            onChange={handleChange}
          />
          <div className="flex flex-col">
          <InputField
            label="New Password"
            name="password"
            type="password"
            isSubmitting={isSubmitting}
            onFocus={() => setIsFocused(true)}
            onBlur={() => setIsFocused(false)}
            onChange={(e) => {
              setPassword(e.target.value);
              handleChange(e);
            }}
          />
          </div>
          {isFocused && <PasswordReqBox password={password} />}

          <div className="flex justify-center">
            <button
              type="submit"
              disabled={isSubmitting}
              className="w-1/3 hover:border-[#F58B00] border-2 p-2 bg-[#F58B00] hover:bg-[#FFC67D] text-balck font-bold py-2 px-4 rounded-md transition disabled:bg-gray-400"
            >
              {isSubmitting ? "Submitting..." : "Reset"}
            </button>
          </div>
        </Form>
      )}
    </Formik>
  );
}

export default UserForgetPasswordForm;
