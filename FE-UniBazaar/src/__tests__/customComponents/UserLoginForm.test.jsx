import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { vi } from "vitest";
import UserLoginForm from "../../customComponents/UserLoginForm";

describe("UserLoginForm Component", () => {
  test("should allow login if credentials are valid", async () => {
    const mockSubmit = vi.fn();
    render(<UserLoginForm handleSubmit={mockSubmit} />);

    // Fill email and password fields
    await userEvent.type(screen.getByLabelText(/Email/i), "test@ufl.edu");
    await userEvent.type(screen.getByLabelText(/Password/i), "Password@123");

    // Click on Login button
    await userEvent.click(screen.getByRole("button", { name: /Login/i }));

    // Check if form is submitted
    expect(mockSubmit).toHaveBeenCalled();
  });

  test("should NOT allow login if fields are empty", async () => {
    const mockSubmit = vi.fn();
    render(<UserLoginForm handleSubmit={mockSubmit} />);

    // Click Login button without filling fields
    await userEvent.click(screen.getByRole("button", { name: /Login/i }));

    // Check that submit function is not called
    expect(mockSubmit).not.toHaveBeenCalled();
  });
});
