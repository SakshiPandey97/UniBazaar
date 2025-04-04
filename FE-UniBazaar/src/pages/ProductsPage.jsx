import React, { useState, useEffect, useRef } from "react";
import { motion } from "framer-motion";
import { FaSearch } from "react-icons/fa";
import { useSearchContext } from "../context/SearchContext";
import useFetchAllProducts from "../hooks/useFetchAllProducts";
import useSearchProducts from "../hooks/useSearchProducts";
import ProductCard from "@/customComponents/ProductCard";

function ProductsPage() {
  const { searchTerm: globalSearchTerm, setSearchTerm: setGlobalSearchTerm } = useSearchContext();
  const [lastId, setLastId] = useState("");
  const loadMoreButtonPositionRef = useRef(0);
  const [isSearchCleared, setIsSearchCleared] = useState(false);

  const limit = 12;
  const searchLimit = 100;

  const { products: searchProducts, loading: searchLoading, error: searchError } = useSearchProducts(globalSearchTerm, searchLimit);
  const { products: allProducts, loading: allLoading, error: allError, hasMoreProducts } = useFetchAllProducts(limit, lastId);

  const products = globalSearchTerm ? searchProducts : allProducts;
  const loading = globalSearchTerm ? searchLoading : allLoading;
  const error = globalSearchTerm ? searchError : allError;

  const loadMoreButtonRef = useRef(null);
  const searchInputRef = useRef(null);

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
    if (loadMoreButtonPositionRef.current !== 0 && !isSearchCleared) {
      setTimeout(() => {
        window.scrollTo({
          top: loadMoreButtonPositionRef.current,
          behavior: "smooth",
        });
      }, 300);
    }
    setIsSearchCleared(false);
  }, [products]);

  const handleSearchChange = (e) => {
    const newValue = e.target.value;
    setGlobalSearchTerm(newValue);
    setLastId("");

    if (!newValue) {
      setIsSearchCleared(true);
      window.scrollTo({ top: 0, behavior: "smooth" });
    }
  };

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
      <div className="mb-6 relative flex items-center">
        <input
          ref={searchInputRef}
          type="text"
          placeholder="Search products..."
          value={globalSearchTerm}
          onChange={handleSearchChange}
          className="w-full p-3 pl-12 border-4 border-[#3B82F6] rounded-lg shadow-md transition-all duration-300 focus:outline-none focus:ring-4 focus:ring-[#3B82F6] focus:border-[#3B82F6] focus:shadow-xl hover:shadow-lg text-base"
        />
        <div className="absolute left-3 text-[#3B82F6]">
          <FaSearch size={20} />
        </div>
      </div>
      <h2 className="text-3xl font-semibold mb-6">
        {globalSearchTerm ? "Search Results" : "All Products"}
      </h2>

      <motion.div
        className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6"
        initial="hidden"
        animate="visible"
        variants={{
          hidden: { opacity: 0 },
          visible: { opacity: 1, transition: { staggerChildren: 0.1, delayChildren: 0.1 } },
        }}
      >
        {products.length > 0 ? (
          products.map((product) => (
            <motion.div
              key={product.productId}
              variants={{
                hidden: { opacity: 0, scale: 0.9 },
                visible: { opacity: 1, scale: 1, transition: { duration: 0.4, delay: 0.05 } },
              }}
            >
              <ProductCard product={product} />
            </motion.div>
          ))
        ) : (
          <div className="text-lg text-gray-500 col-span-full text-center">
            No products found
          </div>
        )}
      </motion.div>

      {!globalSearchTerm && (
        <div className="flex justify-center mt-8">
          {!hasMoreProducts ? (
            <div className="text-lg text-gray-500">No more products</div>
          ) : (
            <div className="flex flex-col items-center">
              <button
                ref={loadMoreButtonRef}
                onClick={loadMoreProducts}
                className="bg-[#F58B00] text-white px-6 py-2 rounded-lg font-semibold transition-all duration-300 cursor-pointer hover:bg-[#FFC67D] hover:shadow-lg focus:outline-none focus:ring-4 focus:ring-[#FFC67D] focus:ring-opacity-50"
                whileHover={{ scale: 1.05, backgroundColor: "#FFC67D" }}
              >
                Load More
              </button>

              {loading && (
                <div className="mt-4">
                  <div className="animate-spin rounded-full h-10 w-10 border-t-4 border-[#F58B00] border-solid"></div>
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default ProductsPage;