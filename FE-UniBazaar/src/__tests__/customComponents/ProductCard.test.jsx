import { render, screen, fireEvent } from '@testing-library/react';
import ProductCard from '../../customComponents/ProductCard';
import { generateStars } from '@/utils/generateStar';
import { MemoryRouter, useLocation } from 'react-router-dom'; // Import useLocation
import { vi } from 'vitest'; // Import vi for mocking

// Mock utility function
vi.mock('@/utils/generateStar', () => ({
  generateStars: vi.fn(),
}));

// Mock react-router-dom hooks used by ProductCard
vi.mock('react-router-dom', async () => {
  const originalModule = await vi.importActual('react-router-dom');
  return {
    ...originalModule,
    useLocation: vi.fn(), // Mock useLocation
  };
});

// Mock utility function used internally by ProductCard for delete
vi.mock('@/utils/getUserId', () => ({
    getCurrentUserId: vi.fn(() => 'currentUser123'), // Provide a mock user ID
}));


describe('ProductCard', () => {
  // Define mock props consistently
  const mockProduct = {
    productId: '1',
    productTitle: 'Sample Product',
    productPrice: 99.99,
    productCondition: 4,
    productImage: 'https://via.placeholder.com/150',
    productDescription: 'A sample product description.',
    userId: 'seller456', // --- ADDED: Crucial for onStartChat ---
  };
  let mockOnStartChat; // Use let to re-initialize in beforeEach

  beforeEach(() => {
    // Reset mocks before each test
    vi.clearAllMocks();
    generateStars.mockClear();

    // Re-initialize mock function for onStartChat
    mockOnStartChat = vi.fn();

    // Default mock for useLocation (can be overridden in specific tests)
    useLocation.mockReturnValue({ pathname: '/' });
  });

  // Helper function for rendering with router context
  const renderWithRouter = (ui, route = '/') => {
    // Set the specific pathname for this render
    useLocation.mockReturnValue({ pathname: route });
    return render(<MemoryRouter initialEntries={[route]}>{ui}</MemoryRouter>);
  };

  it('renders product title correctly', () => {
    renderWithRouter(<ProductCard product={mockProduct} onStartChat={mockOnStartChat} />);
    // Use queryAllByText and check the first one if multiple might exist due to hover state
    const titles = screen.queryAllByText(mockProduct.productTitle);
    expect(titles.length).toBeGreaterThan(0);
    expect(titles[0]).toBeInTheDocument();
  });

  it('renders the correct image with alt text', () => {
    renderWithRouter(<ProductCard product={mockProduct} onStartChat={mockOnStartChat} />);
    const img = screen.getByAltText(mockProduct.productTitle);
    expect(img).toBeInTheDocument();
    expect(img).toHaveAttribute('src', mockProduct.productImage);
  });

  it('renders the stars based on condition', () => {
    // Arrange: Set up the mock return value for generateStars
    const mockStarsElement = <div>⭐⭐⭐⭐</div>;
    generateStars.mockReturnValue(mockStarsElement);

    // Act
    renderWithRouter(<ProductCard product={mockProduct} onStartChat={mockOnStartChat} />);

    // Assert
    expect(generateStars).toHaveBeenCalledWith(mockProduct.productCondition);
    // Check if the content returned by the mock is present
    expect(screen.getByText('⭐⭐⭐⭐')).toBeInTheDocument();
  });

  it('renders the price correctly', () => {
    renderWithRouter(<ProductCard product={mockProduct} onStartChat={mockOnStartChat} />);
    expect(screen.getByText(`$${mockProduct.productPrice}`)).toBeInTheDocument();
  });

  // --- CORRECTED TEST ---
  it('renders Message button on /products and calls onStartChat with sellerId on click', () => {
    // Arrange: Render specifically for the /products route
    renderWithRouter(
        <ProductCard product={mockProduct} onStartChat={mockOnStartChat} />,
        '/products' // Set route to /products
    );

    // Act
    const messageButton = screen.getByRole('button', { name: /message/i }); // Use role for better accessibility
    expect(messageButton).toBeInTheDocument();
    fireEvent.click(messageButton);

    // Assert
    expect(mockOnStartChat).toHaveBeenCalledTimes(1);
    expect(mockOnStartChat).toHaveBeenCalledWith(mockProduct.userId); // Check if called with the correct seller ID
  });

  it('does NOT render Message button on /userproducts', () => {
    renderWithRouter(
        <ProductCard product={mockProduct} onStartChat={mockOnStartChat} />,
        '/userproducts' // Set route to /userproducts
    );
    expect(screen.queryByRole('button', { name: /message/i })).not.toBeInTheDocument();
  });

  it('renders Edit and Delete options only on /userproducts when menu is clicked', () => {
    renderWithRouter(
        <ProductCard product={mockProduct} onStartChat={mockOnStartChat} />,
        '/userproducts' // Set route to /userproducts
    );
    // Find the menu toggle button (assuming it's the only button initially visible besides potential save/cancel in edit mode)
    // Or add a specific test-id or aria-label to the FiMoreVertical button
    const menuButton = screen.getByRole('button'); // This might be fragile, consider adding test-id
    fireEvent.click(menuButton);

    // Assert that Edit and Delete are visible after clicking the menu
    expect(screen.getByText('Edit')).toBeInTheDocument();
    expect(screen.getByText('Delete')).toBeInTheDocument();
  });

  it('does NOT render Edit and Delete options on /products', () => {
    renderWithRouter(
        <ProductCard product={mockProduct} onStartChat={mockOnStartChat} />,
        '/products' // Set route to /products
    );
    // The menu button itself shouldn't even be there on /products
    expect(screen.queryByRole('button', { name: /edit/i })).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /delete/i })).not.toBeInTheDocument();
    // Check specifically for the menu toggle icon/button if possible
    // e.g., expect(screen.queryByLabelText('Product options menu')).not.toBeInTheDocument();
  });

  // Add tests for edit mode functionality (saving, canceling, deleting) if needed
});
