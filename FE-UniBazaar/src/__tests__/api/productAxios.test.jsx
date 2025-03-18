import { describe, test, expect, vi } from "vitest";
import axios from "axios";
import { getAllProductsAPI, postProductAPI } from "../../api/productAxios"; // Adjust import path if needed

// Mock axios
vi.mock("axios");

describe("Product API Tests", () => {
  const mockProducts = [
    { id: "1", name: "Product 1", price: 10 },
    { id: "2", name: "Product 2", price: 20 },
  ];

  afterEach(() => {
    vi.clearAllMocks();
  });

  test("getAllProductsAPI should return product data", async () => {
    axios.get.mockResolvedValue({ data: mockProducts });

    const result = await getAllProductsAPI(10, "");

    expect(result).toEqual(mockProducts);
    expect(axios.get).toHaveBeenCalledWith(expect.stringContaining("/products"), {
      params: { lastId: "", limit: 10 },
    });
  });

  test("getAllProductsAPI should throw an error on failure", async () => {
    axios.get.mockRejectedValue(new Error("Fetch failed"));

    await expect(getAllProductsAPI(10, "")).rejects.toThrow("Fetch failed");
  });

  test("postProductAPI should return success message", async () => {
    const mockResponse = { message: "Product created successfully" };
    axios.post.mockResolvedValue({ data: mockResponse });

    const formData = { name: "New Product", price: 30 };
    const result = await postProductAPI(formData);

    expect(result).toEqual(mockResponse);
    expect(axios.post).toHaveBeenCalledWith(expect.stringContaining("/products"), formData);
  });

  test("postProductAPI should throw an error on failure", async () => {
    axios.post.mockRejectedValue(new Error("Post failed"));

    await expect(postProductAPI({ name: "New Product", price: 30 })).rejects.toThrow("Post failed");
  });
});
