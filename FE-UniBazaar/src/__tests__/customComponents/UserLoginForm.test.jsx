import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { vi } from "vitest";
import UserLoginForm from "../../customComponents/UserLoginForm";

describe("UserLoginForm Component", () => {
  test("should allow login if credentials are valid", async () => {
    const mockSubmit = vi.fn();
    render(<UserLoginForm handleSubmit={mockSubmit} />);

    await userEvent.type(screen.getByLabelText(/Email/i), "test@ufl.edu");
    await userEvent.type(screen.getByLabelText(/Password/i), "Password@123");

    await userEvent.click(screen.getByRole("button", { name: /Login/i }));

    expect(mockSubmit).toHaveBeenCalled();
  });

  test("should NOT allow login if fields are empty", async () => {
    const mockSubmit = vi.fn();
    render(<UserLoginForm handleSubmit={mockSubmit} />);

    await userEvent.click(screen.getByRole("button", { name: /Login/i }));

    expect(mockSubmit).not.toHaveBeenCalled();
  });
});
