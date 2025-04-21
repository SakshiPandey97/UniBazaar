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

  export const syncUserToMessagingAPI = (userData) => {
    if (!CHAT_USERS_BASE_URL) {
      console.error("Messaging service URL (VITE_CHAT_USERS_BASE_URL) is not defined.");
      return Promise.reject(new Error("Messaging service URL is not configured."));
    }
  
    return axios
      .post(CHAT_USERS_BASE_URL + "/api/users/sync", userData) 
      .then((response) => {
        console.log("User sync successful:", response.data);
        return response.data; 
      })
      .catch((error) => {
        console.error("Error syncing user to messaging service:", error);
      });
  };

  export const getUnreadSendersAPI = (userId) => { 
    if (!CHAT_USERS_BASE_URL) {
        console.error("Messaging service URL (VITE_CHAT_USERS_BASE_URL) is not defined.");
        return Promise.reject(new Error("Messaging service URL is not configured."));
    }
    if (!userId) {
        console.error("User ID is required to fetch unread senders.");
        return Promise.reject(new Error("User ID not provided."));
    }

    const url = `${CHAT_USERS_BASE_URL}/api/unread-senders?user_id=${userId}`;

    return axios
        .get(url) 
        .then((response) => {
            if (response.data && Array.isArray(response.data.senderIds)) {
                return response.data.senderIds;
            } else {
                console.warn("Received unexpected data format for unread senders:", response.data);
                return []; 
            }
        })
        .catch((error) => {
            console.error("Error fetching unread senders:", error);
            throw error;
        });
};