import React, { useState } from "react";
import { Formik, Form } from "formik";
import InputField from "./InputField";
import { validationSchema } from "@/utils/validSchema";
import { requestPasswordResetOTP, resetUserPassword } from "@/api/userAxios";

function UserForgetPasswordForm() {
  const [showOTPForm, setShowOTPForm] = useState(false);
  const [email, setEmail] = useState("");
  
  const handleEmailFormSubmit = async (values, actions) => {
    const success = await requestPasswordResetOTP(values.email);
    if (success) {
      setEmail(values.email);
      setShowOTPForm(true);
    } else {
      actions.setSubmitting(false);
    }
  };

  const handleResetFormSubmit = async (values, actions) => {
    const payload = {
      email,
      otp: values.otp,
      password: values.password,
    };

    try {
      const response = await resetUserPassword(payload);
      if (response.success) {
        setShowOTPForm(false);
        setEmail("");
      }
    } catch (error) {
      console.error("Password reset failed:", error);
    } finally {
      actions.setSubmitting(false);
    }
  };

  return (
    <>
      {!showOTPForm ? (
        <Formik
          key="email-form"
          initialValues={{ email: "" }}
          validationSchema={validationSchema.pick(["email"])}
          onSubmit={handleEmailFormSubmit}
        >
          {({ isSubmitting }) => (
            <Form className="w-full">
              <InputField
                data_testid="email-input"
                name="email"
                type="email"
                placeholder="Email"
                isSubmitting={isSubmitting}
              />
              <div className="flex justify-center">
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="w-1/3 hover:border-[#F58B00] border-2 p-2 bg-[#F58B00] hover:bg-[#FFC67D] text-black font-bold py-2 px-4 rounded-md transition disabled:bg-gray-400"
                >
                  {isSubmitting ? "Sending..." : "Get OTP"}
                </button>
              </div>
            </Form>
          )}
        </Formik>
      ) : (
        <Formik
          key="otp-form"
          initialValues={{ otp: "", password: "" }}
          validationSchema={validationSchema.pick(["otp", "password"])}
          onSubmit={handleResetFormSubmit}
        >
          {({ isSubmitting }) => (
            <Form className="w-full">
              <InputField
                data_testid="otp-input"
                name="otp"
                type="text"
                placeholder="Enter OTP"
                isSubmitting={isSubmitting}
              />
              <InputField
                data_testid="password-input"
                name="password"
                type="password"
                placeholder="New Password"
                isSubmitting={isSubmitting}
              />
              <div className="flex justify-center">
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="w-1/3 hover:border-[#F58B00] border-2 p-2 bg-[#F58B00] hover:bg-[#FFC67D] text-black font-bold py-2 px-4 rounded-md transition disabled:bg-gray-400"
                >
                  {isSubmitting ? "Resetting..." : "Reset Password"}
                </button>
              </div>
            </Form>
          )}
        </Formik>
      )}
    </>
  );
}

export default UserForgetPasswordForm;
