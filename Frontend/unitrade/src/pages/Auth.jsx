import { useNavigate, Navigate, useRoutes } from "react-router-dom";
import LoginPage from "./LoginPage";
import RegisterPage from "./RegisterPage";

function Auth() {
  const navigate = useNavigate();

  const handleGotoRegisterPage = () => {
    navigate("/auth/register");
  };

  const handleGoBack = () => {
    navigate("/");
  };

  const routes = useRoutes([
    { path: "/", element: <Navigate to="/auth/login" replace /> },
    { path: "login", element: <LoginPage handleGotoRegisterPage={handleGotoRegisterPage} /> },
    { path: "register", element: <RegisterPage /> },
  ]);

  return (
    <div className="modal-overlay" onClick={handleGoBack}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        {routes}
      </div>
    </div>
  );
}

export default Auth;
