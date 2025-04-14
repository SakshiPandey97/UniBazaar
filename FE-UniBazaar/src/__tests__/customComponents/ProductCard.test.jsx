import { render, screen, fireEvent } from '@testing-library/react';
import ProductCard from '../../customComponents/ProductCard';
import { generateStars } from '@/utils/generateStar';
import { MemoryRouter } from 'react-router-dom';

vi.mock('@/utils/generateStar', () => ({
  generateStars: vi.fn(),
}));

describe('ProductCard', () => {
  const mockProps = {
    product: {
      productId: '1',
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

  const renderWithRouter = (ui, route = '/products') => {
    return render(<MemoryRouter initialEntries={[route]}>{ui}</MemoryRouter>);
  };

  it('renders product title correctly', () => {
    renderWithRouter(<ProductCard {...mockProps} />);
    const [visibleTitle] = screen.getAllByText(mockProps.product.productTitle);
    expect(visibleTitle).toBeInTheDocument();
  });

  it('renders the correct image with alt text', () => {
    renderWithRouter(<ProductCard {...mockProps} />);
    const img = screen.getByAltText(mockProps.product.productTitle);
    expect(img).toHaveAttribute('src', mockProps.product.productImage);
  });

  it('renders the stars based on condition', () => {
    generateStars.mockReturnValue(<div>⭐⭐⭐⭐</div>);
    renderWithRouter(<ProductCard {...mockProps} />);
    expect(generateStars).toHaveBeenCalledWith(mockProps.product.productCondition);
    expect(screen.getByText('⭐⭐⭐⭐')).toBeInTheDocument();
  });

  it('renders the price correctly', () => {
    renderWithRouter(<ProductCard {...mockProps} />);
    expect(screen.getByText(`$${mockProps.product.productPrice}`)).toBeInTheDocument();
  });

  it('renders Message button on /products and handles click event', () => {
    renderWithRouter(<ProductCard {...mockProps} />, '/products');
    const messageButton = screen.getByText('Message');
    expect(messageButton).toBeInTheDocument();
    fireEvent.click(messageButton);
    expect(mockProps.onClick).toHaveBeenCalled();
  });

  it('does NOT render Message button on /userproducts', () => {
    renderWithRouter(<ProductCard {...mockProps} />, '/userproducts');
    expect(screen.queryByText('Message')).not.toBeInTheDocument();
  });

  it('renders Edit and Delete options only on /userproducts', () => {
    renderWithRouter(<ProductCard {...mockProps} />, '/userproducts');
    const menuButton = screen.getByRole('button');
    fireEvent.click(menuButton);
    expect(screen.getByText('Edit')).toBeInTheDocument();
    expect(screen.getByText('Delete')).toBeInTheDocument();
  });

  it('does NOT render Edit and Delete options on /products', () => {
    renderWithRouter(<ProductCard {...mockProps} />, '/products');
    const menuButton = screen.queryByRole('button');
    if (menuButton) {
      fireEvent.click(menuButton);
    }
    expect(screen.queryByText('Edit')).not.toBeInTheDocument();
    expect(screen.queryByText('Delete')).not.toBeInTheDocument();
  });
});
