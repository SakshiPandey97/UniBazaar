import { render, screen } from "@testing-library/react";
import Spinner from "../../customComponents/Spinner"; // Adjust the import according to your file structure

describe("Spinner Component", () => {
  test("renders the spinner element", () => {
    render(<Spinner />);
    
    // Check if the spinner is in the document
    const spinner = screen.getByRole("status"); // We use the "status" role for loading spinners by default
    expect(spinner).toBeInTheDocument();  // This ensures the spinner is rendered
  });

  test("renders with correct className", () => {
    render(<Spinner />);
    
    // Check if the spinner container has the correct class
    const spinnerContainer = screen.getByTestId("spinner-container");
    expect(spinnerContainer).toHaveClass("spinner-container"); // Ensure class name matches
  });
});
