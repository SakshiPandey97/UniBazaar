import { useEffect, useState } from "react";
import { useUserAuth } from "@/hooks/useUserAuth";
import { userProfileDetailsAPI } from "@/api/userAxios";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

function ViewMyProfilePage() {
  const { userID } = useUserAuth();
  const [profile, setProfile] = useState({ name: "", email: "", phone: "" });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchProfile = async () => {
      setLoading(true);
      setError("");

      try {
        if (!userID) throw new Error("No user ID found");

        const res = await userProfileDetailsAPI(userID);

        setProfile({
          name: res.name,
          email: res.email,
          phone: res.phone, // ✅ added phone field
        });
      } catch (err) {
        console.error("Failed to fetch profile:", err);
        setError("Failed to load profile");
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, [userID]);

  return (
    <div className="w-full h-full flex justify-center items-center">
      <div className="flex flex-col w-[600px] h-[500px] loginDiv justify-center">
        <h1 className="font-mono text-5xl font-bold text-[#008080] flex justify-around">
          Profile Details
        </h1>

        <div className="w-full h-15 flex justify-center mt-2 mb-2">
          <Avatar className="h-full w-20">
            <AvatarImage src="https://github.com/shadcn.png" />
            <AvatarFallback>U</AvatarFallback>
          </Avatar>
        </div>

        {loading ? (
          <p className="text-center text-sm">Loading...</p>
        ) : error ? (
          <p className="text-center text-red-500">{error}</p>
        ) : (
          <>
            {/* Name */}
            <div className="w-full flex flex-col">
              <Label htmlFor="name" className="ml-25">Name</Label>
              <div className="w-full flex justify-center">
                <Input
                  type="text"
                  id="name"
                  placeholder="Name"
                  className="w-2/3"
                  value={profile.name}
                  readOnly
                />
              </div>
            </div>

            {/* Email */}
            <div className="w-full flex flex-col">
              <Label htmlFor="email" className="ml-25 mt-2">Email</Label>
              <div className="w-full flex justify-center">
                <Input
                  type="email"
                  id="email"
                  placeholder="Email"
                  className="w-2/3"
                  value={profile.email}
                  readOnly
                />
              </div>
            </div>

            {/* Phone Number */}
            <div className="w-full flex flex-col">
              <Label htmlFor="phone" className="ml-25 mt-2">Phone</Label>
              <div className="w-full flex justify-center">
                <Input
                  type="text"
                  id="phone"
                  placeholder="Phone"
                  className="w-2/3"
                  value={profile.phone}
                  readOnly
                />
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default ViewMyProfilePage;
