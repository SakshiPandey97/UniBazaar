import { useState } from "react";

export const useProductData = () => {
  const [productData, setProductData] = useState({
    productTitle: "",
    productDescription: "",
    productPrice: "",
    productCondition: "",
    productLocation: "",
    productImage: null, 
  });

  // Handle form field changes
  const handleChange = (e) => {
    setProductData({ ...productData, [e.target.name]: e.target.value });
  };

  // Handle file change and update the image in productData
  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    setProductData((prevData) => ({
      ...prevData,
      productImage: selectedFile, // Store the image file in productData
    }));
  };

  return {
    productData,
    handleChange,
    handleFileChange,
    setProductData,
  };
};
