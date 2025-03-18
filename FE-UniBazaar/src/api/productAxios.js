import axios from "axios";
const PRODUCT_BASE_URL = import.meta.env.VITE_PRODUCT_BASE_URL;

export const getAllProductsAPI = (limit, lastId) => {
    const params = {
      lastId: lastId,
      limit: limit,
    };
    return axios
      .get(PRODUCT_BASE_URL + "/products", { params })
      .then((response) => {
        console.error("products:", response);
        return response.data;
      })
      .catch((error) => {
        console.error("Error fetching products:", error);
        throw error;
      });
  };
  
  export const postProductAPI = (formData) => {
    return axios
      .post(PRODUCT_BASE_URL + "/products", formData)
      .then((response) => {
        alert("Product posted successfully!");
        return response.data;
      })
      .catch((error) => {
        console.error("Error posting product:", error);
        alert("Failed to post product. Try again.");
        throw error;
      });
  };
  