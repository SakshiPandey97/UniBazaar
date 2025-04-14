import React from "react";
import { Navigate } from "react-router-dom";
import { useUserAuth } from "@/hooks/useUserAuth";

const PrivateRoute = ({ children }) => {
  const { userState } = useUserAuth();

  if (!userState) {
    return <Navigate to="/" replace />;
  }

  return children;
};

export default PrivateRoute;
