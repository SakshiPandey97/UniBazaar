import { useState, useEffect } from "react";
import { getAllProductsAPI } from "../api/axios";  
import Product from "../modal/product";

const useProducts = () => {
  const [products, setProducts] = useState([]);  
  const [loading, setLoading] = useState(true);   
  const [error, setError] = useState(null);       

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const lastId = "";
        const limit = 10;
  
        const data = await getAllProductsAPI(limit, lastId); 
        
        const mappedProducts = data.map((item) => new Product(item));
        const sortedProducts = mappedProducts.sort(
          (a, b) => new Date(b.productPostDate) - new Date(a.productPostDate)
        );
        
        setProducts(sortedProducts); 
      } catch (err) {
        setError("Error fetching products");
        console.error("Error fetching products:", err);
      } finally {
        setLoading(false); 
      }
    };
  
    fetchProducts();
  }, []);

  return { products, loading, error };  
};

export default useProducts;
