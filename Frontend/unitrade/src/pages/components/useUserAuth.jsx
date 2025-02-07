import { createContext, useContext, useState } from "react";

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [userState, setUserState] = useState(false);

  const toggleUserLogin = () => {
    setUserState((prevState) => !prevState);
  };

  return (
    <AuthContext.Provider value={{ userState, toggleUserLogin }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useUserAuth = () => {
  return useContext(AuthContext);
};
