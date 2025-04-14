import React, { useEffect } from "react";
import { Navigate } from "react-router-dom";
import { useUserAuth } from "@/hooks/useUserAuth";

const PrivateRoute = ({ children }) => {
  const { userState } = useUserAuth();

  useEffect(() => {
    if (!userState) {
      alert("Please login to perform this action");
    }
  }, [userState]);

  if (!userState) {
    return <Navigate to="/" replace />;
  }

  return children;
};

export default PrivateRoute;