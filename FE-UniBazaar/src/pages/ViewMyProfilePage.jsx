import UserRegisterForm from "../customComponents/UserRegisterForm";
import UserLoginForm from "../customComponents/UserLoginForm";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

function ViewMyProfilePage() {
  return (
    <div className="w-full h-full flex justify-center items-center">
      <div className="flex flex-col w-[600px] h-[500px] loginDiv justify-center">
        <h1 className="font-mono text-5xl font-bold text-[#008080] flex justify-around">
          Profile Details
        </h1>
        <div className="w-full h-15 flex justify-center mt-2 mb-2">
          <Avatar className="h-full w-20">
            <AvatarImage src="https://github.com/shadcn.png" />
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
        </div>

        <div className="w-full flex flex-col">
          <Label htmlFor="name" className="ml-25">
            Name
          </Label>
          <div className="w-full flex justify-center">
            <Input type="text" id="name" placeholder="Name" className="w-2/3" />
          </div>
        </div>
        <div className="w-full flex flex-col">
          <Label htmlFor="email" className="ml-25 mt-2">
            Email
          </Label>
          <div className="w-full flex justify-center">
            <Input
              type="email"
              id="email"
              placeholder="Email"
              className="w-2/3"
            />
          </div>
        </div>

        {/* {successMessage && (
          <p className="text-green-600 font-mono text-center mt-2">
            {successMessage}
          </p>
        )} */}

        {/* <div className="flex justify-center mt-[10px]">
              <UserRegisterForm
                handleSubmit={handleSubmit}
                isSubmitting={isSubmitting}
              />
          
        </div> */}

        {/* <p className="font-mono flex flex-row justify-center mt-2">
          Already have an account?"
           <span
            className="font-bold text-[#008080] cursor-pointer ml-1"
            onClick={toggleAuthMode}
          >
            {isRegistering ? "Login" : "Sign-Up"}
          </span>
        </p> */}
      </div>
    </div>
  );
}

export default ViewMyProfilePage;
