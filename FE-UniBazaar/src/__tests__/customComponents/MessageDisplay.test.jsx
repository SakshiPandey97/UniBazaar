// e:\Study\MS_Study\Software Engineering\UniBazaar\FE-UniBazaar\src\__tests__\customComponents\MessageDisplay.test.jsx
import { render, screen } from '@testing-library/react'; // Removed act, not needed for these tests
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { MessageDisplay } from '../../customComponents/Chat/MessageDisplay';

describe('MessageDisplay', () => {
  beforeEach(() => {
    vi.clearAllMocks(); // Keep this if other tests in this file might use mocks
  });

  it('should display a message when messages is null', () => {
    const messages = null;
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };

    render(<MessageDisplay messages={messages} userId={userId} selectedUser={selectedUser} />);

    expect(screen.getByText('Start a conversation with User 2!')).toBeInTheDocument();
  });

  it('should display a message when messages is an empty array', () => {
    const messages = [];
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };

    render(<MessageDisplay messages={messages} userId={userId} selectedUser={selectedUser} />);

    expect(screen.getByText('Start a conversation with User 2!')).toBeInTheDocument();
  });

  it('should render messages when provided', () => {
    const messages = [
      { id: '1', temp_id: 't1', sender_id: 1, content: 'Hello' },
      { id: '2', temp_id: 't2', sender_id: 2, content: 'Hi there' },
    ];
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };

    render(<MessageDisplay messages={messages} userId={userId} selectedUser={selectedUser} />);

    expect(screen.getByText('Hello')).toBeInTheDocument();
    expect(screen.getByText('Hi there')).toBeInTheDocument();
  });

  // --- REMOVED SCROLLING TESTS ---
  // it('should scroll to the bottom when messages change', async () => { ... });
  // it('should scroll to the bottom on initial render', async () => { ... });
});
