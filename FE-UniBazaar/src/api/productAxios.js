import axios from "axios";
import { toast } from "react-toastify";

const PRODUCT_BASE_URL = import.meta.env.VITE_PRODUCT_BASE_URL;

export const getAllProductsAPI = async (limit, lastId) => {
  const params = {
    lastId: lastId,
    limit: limit,
  };
  try {
    const response = await axios
      .get(PRODUCT_BASE_URL + "/products", { params });
    return response.data;
  } catch (error) {
    console.error("Error fetching products:", error);
    toast.error("Failed to fetch products.");
    throw error;
  }
};

export const postProductAPI = async (formData) => {
  try {
    const response = await axios
      .post(PRODUCT_BASE_URL + "/products", formData);
    toast.success("Product posted successfully!");
    return response.data;
  } catch (error) {
    console.error("Error posting product:", error);
    toast.error("Failed to post product. Try again.");
    throw error;
  }
};

export const updateProductAPI = async (userId, productId, formData) => {
  try {
    const response = await axios.put(`${PRODUCT_BASE_URL}/products/${userId}/${productId}`, formData);
    toast.success("Product updated successfully!");
    return response.data;
  } catch (error) {
    console.error("Error updating product:", error);
    toast.error("Failed to update product. Try again.");
    throw error;
  }
};

export const getUserProductsAPI = async (userId, limit, lastId) => {
  const params = {
    lastId: lastId,
    limit: limit,
  };
  try {
    const response = await axios.get(`${PRODUCT_BASE_URL}/products/${userId}`, { params });
    return response.data;
  } catch (error) {
    if (error.response && error.response.status === 404) {
      console.warn("No products found for user:", userId);
      return [];
    }
    toast.error("Failed to fetch your products.");
    throw error;
  }
};


export const searchProductsAPI = async (query, limit) => {
  try {
    const response = await axios.get(`${PRODUCT_BASE_URL}/search/products`, {
      params: { query, limit },
    });
    return response.data;
  } catch (error) {
    if (error.response && error.response.status === 404) {
      console.log("No products found (404).");
      return [];
    }
    toast.error("Failed to search products.");
    throw error;
  }
};

export const deleteProductAPI = async (userId, productId) => {
  try {
    const response = await axios.delete(`${PRODUCT_BASE_URL}/products/${userId}/${productId}`);
    toast.success("Product deleted successfully!");
    return response.data;
  } catch (error) {
    console.error("Error deleting product:", error);
    toast.error("Failed to delete product. Try again.");
    throw error;
  }
};