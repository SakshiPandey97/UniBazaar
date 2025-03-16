import { useState } from "react";
import { Button } from "@/components/ui/button";
import { generateStars } from "@/utils/generateStar";

const ProductCard = ({ product,onClick }) => {
  const [isHovered, setIsHovered] = useState(false);
  console.log(product)
  return (
    <div
      className="relative flex flex-col w-full max-w-sm border border-gray-300 rounded-xl shadow-lg overflow-hidden bg-gray-900 transition-transform transform hover:scale-105 hover:shadow-2xl"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Image Section with Title Overlay */}
      <div className="relative h-64 w-full overflow-hidden">
        <img
          className="w-full h-full object-cover transition-all duration-500"
          src={product.productImage}
          alt={product.productTitle}        />
        <div
          className={`absolute bottom-0 left-0 w-full bg-gradient-to-t from-black to-transparent p-4 text-white transition-opacity ${
            isHovered ? "opacity-0" : "opacity-100"
          }`}
        >
          <h5 className="text-lg font-semibold tracking-tight">{product.productTitle}</h5>
        </div>
      </div>

      {/* Hover View - Detailed Information */}
      <div
        className={`absolute top-0 left-0 w-full h-full bg-black/40 backdrop-blur-md text-white p-5 flex flex-col justify-between opacity-0 transition-opacity duration-500 ${
          isHovered ? "opacity-100" : "opacity-0 pointer-events-none"
        }`}
      >
        <div>
          <h5 className="text-xl font-semibold">{product.productTitle}</h5>
          <p className="text-gray-400 text-sm mt-2">{product.productDescription}</p>
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
            <Button className="bg-[#F58B00] hover:bg-[#FFC67D] text-white font-bold py-2 px-4 rounded-lg shadow-md transition-all hover:shadow-lg hover:text-black"
            onClick={onClick}>
              Message
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProductCard;
