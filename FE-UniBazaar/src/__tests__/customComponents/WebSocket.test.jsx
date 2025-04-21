// src/__tests__/customComponents/WebSocket.test.jsx
import { renderHook } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import useWebSocket from '@/customComponents/WebsocketConnection'; // Adjust path

// --- Mock the global WebSocket ---
const mockWebSocketInstance = {
  onopen: null,
  onmessage: null,
  onerror: null,
  onclose: null,
  close: vi.fn(),
  readyState: WebSocket.OPEN, // Or CLOSED depending on test case
};
const mockWebSocket = vi.fn(() => mockWebSocketInstance);
global.WebSocket = mockWebSocket;

// --- Mock Environment Variable ---
// Do this in your test setup file or at the top of this file if preferred
// vi.stubEnv('VITE_CHAT_USERS_WS_URL', 'wss://unibazaar-messaging.azurewebsites.net');
// OR, if you haven't stubbed it globally, get it directly for the assertion:
const EXPECTED_WS_URL_BASE = 'wss://unibazaar-messaging.azurewebsites.net'; // Or import.meta.env.VITE_CHAT_USERS_WS_URL if setup allows

describe('useWebSocket', () => {
  const userId = '123';
  const onMessageReceived = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    // Reset mock instance properties if needed
    mockWebSocketInstance.close.mockClear();
  });

  it('should connect to WebSocket when userId is provided', () => {
    renderHook(() => useWebSocket(userId, onMessageReceived));

    // --- Updated Assertion ---
    const expectedUrl = `${EXPECTED_WS_URL_BASE}/ws?user_id=${userId}`;
    expect(mockWebSocket).toHaveBeenCalledTimes(1);
    expect(mockWebSocket).toHaveBeenCalledWith(expectedUrl);
  });

  // ... other tests ...
});
