import { renderHook, act } from "@testing-library/react";
import { vi } from "vitest";
import useAllProducts from "../../hooks/useAllProducts"; // Adjust path as needed
import { getAllProductsAPI } from "../../api/productAxios";
import Product from "../../modal/product";

// ✅ Mock API module
vi.mock("../../api/productAxios", () => ({
  getAllProductsAPI: vi.fn(),
}));

describe("useAllProducts Hook", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test("fetches and sets paginated products successfully", async () => {
    const mockData = [
      { productId: "1", productTitle: "Product A" },
      { productId: "2", productTitle: "Product B" },
    ];

    getAllProductsAPI.mockResolvedValue(mockData); // ✅ Mock API response

    const { result } = renderHook(() => useAllProducts(2, ""));

    // ✅ Before API call completes
    expect(result.current.loading).toBe(true);
    expect(result.current.products).toEqual([]);
    expect(result.current.hasMoreProducts).toBe(true);

    // ✅ Wait for hook to finish updating state
    await act(async () => {});

    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
    expect(result.current.products).toHaveLength(2);
    expect(result.current.products[0].productTitle).toBe("Product A");
    expect(result.current.products[1].productTitle).toBe("Product B");
    expect(result.current.hasMoreProducts).toBe(true);
  });

  test("handles API returning fewer products than limit", async () => {
    const mockData = [{ productId: "3", productTitle: "Product C" }];

    getAllProductsAPI.mockResolvedValue(mockData); // ✅ Mock API response

    const { result } = renderHook(() => useAllProducts(5, ""));

    await act(async () => {});

    expect(result.current.loading).toBe(false);
    expect(result.current.products).toHaveLength(1);
    expect(result.current.hasMoreProducts).toBe(false); // ✅ No more products to load
  });

  test("handles API error correctly", async () => {
    getAllProductsAPI.mockRejectedValue(new Error("API Error")); // ✅ Simulate API failure

    const { result } = renderHook(() => useAllProducts(2, ""));

    await act(async () => {});

    expect(result.current.loading).toBe(false);
    expect(result.current.error).toContain("Failed to fetch products"); // ✅ More flexible match
    expect(result.current.products).toEqual([]);
    expect(result.current.hasMoreProducts).toBe(true); // ✅ Keeps trying
  });

  test("handles 404 response and stops fetching", async () => {
    getAllProductsAPI.mockRejectedValue({ response: { status: 404 } }); // ✅ Simulate 404 error

    const { result } = renderHook(() => useAllProducts(2, ""));

    await act(async () => {});

    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
    expect(result.current.products).toEqual([]);
    expect(result.current.hasMoreProducts).toBe(false); // ✅ Stops fetching
  });
});
