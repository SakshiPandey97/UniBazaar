import { useState } from "react";
import { userVerificationAPI } from "@/api/userAxios";

export function useOtpHandler(email) {
  const [otp, setOtp] = useState("");
  const [message, setMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

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

    userVerificationAPI({ email:email, code:otp })
      .then((res) => {
        if (res.success) {
          setMessage("OTP Verified Successfully! ");
        } else {
          setMessage("Invalid OTP. Please try again ");
        }
      })
      .catch((err) => {
        setMessage("Error verifying OTP. Please try again.", err);
      })
      .finally(() => {
        setIsSubmitting(false);
      });
  };

  return { otp, message, isSubmitting, handleChange, handleSubmit };
}
