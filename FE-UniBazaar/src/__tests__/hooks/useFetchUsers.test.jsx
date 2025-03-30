import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import useFetchUsers from '../../hooks/useFetchUsers';
import * as messagingAxios from '../../api/messagingAxios';

// Mock the getAllUsersAPI function
vi.mock('../../api/messagingAxios', () => ({
  getAllUsersAPI: vi.fn(),
}));

describe('useFetchUsers', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.useFakeTimers(); // Use fake timers for better control in other tests
  });

  afterEach(() => {
    vi.useRealTimers(); // Restore real timers after each test
  });

  it('should fetch users successfully', async () => {
    const mockUsers = [{ id: 1, name: 'User 1' }, { id: 2, name: 'User 2' }];
    messagingAxios.getAllUsersAPI.mockResolvedValue(mockUsers);

    const { result } = renderHook(() => useFetchUsers(1));

    expect(result.current.loadingUsers).toBe(true);
    expect(messagingAxios.getAllUsersAPI).toHaveBeenCalledWith(1);

    await act(async () => {
      // Wait for the promise to resolve
      await vi.runAllTimersAsync();
    });

    expect(result.current.loadingUsers).toBe(false);
    expect(result.current.users).toEqual(mockUsers);
  });

  it('should handle API error and retry', async () => {
    const mockError = new Error('API Error');
    messagingAxios.getAllUsersAPI.mockRejectedValueOnce(mockError).mockResolvedValueOnce([{ id: 3, name: 'User 3' }]);

    const { result } = renderHook(() => useFetchUsers(1));

    expect(result.current.loadingUsers).toBe(true);
    expect(messagingAxios.getAllUsersAPI).toHaveBeenCalledWith(1);

    // Wait for the first API call to fail and the retry to be scheduled
    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });

    // Wait for the retry to complete (5 seconds)
    await act(async () => {
      await vi.advanceTimersByTimeAsync(5000);
    });

    expect(result.current.loadingUsers).toBe(false);
    expect(result.current.users).toEqual([{ id: 3, name: 'User 3' }]);
    expect(messagingAxios.getAllUsersAPI).toHaveBeenCalledTimes(2);
  }, 10000); // Increase timeout for this test

  it('should handle multiple errors and retry', async () => {
    const mockError = new Error('API Error');
    messagingAxios.getAllUsersAPI.mockRejectedValue(mockError);

    const { result } = renderHook(() => useFetchUsers(1));

    expect(result.current.loadingUsers).toBe(true);
    expect(messagingAxios.getAllUsersAPI).toHaveBeenCalledWith(1);

    // Wait for the first API call to fail and the retry to be scheduled
    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });

    // Wait for the retry to complete (5 seconds)
    await act(async () => {
      await vi.advanceTimersByTimeAsync(5000);
    });
    expect(messagingAxios.getAllUsersAPI).toHaveBeenCalledTimes(2);
    expect(result.current.loadingUsers).toBe(false);
    expect(result.current.users).toEqual([]);
  }, 10000); // Increase timeout for this test
});
