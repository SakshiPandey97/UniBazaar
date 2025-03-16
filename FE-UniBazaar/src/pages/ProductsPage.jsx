import React, { useState, useEffect, useRef } from "react";
import useAllProducts from "../hooks/useAllProducts";
import ProductCard from "@/customComponents/ProductCard";

function ProductsPage() {
  const [lastId, setLastId] = useState("");
  const limit = 12;

  const { products, loading, error, hasMoreProducts } = useAllProducts(
    limit,
    lastId
  );
  const loadMoreButtonRef = useRef(null);
  const loadMoreButtonPositionRef = useRef(0);
  const loadMoreProducts = () => {
    if (products.length > 0) {
      setLastId(products[products.length - 1].productId);
    }

    if (loadMoreButtonRef.current) {
      const rect = loadMoreButtonRef.current.getBoundingClientRect();
      loadMoreButtonPositionRef.current = rect.top + window.scrollY;
    }
  };

  useEffect(() => {
    if (loadMoreButtonPositionRef.current !== 0) {
      window.scrollTo({
        top: loadMoreButtonPositionRef.current,
        behavior: "smooth",
      });
    }
  }, [products]);

  if (loading) {
    return (
      <div className="flex justify-center items-center h-screen">
        <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-[#F58B00] border-solid"></div>
      </div>
    );
  }

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
      <h2 className="text-3xl font-semibold mb-6">Products</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {products.map((product) => (
          <ProductCard
            key={product.productId}
            product={product}
          />
        ))}
      </div>
      <div className="flex justify-center mt-8">
        {!hasMoreProducts ? (
          <div className="text-lg text-gray-500">No more products</div>
        ) : (
          <button
            ref={loadMoreButtonRef}
            onClick={loadMoreProducts}
            className="bg-[#F58B00] text-white px-6 py-2 rounded-lg font-semibold hover:bg-[#FFC67D]"
          >
            Load More
          </button>
        )}
      </div>

    </div>
  );
}

export default ProductsPage;
