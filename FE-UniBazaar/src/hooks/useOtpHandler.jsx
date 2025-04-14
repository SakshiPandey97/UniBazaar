import { useEffect, useState } from "react";
import { userVerificationAPI } from "@/api/userAxios";

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
      if (res.success) {
        setMessage("OTP Verified Successfully!");
      } else {
        setMessage("Invalid OTP. Please try again.");
      }
    } catch (err) {
      setMessage("Something went wrong.");
    } finally {
      setIsSubmitting(false);
    }
  };

  // Timer countdown
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
  };
}
