import { useState, useEffect } from "react";
import { getAllUsersAPI } from "../api/messagingAxios";

const useFetchUsers = (userId) => {
  const [users, setUsers] = useState([]);
  const [loadingUsers, setLoadingUsers] = useState(true);

  useEffect(() => {
    if (userId) {
      const fetchUsers = async () => {
        try {
          const data = await getAllUsersAPI(userId);
          setUsers(data);
        } catch (error) {
          console.error("Failed to fetch users:", error);
          setTimeout(fetchUsers, 5000);
        } finally {
          setLoadingUsers(false);
        }
      };

      fetchUsers();
    }
  }, [userId]);

  return { users, loadingUsers };
};

export default useFetchUsers;