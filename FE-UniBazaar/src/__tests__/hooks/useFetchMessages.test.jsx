import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { useFetchMessages } from '../../hooks/useFetchMessages';

describe('useFetchMessages', () => {
  let setMessages;
  const originalFetch = global.fetch;
  const originalConsoleError = console.error;

  beforeEach(() => {
    vi.clearAllMocks();
    vi.useFakeTimers();
    setMessages = vi.fn();
    global.fetch = vi.fn();
    console.error = vi.fn();
  });
  afterEach(() => {
    vi.restoreAllMocks();
    vi.useRealTimers();
    global.fetch = originalFetch;
    console.error = originalConsoleError;
  });

  it('should fetch messages successfully', async () => {
    const userId = '1';
    const selectedUser = { id: '2', name: 'User 2' };
    const mockMessages = [
      { id: '1', sender_id: 1, receiver_id: 2, content: 'Hello', sender_name: 'User 1' },
      { id: '2', sender_id: 2, receiver_id: 1, content: 'Hi there', sender_name: 'User 2' },
    ];
    global.fetch.mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue(mockMessages),
    });

    renderHook(() => useFetchMessages(userId, selectedUser, setMessages));

    expect(global.fetch).toHaveBeenCalledWith('http://127.0.0.1:8080/api/conversation/1/2');
    await act(async () => {
      await vi.runAllTimersAsync();
    });
    expect(setMessages).toHaveBeenCalledWith(mockMessages);
  });

  it('should handle API error', async () => {
    const userId = '1';
    const selectedUser = { id: '2', name: 'User 2' };
    const mockError = new Error('API Error');
    global.fetch.mockRejectedValue(mockError);
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    renderHook(() => useFetchMessages(userId, selectedUser, setMessages));
    await act(async () => {
      await vi.runAllTimersAsync();
    });

    expect(setMessages).toHaveBeenCalledWith([]);
    expect(consoleErrorSpy).toHaveBeenCalledWith('Failed to load messages:', mockError);
    consoleErrorSpy.mockRestore();
  });

  it('should handle API error with status not ok', async () => {
    const userId = '1';
    const selectedUser = { id: '2', name: 'User 2' };
    global.fetch.mockResolvedValue({
      ok: false,
      status: 404,
      json: vi.fn().mockResolvedValue({ message: 'Not Found' }),
    });
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    renderHook(() => useFetchMessages(userId, selectedUser, setMessages));
    await act(async () => {
      await vi.runAllTimersAsync();
    });

    expect(setMessages).toHaveBeenCalledWith([]);
    expect(consoleErrorSpy).toHaveBeenCalledWith('Failed to load messages:', new Error('Not Found'));
    consoleErrorSpy.mockRestore();
  });

  it('should handle empty array response', async () => {
    const userId = '1';
    const selectedUser = { id: '2', name: 'User 2' };
    global.fetch.mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue([]),
    });

    renderHook(() => useFetchMessages(userId, selectedUser, setMessages));
    await act(async () => {
      await vi.runAllTimersAsync();
    });

    expect(setMessages).toHaveBeenCalledWith([]);
  });

  it('should handle non-array response', async () => {
    const userId = '1';
    const selectedUser = { id: '2', name: 'User 2' };
    global.fetch.mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({ message: 'Not an array' }),
    });

    renderHook(() => useFetchMessages(userId, selectedUser, setMessages));
    await act(async () => {
      await vi.runAllTimersAsync();
    });

    expect(setMessages).toHaveBeenCalledWith([]);
  });

  it('should not fetch messages if selectedUser or userId is missing', () => {
    renderHook(() => useFetchMessages(null, { id: '2', name: 'User 2' }, setMessages));
    expect(global.fetch).not.toHaveBeenCalled();

    renderHook(() => useFetchMessages('1', null, setMessages));
    expect(global.fetch).not.toHaveBeenCalled();
  });
});
