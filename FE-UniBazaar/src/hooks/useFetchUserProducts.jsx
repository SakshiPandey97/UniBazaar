import { useState, useEffect, useRef } from "react";
import { getUserProductsAPI } from "../api/productAxios";
import Product from "../modal/product";

function useFetchUserProducts(userId, limit, lastId) {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [hasMoreProducts, setHasMoreProducts] = useState(true);

    const fetchedProductIdsRef = useRef(new Set());

    useEffect(() => {
        if (!userId) return;
        console.log(userId)

        const fetchProducts = async () => {
            setLoading(true);
            console.log(`Fetching products for user ${userId}...`);

            try {
                const data = await getUserProductsAPI(userId, limit, lastId);

                console.log("Fetched user products length:", data.length);

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

                newProducts.sort((a, b) => {
                    const dateA = new Date(a.productPostDate);
                    const dateB = new Date(b.productPostDate);
                    return dateB - dateA;
                });

                setProducts((prevProducts) => [...prevProducts, ...newProducts]);
            } catch (err) {
                if (err.response && err.response.status === 404) {
                    setHasMoreProducts(false);
                } else {
                    setError("Failed to fetch user products.");
                }
            } finally {
                setLoading(false);
            }
        };

        fetchProducts();
    }, [userId, limit, lastId]);

    return { products, loading, error, hasMoreProducts };
}

export default useFetchUserProducts;
