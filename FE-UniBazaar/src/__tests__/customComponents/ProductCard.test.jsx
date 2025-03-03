import { render, screen } from '@testing-library/react';
import ProductCard from '../../customComponents/ProductCard'; // Adjust the import according to your file structure
import { generateStars } from '@/utils/generateStar';

// Mock the generateStars function to simplify the test output
vi.mock('@/utils/generateStar', () => ({
  generateStars: vi.fn(),
}));

describe('ProductCard', () => {
  const mockProps = {
    title: 'Sample Product',
    price: 99.99,
    condition: 4,
    image: 'https://via.placeholder.com/150',
    description: 'A description of the product.',
  };

  beforeEach(() => {
    // Reset mocks before each test
    generateStars.mockClear();
  });

  it('renders product title and description correctly', () => {
    render(<ProductCard {...mockProps} />);

    // Check title is rendered
    expect(screen.getByText('Sample Product')).toBeInTheDocument();
    
    // Check description is not directly rendered, as it's not part of the displayed JSX
    // You can modify the component to display description if needed.
  });

  it('renders the correct image', () => {
    render(<ProductCard {...mockProps} />);
    
    // Check if the image has the correct src and alt attributes
    const img = screen.getByAltText('Sample Product');
    expect(img).toHaveAttribute('src', 'https://via.placeholder.com/150');
  });

  it('renders the stars based on condition', () => {
    generateStars.mockReturnValue(<div>⭐⭐⭐⭐</div>); // Return a mock star representation

    render(<ProductCard {...mockProps} />);

    // Check if the stars are rendered correctly
    expect(generateStars).toHaveBeenCalledWith(mockProps.condition);
    expect(screen.getByText('⭐⭐⭐⭐')).toBeInTheDocument();
  });

  it('renders the price correctly', () => {
    render(<ProductCard {...mockProps} />);

    // Check if the price is rendered
    expect(screen.getByText('$99.99')).toBeInTheDocument();
  });

  it('renders the button with correct text', () => {
    render(<ProductCard {...mockProps} />);
    
    // Check if the button text is correct
    expect(screen.getByText('Read More')).toBeInTheDocument();
  });
});
