import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import ProductCard from "@/customComponents/ProductCard";
import { BrowserRouter as Router } from "react-router-dom";
import { vi } from "vitest";

const mockProduct = {
  productId: "1",
  productTitle: "Sample Product",
  productDescription: "This is a sample product",
  productPrice: 99,
  productCondition: 3,
  productImage: "https://via.placeholder.com/150",
};

const renderCard = (props = {}) =>
  render(
    <Router>
      <ProductCard product={mockProduct} {...props} />
    </Router>
  );

describe("ProductCard Component", () => {
  it("renders product title and price", () => {
    renderCard();

    expect(screen.getAllByText("Sample Product")[0]).toBeInTheDocument();
    expect(screen.getByText("$99")).toBeInTheDocument();
  });

  it("shows description on hover", () => {
    renderCard();

    const title = screen.getAllByText("Sample Product")[0];
    fireEvent.mouseEnter(title);

    expect(screen.getByText("This is a sample product")).toBeInTheDocument();
  });

  it("shows menu on /userproducts path", () => {
    const originalLocation = window.location;
    delete window.location;
    window.location = { pathname: "/userproducts" };

    renderCard();

    const moreButton = screen.getByRole("button");
    fireEvent.click(moreButton);

    expect(screen.getByText("Edit")).toBeInTheDocument();
    expect(screen.getByText("Delete")).toBeInTheDocument();

    window.location = originalLocation;
  });

  it("calls onEdit when Edit is clicked", () => {
    const onEditMock = vi.fn();
    const originalLocation = window.location;
    delete window.location;
    window.location = { pathname: "/userproducts" };

    renderCard({ onEdit: onEditMock });

    const moreButton = screen.getByRole("button");
    fireEvent.click(moreButton);

    const editOption = screen.getByText("Edit");
    fireEvent.click(editOption);

    expect(onEditMock).toHaveBeenCalledWith("1");

    window.location = originalLocation;
  });

  it("shows editable fields if propIsEditing is true and allows editing when Edit is clicked", () => {
    const onEditMock = vi.fn();

    const originalLocation = window.location;
    delete window.location;
    window.location = { pathname: "/userproducts" };

    renderCard({ propIsEditing: true, onEdit: onEditMock });

    const moreButton = screen.getByRole("button");
    fireEvent.click(moreButton);

    const editOption = screen.getByText("Edit");
    fireEvent.click(editOption);

    const titleInput = screen.getByPlaceholderText("Product Title");
    const descriptionInput = screen.getByPlaceholderText("Description");
    const priceInput = screen.getByPlaceholderText("Price");

    expect(titleInput).toBeInTheDocument();
    expect(descriptionInput).toBeInTheDocument();
    expect(priceInput).toBeInTheDocument();

    expect(titleInput.value).toBe("Sample Product");
    expect(descriptionInput.value).toBe("This is a sample product");
    expect(priceInput.value).toBe("99");

    expect(onEditMock).toHaveBeenCalledWith("1");

    window.location = originalLocation;
  });


  it("calls onCancel when cancel button is clicked in edit mode", () => {
    const onCancelMock = vi.fn();
  
    const originalLocation = window.location;
    delete window.location;
    window.location = { pathname: "/userproducts" };
  
    renderCard({
      propIsEditing: true,
      onCancel: onCancelMock,
    });
  
    const moreButton = screen.getByRole("button");
    fireEvent.click(moreButton);
  
    const editOption = screen.getByText("Edit");
    fireEvent.click(editOption);
  
    const titleInput = screen.getByPlaceholderText("Product Title");
    const descriptionInput = screen.getByPlaceholderText("Description");
    const priceInput = screen.getByPlaceholderText("Price");
  
    expect(titleInput).toBeInTheDocument();
    expect(descriptionInput).toBeInTheDocument();
    expect(priceInput).toBeInTheDocument();
  
    fireEvent.change(titleInput, { target: { value: "Updated Product Title" } });
    fireEvent.change(descriptionInput, { target: { value: "Updated product description." } });
    fireEvent.change(priceInput, { target: { value: "199" } });
  
    const cancelButton = screen.getByRole("button", { name: /cancel/i });
    fireEvent.click(cancelButton);
  
    expect(onCancelMock).toHaveBeenCalled();
  
    window.location = originalLocation;
  });
  

  it("calls onSave with updated values when Save is clicked", () => {
    const onSaveMock = vi.fn();
  
    const originalLocation = window.location;
    delete window.location;
    window.location = { pathname: "/userproducts" };
  
    renderCard({
      propIsEditing: true,
      onSave: onSaveMock,
    });
  
    const moreButton = screen.getByRole("button");
    fireEvent.click(moreButton);
  
    const editOption = screen.getByText("Edit");
    fireEvent.click(editOption);
  
    const titleInput = screen.getByPlaceholderText("Product Title");
    const descriptionInput = screen.getByPlaceholderText("Description");
    const priceInput = screen.getByPlaceholderText("Price");
  
    expect(titleInput).toBeInTheDocument();
    expect(descriptionInput).toBeInTheDocument();
    expect(priceInput).toBeInTheDocument();
  
    expect(titleInput.value).toBe("Sample Product");
    expect(descriptionInput.value).toBe("This is a sample product");
    expect(priceInput.value).toBe("99");
  
    fireEvent.change(titleInput, { target: { value: "Updated Product Title" } });
    fireEvent.change(descriptionInput, { target: { value: "Updated product description." } });
    fireEvent.change(priceInput, { target: { value: "199" } });
  
    const saveButton = screen.getByText("Save");
    fireEvent.click(saveButton);
  
    window.location = originalLocation;
  });
  


  it("shows 'Message' button on /products path", () => {
    const originalLocation = window.location;
    delete window.location;
    window.location = { pathname: "/products" };

    renderCard();

    expect(screen.getByText("Message")).toBeInTheDocument();

    window.location = originalLocation;
  });
});
