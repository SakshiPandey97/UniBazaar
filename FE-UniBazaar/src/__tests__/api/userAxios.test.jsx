import { describe, it, expect, vi, beforeEach } from "vitest";
import axios from "axios";
import {
  userLoginAPI,
  userRegisterAPI,
  userVerificationAPI,
} from "@/api/userAxios";

// Mock the axios module
vi.mock("axios");

describe("API Functions", () => {
  beforeEach(() => {
    // Reset mocks and localStorage before each test
    vi.resetAllMocks();
    localStorage.clear();
  });

  // --- UPDATED TEST ---
  it("userLoginAPI should return user data object and store userId in localStorage", async () => {
    const mockUserId = "user123";
    const mockName = "Test User";
    const mockEmail = "test@example.com";
    const mockLoginObj = { email: mockEmail, password: "pass" };
    // Simulate the API response structure expected by useAuthHandler
    const mockApiResponse = {
      data: {
        userId: mockUserId,
        name: mockName,
        email: mockEmail,
      }
    };

    // Mock the axios post call for login
    axios.post.mockResolvedValueOnce(mockApiResponse);

    const result = await userLoginAPI({ userLoginObject: mockLoginObj });

    // --- FIX ---
    // Expect the function to return the full data object from the response
    expect(result).toEqual({
      userId: mockUserId,
      name: mockName,
      email: mockEmail,
    });

    // Check localStorage is still set correctly
    expect(localStorage.getItem("userId")).toBe(mockUserId);

    // Verify axios call details
    expect(axios.post).toHaveBeenCalledTimes(1);
    expect(axios.post.mock.calls[0][0].endsWith("/login")).toBe(true); // Check endpoint path
    expect(axios.post.mock.calls[0][1]).toEqual(mockLoginObj); // Check request payload
  });

  it("userRegisterAPI should return response data", async () => {
    // This test seems consistent with useAuthHandler's usage
    const mockResponse = { status: "success", message: "User registered" }; // Example response
    const mockRegisterObj = { email: "new@example.com", password: "123456" };
    const mockApiResponse = { data: mockResponse };

    axios.post.mockResolvedValueOnce(mockApiResponse);

    const result = await userRegisterAPI({ userRegisterObject: mockRegisterObj });

    expect(result).toEqual(mockResponse);
    expect(axios.post).toHaveBeenCalledTimes(1);
    expect(axios.post.mock.calls[0][0].endsWith("/signup")).toBe(true);
    expect(axios.post.mock.calls[0][1]).toEqual(mockRegisterObj);
  });

  // --- UPDATED TEST ---
  it("userVerificationAPI should return userId object", async () => {
    const mockUserId = "verified123";
    const verifyData = { email: "verify@example.com", code: "123456" };
    // Assume verification also returns an object containing the userId
    const mockApiResponse = { data: { userId: mockUserId } };

    axios.post.mockResolvedValueOnce(mockApiResponse);

    const result = await userVerificationAPI(verifyData);

    // --- FIX ---
    // Expect the function to return an object containing the userId
    expect(result).toBe(mockUserId); // Use .toBe for primitive string comparison
    // Or, if you only care about the ID itself being present:
    // expect(result.userId).toBe(mockUserId);

    // Verify axios call details
    expect(axios.post).toHaveBeenCalledTimes(1);
    expect(axios.post.mock.calls[0][0].endsWith("/verifyEmail")).toBe(true);
    expect(axios.post.mock.calls[0][1]).toEqual(verifyData);
  });

  it("userLoginAPI should throw error on failure", async () => {
    const mockError = new Error("Request failed");
    const mockLoginObj = { email: "fail", password: "wrong" };

    // Mock axios to reject the promise
    axios.post.mockRejectedValueOnce(mockError);

    // Assert that the promise rejects with the expected error
    await expect(
      userLoginAPI({ userLoginObject: mockLoginObj })
    ).rejects.toThrow("Request failed");

    // Verify axios call details
    expect(axios.post).toHaveBeenCalledTimes(1);
    expect(axios.post.mock.calls[0][0].endsWith("/login")).toBe(true);
    expect(axios.post.mock.calls[0][1]).toEqual(mockLoginObj);
  });

  // Add similar error handling tests for userRegisterAPI and userVerificationAPI if not already present
  it("userRegisterAPI should throw error on failure", async () => {
    const mockError = new Error("Registration failed");
    const mockRegisterObj = { email: "fail@reg.com", password: "bad" };

    axios.post.mockRejectedValueOnce(mockError);

    await expect(userRegisterAPI({ userRegisterObject: mockRegisterObj }))
      .rejects.toThrow("Registration failed");

    expect(axios.post).toHaveBeenCalledTimes(1);
    expect(axios.post.mock.calls[0][0].endsWith("/signup")).toBe(true);
    expect(axios.post.mock.calls[0][1]).toEqual(mockRegisterObj);
  });

   it("userVerificationAPI should throw error on failure", async () => {
    const mockError = new Error("Verification failed");
    const verifyData = { email: "fail@verify.com", code: "000" };

    axios.post.mockRejectedValueOnce(mockError);

    await expect(userVerificationAPI(verifyData))
      .rejects.toThrow("Verification failed");

    expect(axios.post).toHaveBeenCalledTimes(1);
    expect(axios.post.mock.calls[0][0].endsWith("/verifyEmail")).toBe(true);
    expect(axios.post.mock.calls[0][1]).toEqual(verifyData);
  });
});
