// src/__tests__/customComponents/Chat/ChatPanel.test.jsx
import React from 'react';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { MemoryRouter, useLocation } from 'react-router-dom';
import ChatPanel from '@/customComponents/Chat/ChatPanel';

vi.mock('@/hooks/useFetchMessages', () => ({
  useFetchMessages: vi.fn(),
}));
import { useFetchMessages } from '@/hooks/useFetchMessages'; 

const mockWsSend = vi.fn();
const mockWsClose = vi.fn();
const mockWebSocketRef = { current: { send: mockWsSend, close: mockWsClose, readyState: 1 } };
vi.mock('@/customComponents/WebsocketConnection', () => ({ 
  default: vi.fn(() => mockWebSocketRef),
}));
import useWebSocket from '@/customComponents/WebsocketConnection'; 

const mockSendMessageFn = vi.fn();
vi.mock('@/hooks/useSendMessage', () => ({ 
  default: vi.fn(() => mockSendMessageFn),
}));
import useSendMessage from '@/hooks/useSendMessage'; 

const mockTypingHandler = vi.fn((setInput) => (e) => setInput(e.target.value));
vi.mock('@/hooks/useTypingIndicator', () => ({ 
  useTypingIndicator: vi.fn(() => mockTypingHandler),
}));
import { useTypingIndicator } from '@/hooks/useTypingIndicator'; 

vi.mock('@/utils/getUserId', () => ({ 
  getCurrentUserId: vi.fn(() => '1'), 
}));

vi.mock('react-router-dom', async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    useLocation: vi.fn(),
  };
});

vi.mock('@/customComponents/Chat/ChatHeader', () => ({ default: ({ name }) => <div>ChatHeader: {name}</div> }));
vi.mock('@/customComponents/Chat/MessageList', () => ({ default: ({ messages }) => <div data-testid="message-list">{messages.length} messages</div> }));
vi.mock('@/customComponents/Chat/MessageInput', () => ({ default: ({ input, onChange, onSend }) => (
    <div>
        <input data-testid="message-input" value={input} onChange={onChange} />
        <button onClick={onSend}>Send</button>
    </div>
)}));


describe('ChatPanel', () => {
  const mockUsers = [{ id: '1', name: 'User 1' }, { id: '2', name: 'User 2' }];
  const mockSelectedUser = { id: '2', name: 'User 2' };
  let setSelectedUser;

  beforeEach(() => {
    vi.clearAllMocks();
    setSelectedUser = vi.fn();
    useLocation.mockReturnValue({ search: '', pathname: '/messaging' });
    useFetchMessages.mockClear();
    useWebSocket.mockClear();
    useSendMessage.mockClear();
    useTypingIndicator.mockClear();
  });

  const renderComponent = (props) => {
    return render(
      <MemoryRouter>
        <ChatPanel
          users={mockUsers}
          selectedUser={null}
          setSelectedUser={setSelectedUser}
          {...props}
        />
      </MemoryRouter>
    );
  };

  it('should render placeholder when no user is selected', () => {
    renderComponent();
    expect(screen.getByText(/select a contact/i)).toBeInTheDocument();
    expect(screen.queryByTestId('message-list')).not.toBeInTheDocument();
  });

  it('should render chat components when a user is selected', () => {
    renderComponent({ selectedUser: mockSelectedUser });
    expect(screen.getByText(`ChatHeader: ${mockSelectedUser.name}`)).toBeInTheDocument();
    expect(screen.getByTestId('message-list')).toBeInTheDocument();
    expect(screen.getByTestId('message-input')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Send' })).toBeInTheDocument();
  });

  it('should call useFetchMessages with correct arguments when user is selected', () => {
    renderComponent({ selectedUser: mockSelectedUser });
    expect(useFetchMessages).toHaveBeenCalledWith('1', mockSelectedUser, expect.any(Function));
  });

   it('should call useWebSocket with correct arguments when user is selected', () => {
    renderComponent({ selectedUser: mockSelectedUser });
    expect(useWebSocket).toHaveBeenCalledWith('1', expect.any(Function));
  });

  it('should call useSendMessage with correct arguments', () => {
     renderComponent({ selectedUser: mockSelectedUser });
     expect(useSendMessage).toHaveBeenCalledWith(
        '1',
        mockSelectedUser,
        mockUsers,
        mockWebSocketRef,
        expect.any(String),
        expect.any(Function),
        expect.any(Function)
     );
  });

  it('should update input value on change', () => {
    renderComponent({ selectedUser: mockSelectedUser });
    const input = screen.getByTestId('message-input');
    fireEvent.change(input, { target: { value: 'New message' } });
    expect(mockTypingHandler).toHaveBeenCalled();
  });

   it('should call sendMessage function from useSendMessage hook on send button click', () => {
    renderComponent({ selectedUser: mockSelectedUser });
    const input = screen.getByTestId('message-input');
    fireEvent.change(input, { target: { value: 'Send this' } });
    const sendButton = screen.getByRole('button', { name: 'Send' });
    fireEvent.click(sendButton);
    expect(mockSendMessageFn).toHaveBeenCalledTimes(1);
  });

});
