// src/__tests__/hooks/useFetchMessages.test.jsx
import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useFetchMessages } from '@/hooks/useFetchMessages'; // Adjust path

// --- Mock fetch ---
global.fetch = vi.fn();

// --- Mock Environment Variable ---
// vi.stubEnv('VITE_CHAT_USERS_BASE_URL', 'https://unibazaar-messaging.azurewebsites.net');
// OR:
const EXPECTED_API_URL_BASE = 'https://unibazaar-messaging.azurewebsites.net'; // Or import.meta.env.VITE_CHAT_USERS_BASE_URL

describe('useFetchMessages', () => {
  const userId = '1';
  const selectedUser = { id: '2', name: 'Test User' };
  const setMessages = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    fetch.mockClear(); // Clear fetch mocks
    fetch.mockResolvedValue({ // Default successful mock
        ok: true,
        json: async () => ([{ id: 'msg1', text: 'Hello' }]),
    });
  });

   afterEach(() => {
     vi.restoreAllMocks(); // Restore fetch if needed elsewhere
   });


  it('should fetch messages successfully', async () => {
    renderHook(() => useFetchMessages(userId, selectedUser, setMessages));

    // --- Updated Assertion ---
    const expectedUrl = `${EXPECTED_API_URL_BASE}/api/conversation/${userId}/${selectedUser.id}`;
    expect(global.fetch).toHaveBeenCalledTimes(1);
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl);

    // Wait for the fetch promise to resolve and state update to happen
    await act(async () => {
       // Allow promises microtasks to resolve
       await Promise.resolve();
    });

    expect(setMessages).toHaveBeenCalledWith([{ id: 'msg1', text: 'Hello' }]);
  });

  // ... other tests ...
});
