import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { MemoryRouter } from 'react-router-dom'; 
import ProductCard from '@/customComponents/ProductCard';
import { useUserAuth } from '@/hooks/useUserAuth'; 

vi.mock('@/hooks/useUserAuth');
vi.mock('@/utils/generateStar', () => ({ generateStars: vi.fn(() => <div>Stars</div>) })); 

describe('ProductCard Simple Render', () => {
  const mockProduct = {
    productId: 'simple123',
    productTitle: 'Simple Test Product',
    productDescription: 'Description',
    productPrice: 10.00,
    productCondition: 3, 
    productImage: 'simple-image.jpg',
    userId: 'sellerSimple',
  };

  it('should render the product title', () => {
    useUserAuth.mockReturnValue({ userState: false });

    render(
      <MemoryRouter> 
        <ProductCard product={mockProduct} />
      </MemoryRouter>
    );

    const titleElements = screen.queryAllByText(mockProduct.productTitle);
    expect(titleElements.length).toBeGreaterThan(0);

  });
});
