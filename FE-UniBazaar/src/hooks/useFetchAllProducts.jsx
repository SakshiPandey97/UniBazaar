import { useState, useEffect, useRef } from "react";
import { getAllProductsAPI } from "../api/productAxios";
import Product from "../modal/product";

function useFetchAllProducts(limit, lastId) {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [hasMoreProducts, setHasMoreProducts] = useState(true);


    const fetchedProductIdsRef = useRef(new Set());

    useEffect(() => {
        const fetchProducts = async () => {
            setLoading(true);
            console.log("Fetching Paginated products...");

            try {
                const data = await getAllProductsAPI(limit, lastId);

                console.log("Fetched products length:", data.length);

                if (data.length === 0) {
                    setHasMoreProducts(false);
                    return;
                }

                if (data.length < limit) {
                    setHasMoreProducts(false);
                }

                const mappedProducts = data.map((item) => new Product(item));

                const newProducts = mappedProducts.filter((product) => {
                    if (!fetchedProductIdsRef.current.has(product.productId)) {
                        fetchedProductIdsRef.current.add(product.productId);
                        return true;
                    }
                    return false;
                });

                setProducts((prevProducts) => [...prevProducts, ...newProducts]);
            } catch (err) {
                if (err.response && err.response.status === 404) {
                    setHasMoreProducts(false);
                }
                else {
                    setError("Failed to fetch products.");
                }
            } finally {
                setLoading(false);
            }
        };

        fetchProducts();
    }, [limit, lastId]);

    return { products, loading, error, hasMoreProducts };
}

export default useFetchAllProducts;
