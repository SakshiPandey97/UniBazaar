import { useState } from "react";
import { updateProductAPI } from "../api/productAxios";
import { prepareFormData } from "@/utils/prepareFormData";
import { getCurrentUserId } from "@/utils/getUserId";

const useUpdateProduct = () => {
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState(null);

  const updateProduct = async (editableProduct, newImageFile, condition) => {
    const userId = getCurrentUserId();
    const formData = prepareFormData(editableProduct, newImageFile, condition, false);

    setIsSaving(true);
    setError(null);

    try {
      const response = await updateProductAPI(userId, editableProduct.productId, formData);
      setIsSaving(false);

      return response;
    } catch (error) {
      console.error(error);
      setError("Failed to update product.");
      return null;
    } finally {
      setIsSaving(false);
    }
  };

  return { updateProduct, isSaving, error };
};

export default useUpdateProduct;