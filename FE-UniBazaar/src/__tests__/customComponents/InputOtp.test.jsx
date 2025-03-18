import { render, screen, fireEvent } from "@testing-library/react";
import InputOtp from "../../customComponents/InputOtp";
import { vi } from "vitest";
import { useOtpHandler } from "../../hooks/useOtpHandler";

// ✅ Mock ResizeObserver to prevent test failure
global.ResizeObserver = class {
  observe() {}
  unobserve() {}
  disconnect() {}
};

// ✅ Mock useOtpHandler before tests
vi.mock("../../hooks/useOtpHandler", () => ({
  useOtpHandler: vi.fn(() => ({
    otp: "",
    message: "",
    isSubmitting: false,
    handleChange: vi.fn(),
    handleSubmit: vi.fn(),
  })),
}));

describe("InputOtp Component", () => {
  test("renders correctly", () => {
    render(<InputOtp email="test@example.com" />);
    expect(screen.getByText("Enter OTP")).toBeInTheDocument();
  });

  test("calls handleChange when typing in OTP input", () => {
    const mockHandleChange = vi.fn();
    useOtpHandler.mockReturnValue({
      otp: "",
      message: "",
      isSubmitting: false,
      handleChange: mockHandleChange,
      handleSubmit: vi.fn(),
    });

    render(<InputOtp email="test@example.com" />);
    const otpInputs = screen.getAllByRole("textbox");

    fireEvent.change(otpInputs[0], { target: { value: "1" } });
    expect(mockHandleChange).toHaveBeenCalled();
  });

  test("calls handleSubmit when clicking the submit button", () => {
    const mockHandleSubmit = vi.fn();
    useOtpHandler.mockReturnValue({
      otp: "123456",
      message: "",
      isSubmitting: false,
      handleChange: vi.fn(),
      handleSubmit: mockHandleSubmit,
    });

    render(<InputOtp email="test@example.com" />);
    const submitButton = screen.getByRole("button", { name: "Submit OTP" });

    fireEvent.click(submitButton);
    expect(mockHandleSubmit).toHaveBeenCalledTimes(1);
  });

  test("disables the submit button when isSubmitting is true", () => {
    useOtpHandler.mockReturnValue({
      otp: "123456",
      message: "",
      isSubmitting: true,
      handleChange: vi.fn(),
      handleSubmit: vi.fn(),
    });

    render(<InputOtp email="test@example.com" />);
    const submitButton = screen.getByRole("button", { name: "Verifying..." });

    expect(submitButton).toBeDisabled();
  });
});
