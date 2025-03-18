import { renderHook, act } from "@testing-library/react";
import { vi } from "vitest";
import useProducts from "../../hooks/useProducts";
import { getAllProductsAPI } from "../../api/userAxios";

vi.mock("../../api/axios", () => ({
  getAllProductsAPI: vi.fn(() => Promise.resolve([
    { productId: 1, productTitle: "Product 1" },
    { productId: 2, productTitle: "Product 2" },
    { productId: 3, productTitle: "Product 3" },
  ])),
  postProductAPI: vi.fn(() => Promise.resolve({ productId: 3, productTitle: "Product 3" })),
}));

describe("useProducts Hook", () => {
  it("fetches all products", async () => {
    const { result } = renderHook(() => useProducts());

    await act(async () => {});

    getAllProductsAPI.mockResolvedValueOnce([
      { productId: 1, productTitle: "Product 1" },
      { productId: 2, productTitle: "Product 2" },
      { productId: 3, productTitle: "Product 3" }, 
    ]);

    expect(result.current.products).toHaveLength(3); 
  });
});
