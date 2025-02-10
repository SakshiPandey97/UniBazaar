import React from 'react'

// Use the product data from Products.jsx


const Card = ({ title, price, image }) => {
  return (
    <div className='bg-white drop-shadow-md overflow-hidden rounded-lg mr-4'>
        <img 
            src={image} // product image will change dynamically fetched from this url
            alt={title}
            className='h-80 w-full object-cover'
        />
        <div className='p-15'>
            <h1>{title}</h1>
        </div>
        <div className='px-15 pb-5 flex justify-between items-center'>
            <h3> ${price} </h3>
            {/* <button className='bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700'>
                Bid Now
            </button> */}
        </div>

    </div>
  )
}
export default Card