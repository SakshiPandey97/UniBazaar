import React, { useState } from "react";
import CloudIcon from "../assets/imgs/cloud_icon.svg";
import fileIcon from "../assets/imgs/file_icon.svg";
import {
  PRODUCT_CONDITIONS,
  productConditionMapping,
} from "../utils/productMappings";
import { Button } from "@/ui/button";
import { motion } from "framer-motion";
import { useAnimation } from "../hooks/useAnimation";
import { postProductAPI } from "@/api/axios";
import ShowProductPreview from "@/customComponents/ShowProductPreview";
import { prepareFormData } from "@/utils/prepareFormData";
import { useProductData } from "@/hooks/useProductData";
import { useNavigate } from "react-router-dom";

const SellProductPage = () => {
  const navigate=useNavigate()

  const { productData, handleChange, handleFileChange, setProductData } =
    useProductData();
  const { isAnimating, triggerAnimation } = useAnimation();
  const [isUploaded, setIsUploaded] = useState(false); // Track successful upload

  const handleSubmit = async () => {
    if (!productData.productImage) {
      alert("Please upload a file before listing the product.");
      return;
    }

    const condition = productConditionMapping[productData.productCondition];

    if (!condition) {
      alert("Please select a valid product condition.");
      return;
    }

    try {
      await postProductAPI(prepareFormData(productData, productData.productImage, condition));
      triggerAnimation();
      setIsUploaded(true);
      setTimeout(()=>navigate("/"),3000)
    } catch (error) {
      console.error("Error posting product:", error);
      alert("Failed to post product. Try again.");
    }
  };

  return (
    <div className="flex min-h-screen bg-gray-100">
      {/* Left Side - Form Section */}
      <div className="w-1/2 p-8 bg-white shadow-md">
        <h2 className="text-2xl font-bold text-[#320B34] text-center mb-6">
          Product Information
        </h2>

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
            {PRODUCT_CONDITIONS.map((condition) => (
              <Button
                key={condition}
                onClick={() =>
                  setProductData({
                    ...productData,
                    productCondition: condition,
                  })
                }
                className={`px-4 py-2 rounded-lg border text-balck font-bold ${
                  productData.productCondition === condition
                    ? "bg-[#F58B00]"
                    : "bg-[#FFC67D]"
                }`}
              >
                {condition}
              </Button>
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

        {/* File Upload */}
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
          <Button
            onClick={() => document.getElementById("fileInput").click()}
            className="mt-4 hover:border-[#F58B00] border-2 p-2 bg-[#F58B00] hover:bg-[#FFC67D] text-balck font-bold px-6 py-2 rounded-lg"
          >
            Upload
          </Button>
          {productData.productImage && (
            <p className="text-green-600 mt-2">{productData.productImage.name}</p>
          )}
        </div>
      </div>
      <div className="w-1/2 flex flex-col items-center justify-center bg-gray-200 p-8 rounded-r-lg">
        {/* Cloud Icon */}
        <div className="flex h-3/5 justify-center mb-[-35px]">
          <img
            src={CloudIcon}
            className="rounded-lg shadow-lg w-2/3"
            alt="Cloud"
          />
        </div>
        <div className="flex h-1/5 justify-center">
          <motion.img
            src={fileIcon}
            className="rounded-lg w-10 left-1/2 -translate-x-1/2"
            alt="File"
            animate={{
              y: isAnimating ? [-100, 0] : [0, -100],
              opacity: isAnimating ? [1, 0] : [0, 1],
            }}
            transition={{
              repeat: Infinity,
              duration: 2,
              ease: "easeInOut",
              delay: 0.5,
            }}
          />
        </div>
        <Button
          onClick={handleSubmit}
          className="w-2/3 hover:border-[#F58B00] justify-center border-2 p-2 bg-[#F58B00] hover:bg-[#FFC67D] text-black font-bold py-3 rounded-lg text-lg"
        >
          List Now
        </Button>
      </div>
      {isUploaded && <ShowProductPreview productData={productData} />}
    </div>
  );
};

export default SellProductPage;
