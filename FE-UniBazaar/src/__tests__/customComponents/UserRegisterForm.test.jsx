// import { render, screen, fireEvent, waitFor } from "@testing-library/react";
// import UserRegisterForm from "../../customComponents/UserRegisterForm";

// describe("UserRegisterForm", () => {
//   it("should call handleSubmit with correct values when form is submitted", async () => {
//     const mockHandleSubmit = vi.fn();

//     render(<UserRegisterForm handleSubmit={mockHandleSubmit} />);

//     const nameInput = screen.getByLabelText("Name");
//     const emailInput = screen.getByLabelText("Email");
//     const passwordInput = screen.getByLabelText("Password");
//     const submitButton = screen.getByRole("button", { name: /Register/i });

//     fireEvent.change(nameInput, { target: { value: "Tanmay Testing" } });
//     fireEvent.change(emailInput, { target: { value: "testing@ufl.edu" } });
//     fireEvent.change(passwordInput, { target: { value: "Password@123" } }); // With uppercase

//     // Submit Form
//     fireEvent.click(submitButton);

//     // Wait for async operations
//     await waitFor(() => {
//       expect(mockHandleSubmit).toHaveBeenCalledWith(
//         { name: "Tanmay Testing", email: "testing@ufl.edu", password: "Password@123" },
//         expect.anything()
//       );
//     });
//   });
// });
