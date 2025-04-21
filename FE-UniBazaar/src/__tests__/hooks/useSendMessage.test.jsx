import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import useSendMessage from '../../hooks/useSendMessage';
import { v4 as uuidv4 } from 'uuid';
import { toast } from 'react-toastify';

vi.mock('uuid', () => ({
  v4: vi.fn(() => 'mock-uuid'),
}));

vi.mock('react-toastify', () => ({
  toast: {
    error: vi.fn(),
  },
}));

describe('useSendMessage', () => {
  let mockWs;
  let setInput;
  let setMessages;
  let users;
  const mockTimestamp = 1743124747830;

  beforeEach(() => {
    vi.clearAllMocks();
    vi.spyOn(Date, 'now').mockReturnValue(mockTimestamp);
    mockWs = {
      current: {
        readyState: WebSocket.OPEN,
        send: vi.fn(),
      },
    };
    setInput = vi.fn();
    setMessages = vi.fn();
    users = [
      { id: 1, name: 'User 1' },
      { id: 2, name: 'User 2' },
    ];
    toast.error.mockClear();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should send a message successfully', async () => {
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };
    const input = 'Hello, User 2!';

    const { result } = renderHook(() =>
      useSendMessage(userId, selectedUser, users, mockWs, input, setInput, setMessages)
    );

    const sendMessage = result.current;

    await act(async () => {
      sendMessage();
    });

    expect(mockWs.current.send).toHaveBeenCalledTimes(1);
    expect(mockWs.current.send).toHaveBeenCalledWith(
      JSON.stringify({
        ID: 'mock-uuid',
        sender_id: 1,
        receiver_id: 2,
        content: 'Hello, User 2!',
        timestamp: mockTimestamp,
        read: false,
        sender_name: 'User 1',
      })
    );
    expect(setInput).toHaveBeenCalledWith('');
    expect(toast.error).not.toHaveBeenCalled();
  });

  it('should show toast error if userId is missing', async () => {
    const userId = null;
    const selectedUser = { id: 2, name: 'User 2' };
    const input = 'Hello, User 2!';

    const { result } = renderHook(() =>
      useSendMessage(userId, selectedUser, users, mockWs, input, setInput, setMessages)
    );
    const sendMessage = result.current;

    await act(async () => {
      sendMessage();
    });

    expect(mockWs.current.send).not.toHaveBeenCalled();
    expect(toast.error).toHaveBeenCalledTimes(1);
    expect(toast.error).toHaveBeenCalledWith('Please select a user to chat with!');
  });

  it('should show toast error if selectedUser is missing', async () => {
    const userId = '1';
    const selectedUser = null;
    const input = 'Hello, User 2!';

    const { result } = renderHook(() =>
      useSendMessage(userId, selectedUser, users, mockWs, input, setInput, setMessages)
    );
    const sendMessage = result.current;

    await act(async () => {
      sendMessage();
    });

    expect(mockWs.current.send).not.toHaveBeenCalled();
    expect(toast.error).toHaveBeenCalledTimes(1);
    expect(toast.error).toHaveBeenCalledWith('Please select a user to chat with!');
  });

  it('should not send a message if input is empty', async () => {
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };
    const input = '   ';

    const { result } = renderHook(() =>
      useSendMessage(userId, selectedUser, users, mockWs, input, setInput, setMessages)
    );
    const sendMessage = result.current;

    await act(async () => {
      sendMessage();
    });

    expect(mockWs.current.send).not.toHaveBeenCalled();
    expect(toast.error).not.toHaveBeenCalled();
  });

  it('should not send a message if WebSocket is not ready', async () => {
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };
    const input = 'Hello, User 2!';
    mockWs.current.readyState = WebSocket.CLOSED;

    const { result } = renderHook(() =>
      useSendMessage(userId, selectedUser, users, mockWs, input, setInput, setMessages)
    );
    const sendMessage = result.current;

    await act(async () => {
      sendMessage();
    });

    expect(mockWs.current.send).not.toHaveBeenCalled();
    expect(toast.error).not.toHaveBeenCalled();
  });

  it('should handle error when sending message', async () => {
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };
    const input = 'Hello, User 2!';
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    const sendError = new Error('Simulated WebSocket error');
    mockWs.current.send.mockImplementation(() => {
      throw sendError;
    });

    const { result } = renderHook(() =>
      useSendMessage(userId, selectedUser, users, mockWs, input, setInput, setMessages)
    );
    const sendMessage = result.current;

    await act(async () => {
      sendMessage();
    });

    expect(consoleErrorSpy).toHaveBeenCalledWith('Error sending message:', sendError);
    consoleErrorSpy.mockRestore();
  });
});
