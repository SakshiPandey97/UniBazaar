import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, beforeEach, afterEach, vi, beforeAll, afterAll } from 'vitest';

// Mock the global WebSocket object BEFORE importing anything that might use it
let mockWebSocket;
beforeAll(() => {
  mockWebSocket = vi.fn();
  mockWebSocket.OPEN = 1; // WebSocket.OPEN
  mockWebSocket.CLOSED = 3; // WebSocket.CLOSED
  mockWebSocket.prototype.send = vi.fn();
  mockWebSocket.prototype.close = vi.fn();
  mockWebSocket.prototype.onopen = null;
  mockWebSocket.prototype.onmessage = null;
  mockWebSocket.prototype.onerror = null;
  mockWebSocket.prototype.onclose = null;
  mockWebSocket.prototype.readyState = global.WebSocket.OPEN;

  // Mock the WebSocket constructor
  global.WebSocket = vi.fn((url) => {
    mockWebSocket(url); // Call the mock to track constructor calls
    return mockWebSocket.prototype; // Return the prototype to simulate an instance
  });
  global.WebSocket.OPEN = mockWebSocket.OPEN;
  global.WebSocket.CLOSED = mockWebSocket.CLOSED;
});

import useWebSocket from '../../customComponents/WebsocketConnection';

describe('useWebSocket', () => {
  let onMessageReceived;
  const originalConsoleError = console.error;
  const originalConsoleLog = console.log;

  beforeEach(() => {
    vi.clearAllMocks();

    onMessageReceived = vi.fn();
    console.error = vi.fn();
    console.log = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
    console.error = originalConsoleError;
    console.log = originalConsoleLog;
  });
  afterAll(() => {
    delete global.WebSocket;
  });

  it('should connect to WebSocket when userId is provided', () => {
    const userId = '123';
    renderHook(() => useWebSocket(userId, onMessageReceived));

    expect(mockWebSocket).toHaveBeenCalledWith(`ws://44.211.190.101:8080/ws?user_id=${userId}`);
  });

  it('should not connect to WebSocket when userId is not provided', () => {
    renderHook(() => useWebSocket(null, onMessageReceived));

    expect(mockWebSocket).not.toHaveBeenCalled();
  });

  it('should call onMessageReceived when a message is received', () => {
    const userId = '123';
    renderHook(() => useWebSocket(userId, onMessageReceived));

    const mockMessage = { data: JSON.stringify({ message: 'Hello' }) };
    act(() => {
      mockWebSocket.prototype.onmessage(mockMessage);
    });

    expect(onMessageReceived).toHaveBeenCalledWith({ message: 'Hello' });
  });

  it('should log an error when an error occurs', () => {
    const userId = '123';
    renderHook(() => useWebSocket(userId, onMessageReceived));

    const mockError = new Error('WebSocket error');
    act(() => {
      mockWebSocket.prototype.onerror(mockError);
    });

    expect(console.error).toHaveBeenCalledWith('WebSocket Error:', mockError);
  });

  it('should log when the WebSocket is opened', () => {
    const userId = '123';
    renderHook(() => useWebSocket(userId, onMessageReceived));

    act(() => {
      mockWebSocket.prototype.onopen();
    });

    expect(console.log).toHaveBeenCalledWith(`Connected as User ${userId}`);
  });

  it('should log when the WebSocket is closed', () => {
    const userId = '123';
    renderHook(() => useWebSocket(userId, onMessageReceived));

    const mockEvent = { code: 1000, reason: 'Normal Closure' };
    act(() => {
      mockWebSocket.prototype.onclose(mockEvent);
    });

    expect(console.log).toHaveBeenCalledWith('WebSocket closed:', mockEvent);
  });

  it('should close the WebSocket on unmount', () => {
    const userId = '123';
    const { unmount } = renderHook(() => useWebSocket(userId, onMessageReceived));

    unmount();

    expect(mockWebSocket.prototype.close).toHaveBeenCalled();
  });

  it('should not close the WebSocket if it is not open on unmount', () => {
    const userId = '123';
    mockWebSocket.prototype.readyState = global.WebSocket.CLOSED;
    const { unmount } = renderHook(() => useWebSocket(userId, onMessageReceived));

    unmount();

    expect(mockWebSocket.prototype.close).not.toHaveBeenCalled();
  });
  it('should log when the WebSocket is closing', () => {
    const userId = '123';
    const { unmount } = renderHook(() => useWebSocket(userId, onMessageReceived));

    unmount();

  });
});
