import { createContext, useContext, useState } from "react";

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  // Helper: Get userId from cookie
  const getUserIdFromCookie = () => {
    const match = document.cookie.match(/(?:^|;\s*)userId=([^;]*)/);
    return match ? decodeURIComponent(match[1]) : "";
  };

  const [userID, setUserID] = useState(getUserIdFromCookie());
  const [userState, setUserState] = useState(!!getUserIdFromCookie());

  const loginUser = (userData) => {
    // Extract the actual userId from the object
    const id = userData?.userId; // Use optional chaining just in case

    if (id) { // Only proceed if we have a valid ID
      setUserID(id); // Set state with the primitive ID
      setUserState(true);}

    const expires = new Date(
      Date.now() + 7 * 24 * 60 * 60 * 1000
    ).toUTCString();
    document.cookie = `userId=${encodeURIComponent(
      id
    )}; expires=${expires}; path=/`;
  };

  const logoutUser = () => {
    setUserID("");
    setUserState(false);
    document.cookie = "userId=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    const currentPath = location.pathname;
    if (
      currentPath.startsWith("/sell") ||
      currentPath.startsWith("/messaging")
    ) {
      navigate("/");
    }
  };

  return (
    <AuthContext.Provider value={{ userID, userState, loginUser, logoutUser }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useUserAuth = () => useContext(AuthContext);
