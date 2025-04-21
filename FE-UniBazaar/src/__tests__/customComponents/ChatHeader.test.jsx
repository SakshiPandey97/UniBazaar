import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import ChatHeader from '../../customComponents/Chat/ChatHeader';

describe('ChatHeader', () => {
  it('should render the user name correctly', () => {
    render(<ChatHeader name="John Doe" />);
    expect(screen.getByText('John Doe')).toBeInTheDocument();
  });

  it('should render "Chat" as the title if name is null', () => {
    render(<ChatHeader name={null} />);
    expect(screen.getByRole('heading', { name: 'Chat', level: 2 })).toBeInTheDocument();
  });

  it('should render "Chat" as the title if name is an empty string', () => {
    render(<ChatHeader name="" />);
    expect(screen.getByRole('heading', { name: 'Chat', level: 2 })).toBeInTheDocument();
  });


  it('should render correct initials for a single name', () => {
    render(<ChatHeader name="Alice" />);
    const avatar = screen.getByText('A');
    expect(avatar).toBeInTheDocument();
    expect(avatar).toHaveClass('rounded-full'); 
  });

  it('should render correct initials for multiple names', () => {
    render(<ChatHeader name="Bob Smith" />);
    const avatar = screen.getByText('BS');
    expect(avatar).toBeInTheDocument();
    expect(avatar).toHaveClass('rounded-full');
  });

   it('should render correct initials for names with extra spaces', () => {
    render(<ChatHeader name="  Charlie  Brown  " />);
    const avatar = screen.getByText('CB');
    expect(avatar).toBeInTheDocument();
    expect(avatar).toHaveClass('rounded-full');
  });


  it('should render "?" in the avatar for an empty name string', () => {
    render(<ChatHeader name="" />);
    const avatar = screen.getByText('?');
    expect(avatar).toBeInTheDocument();
    expect(avatar).toHaveClass('rounded-full'); 
  });

  it('should render "?" in the avatar for a null name', () => {
    render(<ChatHeader name={null} />);
    const avatar = screen.getByText('?');
    expect(avatar).toBeInTheDocument();
    expect(avatar).toHaveClass('rounded-full'); 
  });
});
