import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import MessageList from '@/customComponents/Chat/MessageList';


describe('MessageList', () => {
  const userId = '1'; 
  const selectedUser = { id: '2', name: 'User Two' };
  const mockMessages = [
    { id: 'm1', sender_id: 1, receiver_id: 2, content: 'Hello', created_at: new Date().toISOString() },
    { id: 'm2', sender_id: 2, receiver_id: 1, content: 'Hi there', created_at: new Date().toISOString() },
    { id: 'm3', sender_id: 1, receiver_id: 2, content: 'How are you?', created_at: new Date().toISOString() },
  ];

  const renderComponent = (props) => {
    return render(
      <MessageList
        messages={[]}
        userId={userId}
        selectedUser={selectedUser}
        {...props}
      />
    );
  };

  it('should render the correct empty state message when messages array is empty', () => {
    renderComponent();
    expect(screen.getByText(`Start a conversation with ${selectedUser.name}!`)).toBeInTheDocument();
  });

  it('should render the content of each message', () => {
    renderComponent({ messages: mockMessages });
    expect(screen.getByText('Hello')).toBeInTheDocument();
    expect(screen.getByText('Hi there')).toBeInTheDocument();
    expect(screen.getByText('How are you?')).toBeInTheDocument();
  });

  it('should apply correct alignment based on message ownership', () => {
    renderComponent({ messages: mockMessages });

    const getMessageContainer = (text) => screen.getByText(text).closest('div[class*="justify-"]');

    const message1Container = getMessageContainer('Hello');
    expect(message1Container).toHaveClass('justify-end');

    const message2Container = getMessageContainer('Hi there');
    expect(message2Container).toHaveClass('justify-start');

    const message3Container = getMessageContainer('How are you?');
    expect(message3Container).toHaveClass('justify-end');
  });
});
