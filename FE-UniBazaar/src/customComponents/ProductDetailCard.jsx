import React from 'react';
import { reverseProductConditionMapping, conditionColorMap } from '../utils/productMappings';

function ProductDetailCard({ product, onClick, isSelected }) {
  return (
    <div
      className={`bg-gray-100 rounded-lg shadow-lg overflow-hidden transition-transform duration-300 ease-in-out ${isSelected ? 'scale-200 z-20' : 'scale-100'
        }`}
    >
      <div className="relative">
        <div
          onClick={onClick}
          className="absolute inset-0 w-full h-full z-10 cursor-pointer"
        />
        <img
          src={product.productImage}
          alt={product.productTitle}
          className="w-full h-56 object-cover"
        />
      </div>
      <div className="p-4">
        <h3 className="text-xl font-semibold text-gray-900">{product.productTitle}</h3>
        <p className="text-gray-500 text-sm">{product.productDescription}</p>
        <div className="flex justify-between items-center mt-4">
          <span className="text-lg font-bold text-[#F58B00]">{product.getFormattedPrice()}</span>
          <span className="text-sm text-gray-500">{new Date(product.productPostDate).toLocaleDateString()}</span>
        </div>
        <div className="mt-2">
          <span
            className={`px-4 py-1 rounded-full text-sm ${conditionColorMap[product.productCondition]}`}
          >
            {reverseProductConditionMapping[product.productCondition]}
          </span>
        </div>

        <div className="mt-4">
          <button
            data_testid="messageSellerBtn"
            type="button"
            className="w-full hover:border-[#032B54] border-2 p-2 bg-[#032B54] hover:bg-[#021e39] text-white font-bold py-2 px-4 rounded-md transition disabled:bg-gray-400"
            onClick={() => alert('Messaging the seller...')}
          >
            Message Seller
          </button>
        </div>


      </div>
    </div>
  );
}

export default ProductDetailCard;
