import React, { useState, useEffect, useRef } from "react";
import { motion } from "framer-motion";
import useFetchUserProducts from "../hooks/useFetchUserProducts";
import ProductCard from "@/customComponents/ProductCard";
import { getCurrentUserId } from "@/utils/getUserId";

function MyProductsPage() {
    const [lastId, setLastId] = useState("");
    const loadMoreButtonPositionRef = useRef(0);
    const [productsData, setProductsData] = useState([]);

    const limit = 12;
    const userId = getCurrentUserId();
    const { products: fetchedProducts, loading, error, hasMoreProducts } = useFetchUserProducts(userId, limit, lastId);

    const loadMoreButtonRef = useRef(null);

    useEffect(() => {
        if (fetchedProducts) {
            setProductsData((prevProducts) => {
                const newProducts = fetchedProducts.filter(
                    (newProduct) => !prevProducts.some((p) => p.productId === newProduct.productId)
                );
                return [...prevProducts, ...newProducts];
            });
        }
    }, [fetchedProducts]);

    const handleProductUpdated = (updatedProduct) => {
        setProductsData((prevProducts) =>
            prevProducts.map((product) =>
                product.productId === updatedProduct.productId
                    ? { ...product, ...updatedProduct }
                    : product
            )
        );
    };
    

    const handleProductDeleted = (deletedProductId) => {
        setProductsData((prevProducts) =>
            prevProducts.filter((product) => product.productId !== deletedProductId)
        );
    };


    const loadMoreProducts = () => {
        if (productsData.length > 0) {
            setLastId(productsData[productsData.length - 1].productId);
        }

        if (loadMoreButtonRef.current) {
            const rect = loadMoreButtonRef.current.getBoundingClientRect();
            loadMoreButtonPositionRef.current = rect.top + window.scrollY;
        }
    };

    useEffect(() => {
        if (loadMoreButtonPositionRef.current !== 0) {
            setTimeout(() => {
                window.scrollTo({
                    top: loadMoreButtonPositionRef.current,
                    behavior: "smooth",
                });
            }, 300);
        }
    }, [productsData]);

    if (error) {
        return (
            <div className="flex flex-col justify-center items-center h-screen text-center">
                <div className="text-xl text-red-600 font-semibold mb-4">
                    Oops, something went wrong!
                </div>
                <button
                    onClick={() => window.location.reload()}
                    className="bg-red-500 text-white px-6 py-2 rounded-lg font-semibold hover:bg-red-400"
                >
                    Try Again
                </button>
            </div>
        );
    }

    return (
        <div className="container mx-auto p-4">
            <h2 className="text-3xl font-semibold mb-6">My Products</h2>

            <motion.div
                className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6"
                initial="hidden"
                animate="visible"
                variants={{
                    hidden: { opacity: 0 },
                    visible: { opacity: 1, transition: { staggerChildren: 0.1, delayChildren: 0.1 } },
                }}
            >
                {loading && productsData.length === 0 ? (
                    <div className="flex justify-center col-span-full">
                        <div className="animate-spin rounded-full h-10 w-10 border-t-4 border-[#F58B00] border-solid"></div>
                    </div>
                ) : productsData.length > 0 ? (
                    productsData.map((product) => (
                        <motion.div
                            key={product.productId}
                            variants={{
                                hidden: { opacity: 0, scale: 0.9 },
                                visible: { opacity: 1, scale: 1, transition: { duration: 0.4, delay: 0.05 } },
                            }}
                        >
                            <ProductCard
                                product={product}
                                onProductUpdated={handleProductUpdated}
                                onProductDeleted={handleProductDeleted}

                            />
                        </motion.div>
                    ))
                ) : (
                    <div className="text-lg text-gray-500 col-span-full text-center">
                        No products found
                    </div>
                )}
            </motion.div>

            {productsData.length > 0 && (
                <div className="flex justify-center mt-8">
                    {hasMoreProducts ? (
                        <div className="flex flex-col items-center">
                            <button
                                ref={loadMoreButtonRef}
                                onClick={loadMoreProducts}
                                className="bg-[#F58B00] text-white px-6 py-2 rounded-lg font-semibold transition-all duration-300 cursor-pointer hover:bg-[#FFC67D] hover:shadow-lg focus:outline-none focus:ring-4 focus:ring-[#FFC67D] focus:ring-opacity-50"
                            >
                                Load More
                            </button>

                            {loading && (
                                <div className="mt-4">
                                    <div className="animate-spin rounded-full h-10 w-10 border-t-4 border-[#F58B00] border-solid"></div>
                                </div>
                            )}
                        </div>
                    ) : (
                        <div className="text-lg text-gray-500">No more products</div>
                    )}
                </div>
            )}
        </div>
    );
}

export default MyProductsPage;