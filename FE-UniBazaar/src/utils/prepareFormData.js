import { getCurrentUserId } from "./getUserId";
export const prepareFormData = (productData, file, condition) => {
    const formData = new FormData();
    const options = { month: '2-digit', day: '2-digit', year: 'numeric' };
    const productPostDate = new Date().toLocaleDateString("en-US", options).replace(/\//g, "-");

    formData.append("userId", {getCurrentUserId});
    formData.append("productTitle", productData.productTitle);
    formData.append("productDescription", productData.productDescription);
    formData.append("productPrice", productData.productPrice);
    formData.append("productCondition", condition);
    formData.append("productLocation", productData.productLocation);
    formData.append("productPostDate", productPostDate);
    formData.append("productImage", file);
  
    return formData;
  };