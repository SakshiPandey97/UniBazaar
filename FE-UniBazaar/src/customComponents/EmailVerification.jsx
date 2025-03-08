import React from "react";
import { InputOtp } from "./InputOtp";
function EmailVerification() {
  return (
    <div
      className="flex flex-col w-[500px] h-[400px]
            loginDiv justify-center p-6 bg-white rounded-lg "
    >
      <div className="flex justify-center font-bold text-5xl m-2 ">
        Email Verifcation
      </div>
      <div className="flex justify-center text-2xl m-2">
        Please Enter the OTP
      </div>
      <InputOtp className="flex justify-center " />
    </div>
  );
}

export default EmailVerification;
