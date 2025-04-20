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
    // Ensure the base URL is defined
    if (!CHAT_USERS_BASE_URL) {
      console.error("Messaging service URL (VITE_CHAT_USERS_BASE_URL) is not defined.");
      return Promise.reject(new Error("Messaging service URL is not configured."));
    }
  
    // Make a POST request to the /api/users/sync endpoint
    return axios
      .post(CHAT_USERS_BASE_URL + "/api/users/sync", userData) // Send userData in the request body
      .then((response) => {
        console.log("User sync successful:", response.data); // Log success
        return response.data; // Return the response data
      })
      .catch((error) => {
        console.error("Error syncing user to messaging service:", error);
        // It might be okay for this to fail silently in the background,
        // so we log the error but don't necessarily re-throw it to break the login flow.
        // Depending on requirements, you might want to handle specific errors (e.g., 409 Conflict) differently.
        // throw error; // Uncomment if you need the calling code to know about the failure
      });
  };