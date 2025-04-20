import React from 'react';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom'; 
import Products from '@/customComponents/Products';
import useFetchProducts from '@/hooks/useFetchProducts';
import useStartChat from '@/hooks/useStartChat';
import { vi } from 'vitest';

vi.mock('@/hooks/useFetchProducts');
vi.mock('@/hooks/useStartChat');

describe('Products Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();

    useStartChat.mockReturnValue(vi.fn());
  });

  it('displays loading spinner while fetching products', () => {
    useFetchProducts.mockReturnValue({
      products: [],
      loading: true,
      error: null,
    });

    // Act: Render the component wrapped in MemoryRouter
    render(
      <MemoryRouter>
        <Products />
      </MemoryRouter>
    );

    // Assert: Check for the spinner (adjust selector/role if needed)
    expect(screen.getByRole('status')).toBeInTheDocument();
  });

  it('displays error message if there is an error', () => {
    // Arrange: Mock useFetchProducts to return an error
    const errorMessage = 'Network Error';
    useFetchProducts.mockReturnValue({
      products: [],
      loading: false,
      error: errorMessage, // Pass the error message string
    });

    // Act: Render the component wrapped in MemoryRouter
    render(
      <MemoryRouter>
        <Products />
      </MemoryRouter>
    );

    // Assert: Check for the error message
    expect(screen.getByText(errorMessage)).toBeInTheDocument();
  });

  it('displays products when loaded successfully', () => {
    const mockProducts = [
       { productId: '1', productTitle: 'Product 1', userId: 'seller1', /* other props */ },
       { productId: '2', productTitle: 'Product 2', userId: 'seller2', /* other props */ },
    ];
    useFetchProducts.mockReturnValue({
      products: mockProducts,
      loading: false,
      error: null,
    });

    render(
      <MemoryRouter>
        <Products />
      </MemoryRouter>
    );

    const product1Elements = screen.getAllByText('Product 1');
    expect(product1Elements.length).toBeGreaterThan(0); 
    expect(product1Elements[0]).toBeInTheDocument(); 

    const product2Elements = screen.getAllByText('Product 2');
    expect(product2Elements.length).toBeGreaterThan(0);
    expect(product2Elements[0]).toBeInTheDocument();


    expect(screen.queryByRole('status')).not.toBeInTheDocument();
  });
});