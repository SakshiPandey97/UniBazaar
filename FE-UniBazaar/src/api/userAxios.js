import axios from "axios";

const USER_BASE_URL = import.meta.env.VITE_USER_BASE_URL;

export const userLoginAPI = ({ userLoginObject }) => {
  return axios
    .post(USER_BASE_URL + "/login", userLoginObject)
    .then((response) => {
      const userData = response.data;
      if (userData.userId) {
        localStorage.setItem("userId", userData.userId);
      }
      return userData;
    })
    .catch((error) => {
      console.error("Error logging in:", error);
      throw error.response?.data?.message || error;
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

export const userVerificationAPI = (userVerificationObject) => {
  return axios
    .post(USER_BASE_URL + "/verifyEmail", userVerificationObject)
    .then((response) => {
      console.log(response);
      return response.data;
    })
    .catch((error) => {
      console.error("Error Verifying user in:", error);
      throw error;
    });
};

export const userResendAPI = (userVerificationObject) => {
  console.log(userVerificationObject);
  return axios
    .post(USER_BASE_URL + "/resendOTP", userVerificationObject)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      console.error("Error resending OTP:", error);
      throw error;
    });
};

export const userProfileDetailsAPI = (userId) => {
  console.log(userId);
  return axios
    .get(`${USER_BASE_URL}/displayUser/${userId}`)
    .then((response) => response.data)
    .catch((error) => {
      console.error("Error fetching data", error);
      throw error;
    });
};

export const requestPasswordResetOTP = (email) => {
  return axios
    .get(`${USER_BASE_URL}/forgotPassword`, {
      params: { email },
    })
    .then((response) => {
      console.log("OTP request successful:", response.data);
      return true;
    })
    .catch((error) => {
      console.error(
        "Error sending OTP:",
        error.response?.data || error.message
      );
      return false;
    });
};


export const resetUserPassword = ({ email, otp, password }) => {
  return axios
    .post(`${USER_BASE_URL}/updatePassword`, {
      email,
      otp_code: otp,
      new_password: password,
    })
    .then((res) => res.data)
    .catch((err) => {
      console.error("Password reset failed:", err.response?.data || err.message);
      throw err;
    });
};