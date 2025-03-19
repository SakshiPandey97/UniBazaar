import { render, screen } from "@testing-library/react";
import EmailVerification from "../../customComponents/EmailVerification";
import { vi } from "vitest";

// Mock the InputOtp component
vi.mock("../../customComponents/InputOtp", () => ({
  InputOtp: () => <div data-testid="input-otp">Mock OTP Input</div>,
}));

describe("EmailVerification Component", () => {
  test("renders without crashing", () => {
    render(<EmailVerification />);

    expect(
      screen.getByText("Email Verifcation", { exact: false })
    ).toBeInTheDocument();

    expect(
      screen.getByText("Please Enter the OTP", { exact: false })
    ).toBeInTheDocument();
  });

  test("renders the InputOtp component", () => {
    render(<EmailVerification />);
    expect(screen.getByTestId("input-otp")).toBeInTheDocument();
  });
});
