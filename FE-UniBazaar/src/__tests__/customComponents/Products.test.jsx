import { render, screen } from "@testing-library/react";
import { vi } from "vitest";
import Products from "../../customComponents/Products";
import useProducts from "../../hooks/useProducts";

// Mock the useProducts hook
vi.mock("../../hooks/useProducts", () => ({
  default: vi.fn(),
}));

describe("Products Component", () => {
  it("displays loading spinner while fetching products", () => {
    useProducts.mockReturnValue({
      products: [],
      loading: true,
      error: null,
    });

    render(<Products />);

    // Assert the spinner is displayed when loading
    expect(screen.getByTestId("spinner-container")).toBeInTheDocument();
  });

  it("displays error message if there is an error", () => {
    useProducts.mockReturnValue({
      products: [],
      loading: false,
      error: "Error fetching products",
    });

    render(<Products />);

    // Assert the error message is displayed
    expect(screen.getByText("Error fetching products")).toBeInTheDocument();
  });

});
