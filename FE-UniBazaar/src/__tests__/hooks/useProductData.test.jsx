import { renderHook, act } from '@testing-library/react';
import { vi } from 'vitest';
import { useProductData } from '../../hooks/useProductData';  // Adjust the path as necessary

describe('useProductData hook', () => {
  it('should return initial productData state', () => {
    const { result } = renderHook(() => useProductData());

    expect(result.current.productData).toEqual({
      productTitle: "",
      productDescription: "",
      productPrice: "",
      productCondition: "",
      productLocation: "",
      productImage: null,
    });
  });

  it('should update productData when handleChange is called', () => {
    const { result } = renderHook(() => useProductData());
    
    // Simulate change for productTitle
    act(() => {
      result.current.handleChange({ target: { name: 'productTitle', value: 'New Product' } });
    });

    expect(result.current.productData.productTitle).toBe('New Product');
  });

  it('should update productData when handleFileChange is called with a file', () => {
    const { result } = renderHook(() => useProductData());

    // Simulate file input change
    const file = new File(['dummy content'], 'example.jpg', { type: 'image/jpeg' });
    const event = { target: { files: [file] } };
    
    act(() => {
      result.current.handleFileChange(event);
    });

    expect(result.current.productData.productImage).toBe(file);
  });

  it('should allow setProductData to directly update productData', () => {
    const { result } = renderHook(() => useProductData());

    const newProductData = {
      productTitle: "Updated Product",
      productDescription: "This is a great product.",
      productPrice: "199.99",
      productCondition: "New",
      productLocation: "USA",
      productImage: null,
    };

    // Directly set productData
    act(() => {
      result.current.setProductData(newProductData);
    });

    expect(result.current.productData).toEqual(newProductData);
  });
});
