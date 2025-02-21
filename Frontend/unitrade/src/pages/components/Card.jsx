import React from "react";

const Card = ({ title, price, image, description }) => {
  return (
    // <div
    //   className='h-full w-full bg-gray-400 drop-shadow-md overflow-hidden rounded-lg  cursor-pointer'  // Increased right margin
    //   onClick={onClick}
    // >
    //   <img
    //     src={image}
    //     alt={title}
    //     className='h-64 w-64 object-cover'
    //   />
    //   <div className='p-4'>
    //     <h1 className='text-xl font-semibold'>Product Title: {title}</h1>
    //   </div>
    //   <div className='px-4 pb-4 flex justify-between items-center'>
    //     <h3 className='text-lg font-semibold'>Price: ${price}</h3>
    //   </div>
    // </div>
    
      <div class="h-full w-full bg-gray-100 rounded overflow-hidden shadow-lg">
        <img class="w-full" src={image} alt={title} />
        <div class="flex flex-col px-6 py-4">
          <div class="font-bold text-xl mb-2 flex justify-center">{title}</div>
          <p class="text-gray-700 text-base filter blur-sm">
          {description}
          </p>
        </div>
        <div class="px-6 pt-4 pb-2">
          <span class="inline-block bg-gray-200 rounded-full px-3 py-1 text-sm font-semibold text-green-700 mr-2 mb-2 ">
           ${price}
          </span>
        </div>
      </div>
    
  );
};

export default Card;
