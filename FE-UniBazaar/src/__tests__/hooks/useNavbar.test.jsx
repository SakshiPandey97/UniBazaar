import { renderHook, act } from '@testing-library/react';
import { vi } from 'vitest';
import useNavbar from '../../hooks/useNavbar';  // Adjust path as necessary
import { useUserAuth } from '../../hooks/useUserAuth';  // Adjust path as necessary
import { useNavigate } from 'react-router-dom';

// Mock the dependencies
vi.mock('../../hooks/useUserAuth', () => ({
  useUserAuth: vi.fn(),
}));

vi.mock('react-router-dom', () => ({
  useNavigate: vi.fn(),
}));

describe('useNavbar hook', () => {
  beforeEach(() => {
    // Reset mocks before each test
    vi.clearAllMocks();
  });

  it('should return initial menu and dropdown state', () => {
    const toggleLoginModal = vi.fn();
    const { result } = renderHook(() => useNavbar({ toggleLoginModal }));

    expect(result.current.isMenuOpen).toBe(false);
    expect(result.current.isDropdownOpen).toBe(false);
  });

  it('should toggle menu state when toggleMenu is called', () => {
    const toggleLoginModal = vi.fn();
    const { result } = renderHook(() => useNavbar({ toggleLoginModal }));

    // Initially, menu should be closed
    expect(result.current.isMenuOpen).toBe(false);

    // Toggle the menu state
    act(() => {
      result.current.toggleMenu();
    });

    // After toggling, menu should be open
    expect(result.current.isMenuOpen).toBe(true);
  });

  it('should toggle dropdown state when toggleDropdown is called', () => {
    const toggleLoginModal = vi.fn();
    const { result } = renderHook(() => useNavbar({ toggleLoginModal }));

    // Initially, dropdown should be closed
    expect(result.current.isDropdownOpen).toBe(false);

    // Toggle the dropdown state
    act(() => {
      result.current.toggleDropdown();
    });

    // After toggling, dropdown should be open
    expect(result.current.isDropdownOpen).toBe(true);
  });


  it('should call toggleLoginModal when user is not authenticated in handleNavigation', () => {
    const toggleLoginModal = vi.fn();
    const navigate = vi.fn();
    useUserAuth.mockReturnValue({ userState: false });

    const { result } = renderHook(() => useNavbar({ toggleLoginModal }));

    // Call handleNavigation
    act(() => {
      result.current.handleNavigation('/some-path');
    });

    // Expect toggleLoginModal to be called
    expect(toggleLoginModal).toHaveBeenCalled();
  });

  it('should call handleAuthAction and toggle login modal if user is not authenticated', () => {
    const toggleLoginModal = vi.fn();
    useUserAuth.mockReturnValue({ userState: false });

    const { result } = renderHook(() => useNavbar({ toggleLoginModal }));

    // Call handleAuthAction when user is not authenticated
    act(() => {
      result.current.handleAuthAction();
    });

    // Expect toggleLoginModal to be called
    expect(toggleLoginModal).toHaveBeenCalled();
  });

  it('should call handleAuthAction and toggle user login if user is authenticated', () => {
    const toggleLoginModal = vi.fn();
    const userAuthMock = { userState: true, toggleUserLogin: vi.fn(), setUserID: vi.fn() };
    useUserAuth.mockReturnValue(userAuthMock);

    const { result } = renderHook(() => useNavbar({ toggleLoginModal }));

    // Call handleAuthAction when user is authenticated
    act(() => {
      result.current.handleAuthAction();
    });

    // Expect toggleUserLogin and setUserID to be called
    expect(userAuthMock.toggleUserLogin).toHaveBeenCalled();
    expect(userAuthMock.setUserID).toHaveBeenCalledWith('');
  });
});
