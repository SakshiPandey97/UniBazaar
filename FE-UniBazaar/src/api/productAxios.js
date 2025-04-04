import axios from "axios";
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
    throw error;
  }
  };
  
  export const postProductAPI = async (formData) => {
    try {
      const response = await axios
        .post(PRODUCT_BASE_URL + "/products", formData);
      alert("Product posted successfully!");
      return response.data;
    } catch (error) {
      console.error("Error posting product:", error);
      alert("Failed to post product. Try again.");
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
        console.error("Error searching products:", error);
        throw error;
    }
};