import axios from "axios";

const USER_BASE_URL = import.meta.env.VITE_USER_BASE_URL;


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

export const userVerificationAPI = ( userVerificationObject ) => {
  console.log(userVerificationObject)
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

