import { renderHook, act } from "@testing-library/react";
import { vi } from "vitest";
import useFetchProducts from "../../hooks/useFetchProducts"; // Adjust the path based on your project
import { getAllProductsAPI } from "../../api/productAxios";
import Product from "../../modal/product";

// ✅ Mock the API module
vi.mock("../../api/productAxios", () => ({
  getAllProductsAPI: vi.fn(),
}));

describe("useFetchProducts Hook", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test("fetches and sets products successfully", async () => {
    const mockData = [
      { productId: "1", productTitle: "Product A", productPostDate: "2024-03-10" },
      { productId: "2", productTitle: "Product B", productPostDate: "2024-03-15" },
    ];

    getAllProductsAPI.mockResolvedValue(mockData); // ✅ Mock API response

    const { result } = renderHook(() => useFetchProducts());

    // ✅ Before API call completes
    expect(result.current.loading).toBe(true);
    expect(result.current.products).toEqual([]);

    // ✅ Wait for hook to finish updating state
    await act(async () => {});

    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();

    // ✅ Verify sorted product order
    expect(result.current.products).toHaveLength(2);
    expect(result.current.products[0].productTitle).toBe("Product B"); // Most recent first
    expect(result.current.products[1].productTitle).toBe("Product A");
  });

  test("handles API error correctly", async () => {
    getAllProductsAPI.mockRejectedValue(new Error("API Error")); // ✅ Simulate API failure

    const { result } = renderHook(() => useFetchProducts());

    // ✅ Before API call completes
    expect(result.current.loading).toBe(true);
    expect(result.current.error).toBeNull();

    // ✅ Wait for hook to finish updating state
    await act(async () => {});

    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBe("Error fetching products");
    expect(result.current.products).toEqual([]); // No products due to error
  });
});
