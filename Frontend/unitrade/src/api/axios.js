import axios from "axios";

const USER_BASE_URL = "http://127.0.0.1:4000";
const PRODUCT_BASE_URL = "https://unibazaar-products.azurewebsites.net";

export const userLoginAPI = ({ userLoginObject }) => {
  return axios
    .post(USER_BASE_URL + "/login", userLoginObject)
    .then((response) => {
      return response.data.userId;
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

export const getAllProductsAPI = () => {
  return axios
    .get(PRODUCT_BASE_URL + "/products")
    .then((response) => {
      console.error("products:", response);
      return response.data;
    })
    .catch((error) => {
      console.error("Error fetching products:", error);
      throw error;
    });
};