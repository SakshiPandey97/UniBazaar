import { useState } from "react";
import { userLoginAPI, userRegisterAPI } from "../api/axios";
import { useUserAuth } from "./useUserAuth";

export function useAuthHandler({ toggleModal }) {
  const useAuth = useUserAuth();
  const [isRegistering, setIsRegistering] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

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
          setSuccessMessage("Registration successful! Redirecting to login...");
          setTimeout(() => {
            setIsRegistering(false);
            setSuccessMessage("");
            setIsSubmitting(false);
          }, 3000);
        } else {
          setSuccessMessage("Login successful! Redirecting to Home...");
          setTimeout(() => {
            setSuccessMessage("");
            useAuth.toggleUserLogin()
            toggleModal();
          }, 3000);
        }
      })
      .catch((err) => {
        console.error(
          `${isRegistering ? "Registration" : "Login"} failed:`,
          err
        );
        setSubmitting(false);
        setIsSubmitting(false);
      });
  };

  const toggleAuthMode = () => setIsRegistering((prev) => !prev);

  return {
    isRegistering,
    isSubmitting,
    successMessage,
    handleSubmit,
    toggleAuthMode,
  };
}
