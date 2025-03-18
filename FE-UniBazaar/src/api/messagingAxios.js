import axios from "axios";

const CHAT_USERS_BASE_URL = import.meta.env.VITE_CHAT_USERS_BASE_URL;

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