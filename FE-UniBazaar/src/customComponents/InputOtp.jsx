import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from "@/components/ui/input-otp";

export function InputOtp() {
  const REGEXP_ONLY_DIGITS = "^\\d+$";
  return (
    <InputOTP maxLength={6} pattern={REGEXP_ONLY_DIGITS}>
      <InputOTPGroup>
        <InputOTPSlot
          index={0}
          className="w-14 h-14 text-2xl focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <InputOTPSlot
          index={1}
          className="w-14 h-14 text-2xl focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <InputOTPSlot
          index={2}
          className="w-14 h-14 text-2xl focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </InputOTPGroup>
      <InputOTPSeparator />
      <InputOTPGroup>
        <InputOTPSlot
          index={3}
          className="w-14 h-14 text-2xl  focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <InputOTPSlot
          index={4}
          className="w-14 h-14 text-2xl  focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <InputOTPSlot
          index={5}
          className="w-14 h-14 text-2xl  focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </InputOTPGroup>
    </InputOTP>
  );
}
