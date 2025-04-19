import React, { useState } from "react";
import {
  PRODUCT_CONDITIONS,
} from "../utils/productMappings";
import { Button } from "@/ui/button";
import { useProductData } from "@/hooks/useProductData";
import { useNavigate } from "react-router-dom";
import { motion } from "framer-motion";
import { useAnimation } from "../hooks/useAnimation";
import { useCreateProduct } from "../hooks/useCreateProduct";
import { FiUploadCloud, FiCheckCircle } from "react-icons/fi";
import {
  productConditionMapping,
  conditionColorMap,
} from "../utils/productMappings";

const SellProductPage = () => {
  const navigate = useNavigate();
  const { productData, handleChange, setProductData } = useProductData();
  const { isAnimating, triggerAnimation, animationControls } = useAnimation();
  const [isUploaded, setIsUploaded] = useState(false);
  const [imagePreview, setImagePreview] = useState(null);
  const { isLoading, createProduct } = useCreateProduct();


  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      setProductData({ ...productData, productImage: file });
      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result);
      };
      reader.readAsDataURL(file);
    } else {
      setImagePreview(null);
      setProductData({ ...productData, productImage: null });
    }
  };

  const handleSubmit = () => {
    createProduct(productData, () => {
      triggerAnimation();
      setIsUploaded(true);
      setTimeout(() => navigate("/userproducts"), 3000);
    });
  };

  return (
    <div className="bg-gray-100 min-h-screen py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto bg-white shadow-xl rounded-lg p-8">
        <h2 className="text-3xl font-extrabold text-[#320B34] text-center mb-8">
          List Your Item
        </h2>

        <div className="space-y-6">
          {/* Product Title */}
          <div>
            <label htmlFor="productTitle" className="block text-sm font-medium text-gray-700">
              Product Title
            </label>
            <input
              type="text"
              id="productTitle"
              name="productTitle"
              placeholder="e.g., Used Textbook: Intro to Calculus"
              value={productData.productTitle}
              onChange={handleChange}
              className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            />
          </div>

          {/* Description */}
          <div>
            <label htmlFor="productDescription" className="block text-sm font-medium text-gray-700">
              Description
            </label>
            <textarea
              id="productDescription"
              name="productDescription"
              rows="3"
              placeholder="Describe your product in detail..."
              value={productData.productDescription}
              onChange={handleChange}
              maxLength={250}
              className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            />
            <p className="text-sm text-gray-500 text-right">
              {productData.productDescription.length}/250 characters
            </p>
          </div>


          {/* Price */}
          <div>
            <label htmlFor="productPrice" className="block text-sm font-medium text-gray-700">
              Price ($)
            </label>
            <div className="relative rounded-md shadow-sm">
              <div className="pointer-events-none absolute inset-y-0 left-0 pl-3 flex items-center">
                <span className="text-gray-500 sm:text-sm">$</span>
              </div>
              <input
                type="number"
                id="productPrice"
                name="productPrice"
                placeholder="e.g., 25.00"
                value={productData.productPrice}
                onChange={handleChange}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 pl-7 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>
          </div>

          {/* Product Condition */}
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Product Condition
            </label>
            <div className="mt-2 grid grid-cols-2 md:grid-cols-3 gap-3">
              {PRODUCT_CONDITIONS.map((condition) => {
                const isSelected = productData.productCondition === condition;
                const conditionValue = productConditionMapping[condition];
                const selectedColorClass = isSelected ? conditionColorMap[conditionValue] : "bg-white hover:bg-gray-100 text-gray-700";

                return (
                  <Button
                    key={condition}
                    onClick={() =>
                      setProductData({ ...productData, productCondition: condition })
                    }
                    className={`px-4 py-2 rounded-md border font-medium transition-all shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 ${selectedColorClass}`}
                  >
                    {condition}
                  </Button>
                );
              })}
            </div>
          </div>

          {/* Location */}
          <div>
            <label htmlFor="productLocation" className="block text-sm font-medium text-gray-700">
              Location
            </label>
            <input
              type="text"
              id="productLocation"
              name="productLocation"
              placeholder="e.g., Gainesville, FL"
              value={productData.productLocation}
              onChange={handleChange}
              className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            />
          </div>

          {/* File Upload with Preview */}
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Upload Product Image
            </label>
            <div className="mt-1 flex flex-col items-center justify-center px-6 pt-5 pb-6 border-2 border-dashed border-gray-300 rounded-md">
              {imagePreview ? (
                <div className="mb-4">
                  <div className="w-72 h-60 rounded-md overflow-hidden mx-auto">
                    <img
                      src={imagePreview}
                      alt="Image Preview"
                      className="object-cover w-full h-full"
                    />
                  </div>
                  <div className="mt-2 text-center">
                    <button
                      type="button"
                      onClick={() => {
                        setImagePreview(null);
                        setProductData({ ...productData, productImage: null });
                        const fileInput = document.getElementById('file-upload');
                        if (fileInput) {
                          fileInput.value = '';
                        }
                      }}
                      className="mt-2 inline-flex items-center px-4 py-2 bg-red-500 text-white text-sm font-medium rounded-md shadow-sm hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500"
                    >
                      Remove Image
                    </button>
                  </div>
                </div>
              ) : (
                <div className="space-y-1 text-center">
                  <FiUploadCloud className="mx-auto h-12 w-12 text-gray-400" />
                  <div className="flex text-sm text-gray-600 justify-center">
                    <label
                      htmlFor="file-upload"
                      className="relative cursor-pointer bg-white rounded-md font-medium text-indigo-600 hover:text-indigo-500 focus-within:outline-none focus-within:ring-2 focus-within:ring-offset-2 focus-within:ring-indigo-500"
                    >
                      <span>Upload a file</span>
                      <input
                        id="file-upload"
                        name="productImage"
                        type="file"
                        accept="image/*"
                        className="sr-only"
                        onChange={handleFileChange}
                      />
                    </label>
                    <p className="pl-1">or drag and drop</p>
                  </div>
                  <p className="text-xs text-gray-500">PNG, JPG, GIF up to 10MB</p>
                </div>
              )}


              {productData.productImage && !imagePreview && (
                <p className="text-green-600 mt-2">{productData.productImage.name}</p>
              )}
            </div>
          </div>

          {/* Submit Button */}
          <div>
            <Button
              onClick={handleSubmit}
              disabled={isLoading}
              className={`w-full ${isLoading ? "bg-gray-400 cursor-not-allowed" : "bg-[#F58B00] hover:cursor-pointer hover:bg-[#FFC67D]"} text-white font-bold py-3 rounded-md shadow-md transition-all text-lg focus:outline-none focus:ring-indigo-500`}
            >
              {isLoading ? "Listing..." : "List Now"}
            </Button>
          </div>

          {/* Success Message */}
          {isUploaded && (
            <motion.div
              initial={{ opacity: 0, y: -20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="mt-4 p-4 bg-gray-100 text-gray-700 rounded-md flex items-center"
            >
              <FiCheckCircle className="mr-2 h-5 w-5 text-gray-500" />
              Redirecting you to your listings...
            </motion.div>
          )}
        </div>
      </div>
    </div>
  );
};

export default SellProductPage;