import { useState, useEffect } from "react";
import { searchProductsAPI } from "../api/productAxios";  
import Product from "../modal/product";

const useSearchProducts = (query, limit = 10, debounceDelay = 750) => {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const [debouncedQuery, setDebouncedQuery] = useState(query);

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedQuery(query);
        }, debounceDelay);

        return () => clearTimeout(timer);
    }, [query, debounceDelay]);

    useEffect(() => {
        if (!debouncedQuery) {
            setProducts([]);
            return;
        }

        const fetchSearchResults = async () => {
            setLoading(true);
            setError(null);

            try {
                const data = await searchProductsAPI(debouncedQuery, limit);
                if (data.length === 0) {                
                    setProducts([]);
                } else {
                    const mappedProducts = data.map((item) => new Product(item));
                    setProducts(mappedProducts);
                }
            } catch (err) {
                setError("Failed to fetch search results.");
            } finally {
                setLoading(false);
            }
        };

        fetchSearchResults();
    }, [debouncedQuery, limit]);

    return { products, loading, error };
};

export default useSearchProducts;
