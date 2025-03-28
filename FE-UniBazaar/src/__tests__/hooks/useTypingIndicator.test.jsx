import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useTypingIndicator } from '../../hooks/useTypingIndicator';

describe('useTypingIndicator', () => {
  let setInput;

  beforeEach(() => {
    vi.clearAllMocks();
    setInput = vi.fn();
  });

  it('should update setInput with the input value', () => {
    const { result } = renderHook(() => useTypingIndicator(setInput));
    const handleTyping = result.current;

    const mockEvent = {
      target: {
        value: 'Hello, world!',
      },
    };

    act(() => {
      handleTyping(mockEvent);
    });

    expect(setInput).toHaveBeenCalledWith('Hello, world!');
  });

  it('should update setInput with an empty string', () => {
    const { result } = renderHook(() => useTypingIndicator(setInput));
    const handleTyping = result.current;

    const mockEvent = {
      target: {
        value: '',
      },
    };

    act(() => {
      handleTyping(mockEvent);
    });

    expect(setInput).toHaveBeenCalledWith('');
  });
});
