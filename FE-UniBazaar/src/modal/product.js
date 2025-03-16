class Product {
    constructor({ sellerID, productId, productTitle, productPrice, productCondition, productImage, productPostDate, productDescription }) {
      this.productSellerID = sellerID;
      this.productId = productId;
      this.productTitle = productTitle;
      this.productPrice = productPrice;
      this.productCondition = productCondition;
      this.productImage = productImage;
      this.productPostDate = productPostDate;
      this.productDescription = productDescription || "No description available.";  
    }
  
    getFormattedPrice() {
      return `$${this.productPrice.toFixed(2)}`;
    }
  }
  
  export default Product;
  