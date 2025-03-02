export const prepareFormData = (productData, file, condition) => {
    const formData = new FormData();
    const productPostDate = new Date().toLocaleDateString("en-GB").replace(/\//g, "-");
  
    formData.append("userId", 1);
    formData.append("productTitle", productData.productTitle);
    formData.append("productDescription", productData.productDescription);
    formData.append("productPrice", productData.productPrice);
    formData.append("productCondition", condition);
    formData.append("productLocation", productData.productLocation);
    formData.append("productPostDate", productPostDate);
    formData.append("productImage", file);
  
    return formData;
  };