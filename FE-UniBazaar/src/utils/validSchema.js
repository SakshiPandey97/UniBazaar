
import * as Yup from "yup";

export const validationSchema = Yup.object({
  email: Yup.string().email(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/,"Invalid email address").required("Email is required"),
  password: Yup.string()
    .min(6, "At least 6 characters")
    .matches(/[A-Z]/, "At least one uppercase letter")
    .matches(/[a-z]/, "At least one lowercase letter")
    .matches(/\d/, "At least one number")
    .matches(/[@$!%*?&]/, "At least one special character (@, $, !, %, *, ?, &)")
    .required("Password is required"),
});