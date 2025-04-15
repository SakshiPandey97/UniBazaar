import { useState } from "react";
import { toast } from "react-toastify";
import { postProductAPI } from "@/api/productAxios";
import { prepareFormData } from "@/utils/prepareFormData";
import { productConditionMapping } from "@/utils/productMappings";

export const useCreateProduct = () => {
  const [isLoading, setIsLoading] = useState(false);

  const createProduct = async (productData, onSuccess) => {
    if (!productData.productImage) {
      toast.error("Please upload a file before listing the product.");
      return;
    }

    if (!productData.productTitle || productData.productTitle.trim() === "") {
      toast.error("Product title cannot be empty.");
      return;
    }

    if (
      !productData.productPrice ||
      isNaN(productData.productPrice) ||
      Number(productData.productPrice) <= 0
    ) {
      toast.error("Please enter a valid price greater than 0.");
      return;
    }

    const condition = productConditionMapping[productData.productCondition];

    if (!condition) {
      toast.error("Please select a valid product condition.");
      return;
    }

    try {
      setIsLoading(true);
      await postProductAPI(
        prepareFormData(productData, productData.productImage, condition, true)
      );
      onSuccess?.();
    } catch (error) {
      console.error("Error posting product:", error);
      toast.error("Failed to post product. Try again.");
    } finally {
      setIsLoading(false);
    }
  };

  return {
    isLoading,
    createProduct,
  };
};
