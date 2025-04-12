import { useState, useRef, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { generateStars } from "@/utils/generateStar";
import { useLocation } from "react-router-dom";
import { FiMoreVertical } from "react-icons/fi";

const ProductCard = ({ product, onClick }) => {
  const [isHovered, setIsHovered] = useState(false);
  const [menuVisible, setMenuVisible] = useState(false);
  const location = useLocation();
  const cardRef = useRef(null);

  useEffect(() => {
    const handleClick = (event) => {
      if (cardRef.current && !cardRef.current.contains(event.target)) {
        setMenuVisible(false);
      }
    };

    document.addEventListener("mousedown", handleClick);
    return () => {
      document.removeEventListener("mousedown", handleClick);
    };
  }, []);

  return (
    <div
      ref={cardRef}
      className="relative flex flex-col w-full max-w-sm border border-gray-300 rounded-xl shadow-lg overflow-hidden bg-gray-900 transition-transform transform hover:scale-105 hover:shadow-2xl"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Three Dots Menu (always visible on /userproducts, not dependent on hover) */}
      {location.pathname === "/userproducts" && (
        <div className="absolute top-3 right-3 z-20">
          <button
            onClick={() => setMenuVisible(!menuVisible)}
            className="text-white hover:text-gray-400 cursor-pointer"
          >
            <FiMoreVertical size={20} />
          </button>
          {menuVisible && (
            <div className="absolute right-0 mt-2 w-24 bg-gray-800 text-white rounded-md shadow-lg z-30">
              <div className="px-4 py-2 hover:bg-gray-700 cursor-pointer">Edit</div>
              <div className="px-4 py-2 hover:bg-gray-700 cursor-pointer text-red-400">Delete</div>
            </div>
          )}
        </div>
      )}

      {/* Image Section with Title Overlay */}
      <div className="relative h-64 w-full overflow-hidden">
        <img
          className="w-full h-full object-cover transition-all duration-500"
          src={product.productImage}
          alt={product.productTitle}
        />
        <div
          className={`absolute bottom-0 left-0 w-full bg-gradient-to-t from-black to-transparent p-4 text-white transition-opacity ${isHovered ? "opacity-0" : "opacity-100"
            }`}
        >
          <h5 className="text-lg font-semibold tracking-tight">{product.productTitle}</h5>
        </div>
      </div>

      {/* Hover View - Detailed Information */}
      <div
        className={`absolute top-0 left-0 w-full h-full bg-black/40 backdrop-blur-md text-white p-5 flex flex-col justify-between opacity-0 transition-opacity duration-500 ${isHovered ? "opacity-100" : "opacity-0 pointer-events-none"
          }`}
      >
        <div>
          <h5 className="text-xl font-semibold">{product.productTitle}</h5>
          <p className="text-white text-sm mt-2 font-medium">{product.productDescription}</p>
        </div>

        <div>
          {/* Rating */}
          <div className="flex items-center">
            <div className="flex space-x-1">{generateStars(product.productCondition)}</div>
            <span className="ml-2 bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-0.5 rounded-md">
              5.0
            </span>
          </div>

          {/* Price & Button */}
          <div className="flex items-center justify-between mt-4">
            <span className="text-2xl font-bold text-[#F58B00]">${product.productPrice}</span>

            {/* Show the button only on /products page */}
            {location.pathname === "/products" && (
              <Button
                className="bg-[#F58B00] hover:bg-[#FFC67D] text-white font-bold py-2 px-4 rounded-lg shadow-md transition-all hover:shadow-lg hover:text-black cursor-pointer"
                onClick={onClick}
              >
                Message
              </Button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProductCard;
