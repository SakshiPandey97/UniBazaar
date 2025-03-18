export const getCurrentUserId = () => {
    return localStorage.getItem("userId") || null;
  };
  