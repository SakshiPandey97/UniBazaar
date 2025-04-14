import { getCurrentUserId } from "./getUserId";

export const prepareFormData = (productData, file, condition, isPostRequest = true) => {
  const formData = new FormData();
  const options = { month: '2-digit', day: '2-digit', year: 'numeric' };
  const productPostDate = new Date().toLocaleDateString("en-US", options).replace(/\//g, "-");

  const userId = getCurrentUserId();

  formData.append("userId", userId);
  formData.append("productTitle", productData.productTitle);
  formData.append("productDescription", productData.productDescription);
  formData.append("productPrice", productData.productPrice);
  formData.append("productCondition", condition);
  formData.append("productLocation", productData.productLocation);
  formData.append("productPostDate", productPostDate);

  if (file) {
    formData.append("productImage", file);
  }

  if (!isPostRequest && productData.productId) {
    formData.append("productId", productData.productId);
  }

  return formData;
};
