import { render, screen, fireEvent } from '@testing-library/react';
import ProductCard from '../../customComponents/ProductCard'; // Adjust import path
import { generateStars } from '@/utils/generateStar';

// Mock the generateStars function
vi.mock('@/utils/generateStar', () => ({
  generateStars: vi.fn(),
}));

describe('ProductCard', () => {
  const mockProps = {
    product: {
      productTitle: 'Sample Product',
      productPrice: 99.99,
      productCondition: 4,
      productImage: 'https://via.placeholder.com/150',
      productDescription: 'A sample product description.',
    },
    onClick: vi.fn(),
  };

  beforeEach(() => {
    generateStars.mockClear();
  });

  it('renders product title correctly', () => {
    render(<ProductCard {...mockProps} />);

    // Ensure the product title is displayed in the non-hover state
    const [visibleTitle] = screen.getAllByText(mockProps.product.productTitle);
    expect(visibleTitle).toBeInTheDocument();
  });

  it('renders the correct image with alt text', () => {
    render(<ProductCard {...mockProps} />);
    
    // Check if the image has the correct src and alt attributes
    const img = screen.getByAltText(mockProps.product.productTitle);
    expect(img).toHaveAttribute('src', mockProps.product.productImage);
  });

  it('renders the stars based on condition', () => {
    generateStars.mockReturnValue(<div>⭐⭐⭐⭐</div>); // Mock star representation

    render(<ProductCard {...mockProps} />);

    // Ensure the function is called with the correct condition
    expect(generateStars).toHaveBeenCalledWith(mockProps.product.productCondition);
    expect(screen.getByText('⭐⭐⭐⭐')).toBeInTheDocument();
  });

  it('renders the price correctly', () => {
    render(<ProductCard {...mockProps} />);

    // Check if the price is displayed correctly
    expect(screen.getByText(`$${mockProps.product.productPrice}`)).toBeInTheDocument();
  });

  it('renders the Message button and handles click event', () => {
    render(<ProductCard {...mockProps} />);

    // Check if the button is present
    const messageButton = screen.getByText('Message');
    expect(messageButton).toBeInTheDocument();

    // Simulate click event
    fireEvent.click(messageButton);
    expect(mockProps.onClick).toHaveBeenCalled();
  });
});
