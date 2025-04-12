import { describe, it, expect, vi, beforeEach } from "vitest";
import axios from "axios";
import {
  userLoginAPI,
  userRegisterAPI,
  userVerificationAPI,
} from "@/api/userAxios";

vi.mock("axios");

describe("API Functions", () => {
  beforeEach(() => {
    vi.resetAllMocks();
    localStorage.clear();
  });

  it("userLoginAPI should return userId and store in localStorage", async () => {
    const mockUserId = "user123";
    const mockLoginObj = { email: "test@example.com", password: "pass" };

    axios.post.mockResolvedValueOnce({ data: { userId: mockUserId } });

    const result = await userLoginAPI({ userLoginObject: mockLoginObj });

    expect(result).toBe(mockUserId);
    expect(localStorage.getItem("userId")).toBe(mockUserId);

    // ✅ Check path ending only (less brittle)
    expect(axios.post.mock.calls[0][0].endsWith("/login")).toBe(true);
    expect(axios.post.mock.calls[0][1]).toEqual(mockLoginObj);
  });

  it("userRegisterAPI should return response data", async () => {
    const mockResponse = { status: "success" };
    const mockRegisterObj = { email: "new@example.com", password: "123456" };

    axios.post.mockResolvedValueOnce({ data: mockResponse });

    const result = await userRegisterAPI({ userRegisterObject: mockRegisterObj });

    expect(result).toEqual(mockResponse);
    expect(axios.post.mock.calls[0][0].endsWith("/signup")).toBe(true);
    expect(axios.post.mock.calls[0][1]).toEqual(mockRegisterObj);
  });

  it("userVerificationAPI should return userId", async () => {
    const mockUserId = "verified123";
    const verifyData = { email: "verify@example.com", code: "123456" };

    axios.post.mockResolvedValueOnce({ data: { userId: mockUserId } });

    const result = await userVerificationAPI(verifyData);

    expect(result).toBe(mockUserId);
    expect(axios.post.mock.calls[0][0].endsWith("/verifyEmail")).toBe(true);
    expect(axios.post.mock.calls[0][1]).toEqual(verifyData);
  });

  it("userLoginAPI should throw error on failure", async () => {
    axios.post.mockRejectedValueOnce(new Error("Request failed"));

    await expect(
      userLoginAPI({ userLoginObject: { email: "fail", password: "wrong" } })
    ).rejects.toThrow("Request failed");
  });
});
