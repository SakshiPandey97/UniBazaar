import Cookies from "js-cookie";

export function getCurrentUserId() {
    const userId = Cookies.get("userId");
    return userId || "";
}
