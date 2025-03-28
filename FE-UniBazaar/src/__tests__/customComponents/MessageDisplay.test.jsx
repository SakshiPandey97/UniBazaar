import { render, screen, act } from '@testing-library/react';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { MessageDisplay } from '../../customComponents/MessageDisplay';

describe('MessageDisplay', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should display a message when messages is null', () => {
    const messages = null;
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };

    render(<MessageDisplay messages={messages} userId={userId} selectedUser={selectedUser} />);

    expect(screen.getByText('Start a conversation with User 2!')).toBeInTheDocument();
  });

  it('should scroll to the bottom when messages change', async () => {
    const messages = [
      { id: '1', sender_id: 1, receiver_id: 2, content: 'Hello', sender_name: 'User 1' },
    ];
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };
    const mockScrollIntoView = vi.fn();
    const originalScrollIntoView = Element.prototype.scrollIntoView;
    Element.prototype.scrollIntoView = mockScrollIntoView;

    const { rerender } = render(<MessageDisplay messages={messages} userId={userId} selectedUser={selectedUser} />);

    const newMessages = [
      ...messages,
      { id: '2', sender_id: 2, receiver_id: 1, content: 'Hi there', sender_name: 'User 2' },
    ];
    await act(async () => {
      rerender(<MessageDisplay messages={newMessages} userId={userId} selectedUser={selectedUser} />);
      await new Promise(resolve => setTimeout(resolve, 20));
    });

    expect(mockScrollIntoView).toHaveBeenCalled();
    Element.prototype.scrollIntoView = originalScrollIntoView;
  });
  it('should scroll to the bottom on initial render', async () => {
    const messages = [
      { id: '1', sender_id: 1, receiver_id: 2, content: 'Hello', sender_name: 'User 1' },
    ];
    const userId = '1';
    const selectedUser = { id: 2, name: 'User 2' };
    const mockScrollIntoView = vi.fn();
    const originalScrollIntoView = Element.prototype.scrollIntoView;
    Element.prototype.scrollIntoView = mockScrollIntoView;

    render(<MessageDisplay messages={messages} userId={userId} selectedUser={selectedUser} />);
    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 20));
    });

    expect(mockScrollIntoView).toHaveBeenCalled();
    Element.prototype.scrollIntoView = originalScrollIntoView;
  });
});
