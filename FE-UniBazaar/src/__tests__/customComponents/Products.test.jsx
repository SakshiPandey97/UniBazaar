import { render, screen } from "@testing-library/react";
import { vi } from "vitest";
import Products from "../../customComponents/Products";
import useFetchProducts from "../../hooks/useFetchProducts";

// Mock the useFetchProducts hook
vi.mock("../../hooks/useFetchProducts", () => ({
  default: vi.fn(),
}));

describe("Products Component", () => {
  it("displays loading spinner while fetching products", () => {
    useFetchProducts.mockReturnValue({
      products: [],
      loading: true,
      error: null,
    });

    render(<Products />);

    // Assert the spinner is displayed when loading
    expect(screen.getByTestId("spinner-container")).toBeInTheDocument();
  });

  it("displays error message if there is an error", () => {
    useFetchProducts.mockReturnValue({
      products: [],
      loading: false,
      error: "Error fetching products",
    });

    render(<Products />);

    // Assert the error message is displayed
    expect(screen.getByText("Error fetching products")).toBeInTheDocument();
  });

});
