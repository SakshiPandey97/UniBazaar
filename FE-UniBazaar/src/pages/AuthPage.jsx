import React from "react";
import { motion, AnimatePresence } from "framer-motion";
import Confetti from "react-confetti";
import { useWindowSize } from "react-use";
import { useAuthHandler } from "@/hooks/useAuthHandler";
import InputOtp from "@/customComponents/InputOtp";
import UserLoginView from "@/customComponents/UserLoginView";
import UserRegisterView from "@/customComponents/UserRegisterView";

function AuthPage({ toggleModal }) {
  const { width, height } = useWindowSize();

  const {
    isRegistering,
    showConfetti,
    isSubmitting,
    successMessage,
    isVerifyingOTP,
    registeredEmail,
    toggleAuthMode,
    handleSubmit,
  } = useAuthHandler({ toggleModal });

  return (
    <>
      {showConfetti && <Confetti width={width} height={height} />}

      <div className="w-full h-full flex justify-center items-center relative">
          {isVerifyingOTP ? (
              <InputOtp email={registeredEmail} />
          ) : isRegistering ? (
              <UserRegisterView
                authHandlers={{
                  isSubmitting,
                  successMessage,
                  toggleAuthMode,
                  handleSubmit,
                }}
              />
          ) : (
              <UserLoginView
                authHandlers={{
                  isSubmitting,
                  successMessage,
                  toggleAuthMode,
                  handleSubmit,
                }}
              />
          )}
      </div>
    </>
  );
}

export default AuthPage;
