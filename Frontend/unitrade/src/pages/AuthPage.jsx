import { useAuthHandler } from "../hooks/useAuthHandler";
import UserLoginForm from "./components/UserLoginForm";
import UserRegisterForm from "./components/UserRegisterForm";

function AuthPage({toggleModal}) {
  const {
    isRegistering,
    isSubmitting,
    successMessage,
    handleSubmit,
    toggleAuthMode,
  } = useAuthHandler({toggleModal});

  return (
    <div className="w-full h-full flex justify-center items-center">
      <div className="flex flex-col w-[600px] h-[500px] loginDiv justify-center">
        <h1 className="ml-[50px] font-mono text-5xl font-bold text-[#008080]">
          {isRegistering ? "Sign Up" : "Login"}
        </h1>
        <p className="ml-[50px] mt-[40px] font-mono text-gray-500">
          {isRegistering ? "Create your account" : "Welcome Back!!"}
        </p>

        {successMessage && (
          <p className="text-green-600 font-mono text-center mt-2">{successMessage}</p>
        )}

        <div className="flex justify-center mt-[10px]">
          {isRegistering ? (
            <UserRegisterForm handleSubmit={handleSubmit} isSubmitting={isSubmitting} />
          ) : (
            <UserLoginForm handleSubmit={handleSubmit} isSubmitting={isSubmitting} />
          )}
        </div>

        <p className="font-mono flex flex-row justify-center mt-2">
          {isRegistering ? "Already have an account?" : "I don’t have an account?"}
          <span
            className="font-bold text-[#008080] cursor-pointer ml-1"
            onClick={toggleAuthMode}
          >
            {isRegistering ? "Login" : "Sign-Up"}
          </span>
        </p>
      </div>
    </div>
  );
}

export default AuthPage;
