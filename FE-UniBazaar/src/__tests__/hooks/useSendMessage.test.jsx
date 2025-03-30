import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import useSendMessage from '../../hooks/useSendMessage';
import { v4 as uuidv4 } from 'uuid';

// Mock uuidv4 to return a consistent value for testing
vi.mock('uuid', () => ({
  v4: vi.fn(() => 'mock-uuid'),
}));

describe('useSendMessage', () => {
  let mockWs;
  let setInput;
  let setMessages;
  let users;
  const mockTimestamp = 1743124747830; // Consistent timestamp for testing

  beforeEach(() => {
    vi.clearAllMocks();

    // Mock Date.now()
    vi.spyOn(Date, 'now').mockReturnValue(mockTimestamp);

    // Mock WebSocket
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

    await act(async () => {
      result.current();
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
  });

  it('should not send a message if userId is missing', async () => {
    const userId = null;
    const selectedUser = { id: 2, name: 'User 2' };
    const input = 'Hello, User 2!';

    const { result } = renderHook(() =>
      useSendMessage(userId, selectedUser, users, mockWs, input, setInput, setMessages)
    );

    // Mock alert
    const originalAlert = window.alert;
    window.alert = vi.fn();

    await act(async () => {
      result.current();
    });

    expect(mockWs.current.send).not.toHaveBeenCalled();
    expect(window.alert).toHaveBeenCalledWith('Please select a user to chat with!');

    // Restore alert
    window.alert = originalAlert;
  });

  it('should not send a message if selectedUser is missing', async () => {
    const userId = '1';
    const selectedUser = null;
    const input = 'Hello, User 2!';

    const { result } = renderHook(() =>
      useSendMessage(userId, selectedUser, users, mockWs, input, setInput, setMessages)
    );

    // Mock alert
    const originalAlert = window.alert;
    window.alert = vi.fn();

    await act(async () => {
      result.current();
    });

    expect(mockWs.current.send).not.toHaveBeenCalled();
    expect(window.alert).toHaveBeenCalledWith('Please select a user to chat with!');

    // Restore alert
    window.alert = originalAlert;
  });

  it('should not send a message if input is empty', async () => {
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };
    const input = '   ';

    const { result } = renderHook(() =>
      useSendMessage(userId, selectedUser, users, mockWs, input, setInput, setMessages)
    );

    await act(async () => {
      result.current();
    });

    expect(mockWs.current.send).not.toHaveBeenCalled();
  });

  it('should not send a message if WebSocket is not ready', async () => {
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };
    const input = 'Hello, User 2!';
    mockWs.current.readyState = WebSocket.CLOSED;

    const { result } = renderHook(() =>
      useSendMessage(userId, selectedUser, users, mockWs, input, setInput, setMessages)
    );

    await act(async () => {
      result.current();
    });

    expect(mockWs.current.send).not.toHaveBeenCalled();
  });
  it('should handle error when sending message', async () => {
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };
    const input = 'Hello, User 2!';
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    mockWs.current.send.mockImplementation(() => {
      throw new Error('Simulated WebSocket error');
    });

    const { result } = renderHook(() =>
      useSendMessage(userId, selectedUser, users, mockWs, input, setInput, setMessages)
    );

    await act(async () => {
      result.current();
    });

    expect(consoleErrorSpy).toHaveBeenCalledWith('Error sending message:', new Error('Simulated WebSocket error'));
    consoleErrorSpy.mockRestore();
  });
});
