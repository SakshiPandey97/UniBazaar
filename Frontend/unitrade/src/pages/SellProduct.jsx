import React, { useState } from "react";
import { Switch } from "@headlessui/react"; // Ensure @headlessui/react is installed
import cartGuy from "../assets/imgs/cartGuy.svg"; // Ensure the path is correct

const SellProduct = () => {
  const [productData, setProductData] = useState({
    productTitle: "",
    productDescription: "",
    productPrice: "",
    productCondition: "",
    productLocation: "",
  });

  const [file, setFile] = useState(null);
  const [contactDetails, setContactDetails] = useState(true);
  const [allowChat, setAllowChat] = useState(false);
  const [allowOffer, setAllowOffer] = useState(false);

  const handleChange = (e) => {
    setProductData({ ...productData, [e.target.name]: e.target.value });
  };

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleSubmit = async () => {
    if (!file) {
      alert("Please upload a file before listing the product.");
      return;
    }
  
    const formData = new FormData();
    formData.append("userId", "1");
    formData.append("productTitle", productData.productTitle);
    formData.append("productDescription", productData.productDescription);
    formData.append("productPrice", productData.productPrice);
    formData.append("productCondition", productData.productCondition);
    formData.append("productLocation", productData.productLocation);
    formData.append("productPostDate", new Date().toISOString().split("T")[0]);
    formData.append("productImage", file);
    formData.append("productImageType", file.type); // Include image type
  
    try {
      await fetch("http://192.168.0.203:8080/products", {
        method: "POST",
        body: formData,
      });
      alert("Product posted successfully!");
    } catch (error) {
      console.error("Error posting product:", error);
      alert("Failed to post product. Try again.");
    }
  };
  
  

  return (
    <div className="flex min-h-screen bg-gray-100">
      {/* Left Side - Form Section */}
      <div className="w-1/2 p-8 bg-white shadow-md">
        <h2 className="text-2xl font-bold text-teal-600 text-center mb-6">
          Product Information
        </h2>

        {/* Product Form Fields */}
        <input
          type="text"
          name="productTitle"
          placeholder="Product Title"
          value={productData.productTitle}
          onChange={handleChange}
          className="border p-3 w-full rounded-lg mb-3"
        />
        <textarea
          name="productDescription"
          placeholder="Description"
          value={productData.productDescription}
          onChange={handleChange}
          className="border p-3 w-full rounded-lg mb-3"
        />
        <input
          type="number"
          name="productPrice"
          placeholder="Price ($)"
          value={productData.productPrice}
          onChange={handleChange}
          className="border p-3 w-full rounded-lg mb-3"
        />
        <div className="mb-4">
          <h3 className="text-lg font-semibold mb-2">Product Condition</h3>
          <div className="flex gap-2">
            {["Excellent", "Very Good", "Good", "Fair", "Poor"].map((condition) => (
              <button
                key={condition}
                onClick={() =>
                  setProductData({ ...productData, productCondition: condition })
                }
                className={`px-4 py-2 rounded-lg border ${
                  productData.productCondition === condition
                    ? "bg-green-500 text-white"
                    : "bg-gray-200 text-gray-700"
                }`}
              >
                {condition}
              </button>
            ))}
          </div>
        </div>
        <input
          type="text"
          name="productLocation"
          placeholder="Location"
          value={productData.productLocation}
          onChange={handleChange}
          className="border p-3 w-full rounded-lg mb-3"
        />

        {/* File Upload - Replaced with Upload Button */}
        <div className="border-2 border-dashed border-gray-400 p-6 flex flex-col items-center">
          <input
            type="file"
            accept="image/*"
            onChange={handleFileChange}
            className="hidden"
            id="fileInput"
          />
          <label
            htmlFor="fileInput"
            className="cursor-pointer flex flex-col items-center"
          >
            <div className="flex items-center space-x-2">
              <span role="img" aria-label="camera">
                📷
              </span>
              <span role="img" aria-label="video">
                📹
              </span>
            </div>
            <p className="text-gray-500 mt-2">Drag and drop files</p>
          </label>
          <button
            onClick={() => document.getElementById("fileInput").click()}
            className="mt-4 bg-green-600 text-white px-6 py-2 rounded-lg"
          >
            Upload
          </button>
          {file && <p className="text-green-600 mt-2">{file.name}</p>}
        </div>
      </div>

      {/* Right Side - Settings & Illustration */}
      <div className="w-1/2 flex flex-col items-center justify-center bg-gray-200 p-8 rounded-r-lg">
        {/* Toggle Options */}
        <div className="w-full max-w-md">
          <div className="flex justify-between items-center mb-4">
            <span className="text-lg">Share Contact Details</span>
            <Switch
              checked={contactDetails}
              onChange={setContactDetails}
              className={`relative inline-flex items-center h-6 rounded-full w-12 transition ${
                contactDetails ? "bg-green-500" : "bg-gray-300"
              }`}
            >
              <span
                className={`inline-block w-6 h-6 bg-white rounded-full transform transition ${
                  contactDetails ? "translate-x-6" : "translate-x-0"
                }`}
              />
            </Switch>
          </div>

          <div className="flex justify-between items-center mb-4">
            <span className="text-lg">Allow to Chat</span>
            <Switch
              checked={allowChat}
              onChange={setAllowChat}
              className={`relative inline-flex items-center h-6 rounded-full w-12 transition ${
                allowChat ? "bg-green-500" : "bg-gray-300"
              }`}
            >
              <span
                className={`inline-block w-6 h-6 bg-white rounded-full transform transition ${
                  allowChat ? "translate-x-6" : "translate-x-0"
                }`}
              />
            </Switch>
          </div>

          <div className="flex justify-between items-center mb-6">
            <span className="text-lg">Allow to Make an Offer</span>
            <Switch
              checked={allowOffer}
              onChange={setAllowOffer}
              className={`relative inline-flex items-center h-6 rounded-full w-12 transition ${
                allowOffer ? "bg-green-500" : "bg-gray-300"
              }`}
            >
              <span
                className={`inline-block w-6 h-6 bg-white rounded-full transform transition ${
                  allowOffer ? "translate-x-6" : "translate-x-0"
                }`}
              />
            </Switch>
          </div>

          {/* List Now Button with Submit Functionality */}
          <button
            onClick={handleSubmit}
            className="w-full bg-green-600 text-white py-3 rounded-lg text-lg font-semibold hover:bg-green-700"
          >
            List Now
          </button>
        </div>

        {/* Illustration */}
        <div className="mt-8">
          <img src={cartGuy} className="rounded-lg shadow-lg" alt="Illustration" />
        </div>
      </div>
    </div>
  );
};

export default SellProduct;
