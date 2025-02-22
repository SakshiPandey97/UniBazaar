class Product {
    constructor({ productId, productTitle, productPrice, productImage, productPostDate, productDescription }) {
      this.productId = productId;
      this.productTitle = productTitle;
      this.productPrice = productPrice;
      this.productImage = productImage;
      this.productPostDate = productPostDate;
      this.productDescription = productDescription || "No description available.";  
    }
  
    getFormattedPrice() {
      return `$${this.productPrice.toFixed(2)}`;
    }
  }
  
  export default Product;
  