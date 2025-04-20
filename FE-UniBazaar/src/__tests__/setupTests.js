import '@testing-library/jest-dom/vitest';
import { vi } from 'vitest'; // Import vi for mocking

// Mock window.matchMedia for libraries like embla-carousel
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false, // Default value, usually sufficient for tests
    media: query,
    onchange: null,
    addListener: vi.fn(), // Deprecated but sometimes used
    removeListener: vi.fn(), // Deprecated but sometimes used
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

// You might also need mocks for IntersectionObserver and ResizeObserver
// if embla-carousel or other components use them. Uncomment if needed.

// --- FIX: Uncomment the IntersectionObserver mock ---
const IntersectionObserverMock = vi.fn(() => ({
  disconnect: vi.fn(),
  observe: vi.fn(),
  takeRecords: vi.fn(),
  unobserve: vi.fn(),
}));
vi.stubGlobal('IntersectionObserver', IntersectionObserverMock);

const ResizeObserverMock = vi.fn(() => ({
  disconnect: vi.fn(),
  observe: vi.fn(),
  unobserve: vi.fn(),
}));
vi.stubGlobal('ResizeObserver', ResizeObserverMock);
