import { renderHook, act } from '@testing-library/react';
import { useAnimation } from '../../hooks/useAnimation';
import { vi } from 'vitest';

describe('useAnimation', () => {
  beforeEach(() => {
    vi.useFakeTimers(); // Mock timers
  });

  afterEach(() => {
    vi.useRealTimers(); // Restore timers after each test
  });

  it('should initialize with isAnimating as false', () => {
    const { result } = renderHook(() => useAnimation());
    expect(result.current.isAnimating).toBe(false); // Initial state should be false
  });

  it('should set isAnimating to true when triggerAnimation is called', () => {
    const { result } = renderHook(() => useAnimation());

    act(() => {
      result.current.triggerAnimation(); // Trigger the animation
    });

    expect(result.current.isAnimating).toBe(true); // isAnimating should be true after trigger
  });

});
