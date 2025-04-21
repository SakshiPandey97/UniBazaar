// src/__tests__/customComponents/Chat/ContactList.test.jsx
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { MemoryRouter } from 'react-router-dom'; // Needed because of useNavigate
import ContactList from '../../customComponents/Chat/ContactList'; // Adjust path

// Mock useNavigate
const mockedNavigate = vi.fn();
vi.mock('react-router-dom', async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    useNavigate: () => mockedNavigate,
  };
});

describe('ContactList', () => {
  const mockUsers = [
    { id: '1', name: 'User One' },
    { id: '2', name: 'User Two' },
    { id: '3', name: 'User Three' },
  ];
  const currentUserId = '1';
  let onSelect;

  beforeEach(() => {
    onSelect = vi.fn();
    mockedNavigate.mockClear();
  });

  const renderComponent = (props) => {
    return render(
      <MemoryRouter>
        <ContactList
          users={mockUsers}
          loading={false}
          currentUserId={currentUserId}
          selectedUserId={null}
          onSelect={onSelect}
          unreadSenders={new Set()}
          {...props}
        />
      </MemoryRouter>
    );
  };

  it('should render loading state', () => {
    renderComponent({ loading: true });
    expect(screen.getByText('Loading users...')).toBeInTheDocument();
  });

  it('should render "No other contacts found" when users array is empty (after filtering)', () => {
    renderComponent({ users: [{ id: '1', name: 'Current User Only' }] });
    expect(screen.getByText('No other contacts found.')).toBeInTheDocument();
  });

  it('should render list of users, filtering out the current user', () => {
    renderComponent();
    expect(screen.queryByText('User One')).not.toBeInTheDocument(); // Current user
    expect(screen.getByText('User Two')).toBeInTheDocument();
    expect(screen.getByText('User Three')).toBeInTheDocument();
  });

  it('should highlight the selected user', () => {
    renderComponent({ selectedUserId: '2' });
    const userTwoItem = screen.getByText('User Two').closest('li');
    expect(userTwoItem).toHaveClass('bg-blue-100');
  });

  it('should display notification dot for unread senders', () => {
    renderComponent({ unreadSenders: new Set(['3']) });
    const userThreeItem = screen.getByText('User Three').closest('li');
    const notificationDot = userThreeItem.querySelector('span[aria-label="Unread messages"]');
    expect(notificationDot).toBeInTheDocument();
    expect(notificationDot).toHaveClass('bg-red-500'); // Or blue, depending on your final style

    const userTwoItem = screen.getByText('User Two').closest('li');
    expect(userTwoItem.querySelector('span[aria-label="Unread messages"]')).not.toBeInTheDocument();
  });

  it('should call onSelect and navigate when a user is clicked', () => {
    renderComponent();
    const userTwoItem = screen.getByText('User Two').closest('li');
    fireEvent.click(userTwoItem);

    expect(onSelect).toHaveBeenCalledTimes(1);
    expect(onSelect).toHaveBeenCalledWith(mockUsers.find(u => u.id === '2'));
    expect(mockedNavigate).toHaveBeenCalledTimes(1);
    expect(mockedNavigate).toHaveBeenCalledWith('/messaging', { replace: true });
  });
});
