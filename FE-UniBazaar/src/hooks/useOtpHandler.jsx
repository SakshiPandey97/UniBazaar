import { useEffect, useState } from "react";
import { userVerificationAPI, userResendAPI } from "@/api/userAxios";

export function useOtpHandler(email) {
  const [otp, setOtp] = useState("");
  const [message, setMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [timeLeft, setTimeLeft] = useState(60);

  const handleChange = (value) => {
    setOtp(value);
  };

  const handleSubmit = async () => {
    if (otp.length !== 6) {
      setMessage("Please enter a valid 6-digit OTP.");
      return;
    }

    setIsSubmitting(true);
    setMessage("");

    try {
      const res = await userVerificationAPI({ email, code: otp });
      console.log(res)
      if (res.verified) {
        setMessage("OTP Verified Successfully!");
        setTimeout(() => {
          window.location.href = "/"; 
        }, 1000);
      } else {
        setMessage("Invalid OTP. Please try again.");
      }
    } catch (err) {
      console.error("OTP verification failed:", err);
      setMessage("Something went wrong.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleResend = async () => {
    setIsSubmitting(true);
    setMessage("");
    try {
      await userResendAPI({ email });
      setMessage("OTP resent successfully!");
      setTimeLeft(60);
      setOtp("");
    } catch (error) {
      console.error("Resend OTP failed:", error);
      setMessage("Failed to resend OTP. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  useEffect(() => {
    if (timeLeft <= 0) return;
    const timer = setInterval(() => {
      setTimeLeft((prev) => prev - 1);
    }, 1000);
    return () => clearInterval(timer);
  }, [timeLeft]);

  return {
    otp,
    message,
    isSubmitting,
    handleChange,
    handleSubmit,
    timeLeft,
    resetTimer: handleResend,
  };
}
