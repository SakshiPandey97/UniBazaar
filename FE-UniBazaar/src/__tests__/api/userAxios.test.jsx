import { describe, test, expect, vi } from "vitest";
import axios from "axios";
import {
  userLoginAPI,
  userRegisterAPI,
  userVerificationAPI,
} from "../../api/userAxios"; // Adjust import path if needed

// Mock axios
vi.mock("axios");

describe("User API Tests", () => {
  const mockUserLoginObject = { email: "test@example.com", password: "password123" };
  const mockUserRegisterObject = { name: "John Doe", email: "john@example.com", password: "password123" };
  const mockUserVerificationObject = { email: "john@example.com", otp: "123456" };

  afterEach(() => {
    vi.clearAllMocks();
  });

  test("userLoginAPI should return userId and store it in localStorage", async () => {
    const mockUserId = "12345";
    axios.post.mockResolvedValue({ data: { userId: mockUserId } });

    const result = await userLoginAPI({ userLoginObject: mockUserLoginObject });

    expect(result).toBe(mockUserId);
    expect(localStorage.getItem("userId")).toBe(mockUserId);
    expect(axios.post).toHaveBeenCalledWith(expect.stringContaining("/login"), mockUserLoginObject);
  });

  test("userLoginAPI should throw an error on failure", async () => {
    axios.post.mockRejectedValue(new Error("Login failed"));

    await expect(userLoginAPI({ userLoginObject: mockUserLoginObject })).rejects.toThrow("Login failed");
  });

  test("userRegisterAPI should return response data", async () => {
    const mockResponse = { message: "User registered successfully" };
    axios.post.mockResolvedValue({ data: mockResponse });

    const result = await userRegisterAPI({ userRegisterObject: mockUserRegisterObject });

    expect(result).toEqual(mockResponse);
    expect(axios.post).toHaveBeenCalledWith(expect.stringContaining("/signup"), mockUserRegisterObject);
  });

  test("userRegisterAPI should throw an error on failure", async () => {
    axios.post.mockRejectedValue(new Error("Registration failed"));

    await expect(userRegisterAPI({ userRegisterObject: mockUserRegisterObject })).rejects.toThrow("Registration failed");
  });

  test("userVerificationAPI should return userId", async () => {
    const mockUserId = "67890";
    axios.post.mockResolvedValue({ data: { userId: mockUserId } });

    const result = await userVerificationAPI({ userVerificationObject: mockUserVerificationObject });

    expect(result).toBe(mockUserId);
    expect(axios.post).toHaveBeenCalledWith(expect.stringContaining("/verifyEmail"), mockUserVerificationObject);
  });

  test("userVerificationAPI should throw an error on failure", async () => {
    axios.post.mockRejectedValue(new Error("Verification failed"));

    await expect(userVerificationAPI({ userVerificationObject: mockUserVerificationObject })).rejects.toThrow("Verification failed");
  });
});
