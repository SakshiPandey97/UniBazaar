import { userLoginAPI, userRegisterAPI } from "@/api/axios";
import { useUserAuth } from "@/hooks/useUserAuth";
import { useState } from "react";

export function useAuthHandler({ toggleModal }) {
  const useAuth = useUserAuth();
  const [successMessage, setSuccessMessage] = useState("");
  const [registeredEmail, setRegisteredEmail] = useState("");
  const [isRegistering, setIsRegistering] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [showConfetti, setShowConfetti] = useState(false);
  const [isVerifyingOTP, setIsVerifyingOTP] = useState(false);

  const handleRegisterSuccess = (email) => {
    setRegisteredEmail(email);
    setIsVerifyingOTP(true);
  };

  const toggleAuthMode = () => {
    setIsRegistering((prev) => !prev);
  };

  const handleSubmit = (values, { setSubmitting }) => {
    setSubmitting(true);
    setIsSubmitting(true);

    const apiCall = isRegistering ? userRegisterAPI : userLoginAPI;
    const requestObject = isRegistering
      ? { userRegisterObject: values }
      : { userLoginObject: values };

    apiCall(requestObject)
      .then((data) => {
        console.log(
          `${isRegistering ? "Registration" : "Login"} successful`,
          data
        );

        if (isRegistering) {
          setSuccessMessage(
            "Registration successful! Redirecting to Email Verification..."
          );
          setShowConfetti(true);
          handleRegisterSuccess(values.email);
          setTimeout(() => {
            setShowConfetti(false);
          }, 4000);
        } else {
          setSuccessMessage("Login successful! Redirecting to Home...");
          setTimeout(() => {
            setSuccessMessage("");
            useAuth.toggleUserLogin();
            toggleModal();
          }, 3000);
        }
      })
      .catch((err) => {
        console.error(
          `${isRegistering ? "Registration" : "Login"} failed:`,
          err
        );
        setSuccessMessage(err);
        setTimeout(() => setSuccessMessage(""), 3000);
        setSubmitting(false);
        setIsSubmitting(false);
      });
  };

  return {
    isRegistering,
    showConfetti,
    isSubmitting,
    successMessage,
    isVerifyingOTP,
    toggleAuthMode,
    handleSubmit,
  };
}
