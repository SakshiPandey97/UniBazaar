import axios from "axios";

const USER_BASE_URL = "http://127.0.0.1:4000";
const PRODUCT_BASE_URL = "https://unibazaar-products.azurewebsites.net";
const CHAT_USERS_BASE_URL = "http://127.0.0.1:8080";
export const userLoginAPI = ({ userLoginObject }) => {
  return axios
    .post(USER_BASE_URL + "/login", userLoginObject)
    .then((response) => {
      const userId = response.data.userId;
      localStorage.setItem("userId", userId);
      return userId;
    })
    .catch((error) => {
      console.error("Error logging in:", error);
      throw error;
    });
};

export const userRegisterAPI = ({ userRegisterObject }) => {
  console.log("Register Recevied Obj", userRegisterObject);
  return axios
    .post(USER_BASE_URL + "/signup", userRegisterObject)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      console.error("Error Registering: ", error);
      throw error;
    });
};

export const userVerificationAPI = ({ userVerificationObject }) => {
  return axios
    .post(USER_BASE_URL + "/verifyEmail", userVerificationObject)
    .then((response) => {
      return response.data.userId;
    })
    .catch((error) => {
      console.error("Error Verifying user in:", error);
      throw error;
    });
};

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
export const getAllUsersAPI = () => {
  return axios
    .get(CHAT_USERS_BASE_URL + "/users")
    .then((response) => {
      return response.data; 
    })
    .catch((error) => {
      console.error("Error fetching users:", error);
      throw error;
    });
};