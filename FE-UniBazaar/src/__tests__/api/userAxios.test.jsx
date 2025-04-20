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

  it("userRegisterAPI should return response data", async () => {
    const mockResponse = { status: "success" };
    const mockRegisterObj = { email: "new@example.com", password: "123456" };

    axios.post.mockResolvedValueOnce({ data: mockResponse });

    const result = await userRegisterAPI({ userRegisterObject: mockRegisterObj });

    expect(result).toEqual(mockResponse);
    expect(axios.post.mock.calls[0][0].endsWith("/signup")).toBe(true);
    expect(axios.post.mock.calls[0][1]).toEqual(mockRegisterObj);
  });
  it("should return userId on successful verification", async () => {
    const mockUserId = "verified123";
    const verifyData = { email: "verify@example.com", code: "123456" };

    // Mock axios.post to resolve to a data object with userId
    axios.post.mockResolvedValueOnce({ data: { userId: mockUserId } });

    const result = await userVerificationAPI(verifyData);

    expect(result.userId).toBe(mockUserId);
    expect(axios.post).toHaveBeenCalledWith(
      expect.stringMatching(/\/verifyEmail$/),
      verifyData
    );
  });

  it("should throw an error when verification fails", async () => {
    const errorMessage = "Request failed";
    axios.post.mockRejectedValueOnce(new Error(errorMessage));

    const verifyData = { email: "fail@example.com", code: "000000" };

    await expect(userVerificationAPI(verifyData)).rejects.toThrow(errorMessage);
    expect(axios.post).toHaveBeenCalledWith(
      expect.stringMatching(/\/verifyEmail$/),
      verifyData
    );
  });
});
