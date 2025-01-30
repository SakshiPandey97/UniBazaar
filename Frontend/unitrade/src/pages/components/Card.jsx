import React from 'react'
import prodImg1 from '../../assets/imgs/Prodcard-3.svg'

const Card = () => {
  return (
    <div className='bg-white drop-shadow-md overflow-hidden rounded-lg mr-4'>
        <img src={prodImg1}
            className='h-80 w-full object-cover'
        />
        <div className='p-5'>
            <h1> Nike Shoes</h1>
        </div>
        <div className='px-5 pb-5 flex justify-between items-center'>
            <h3> $500 </h3>
            <button className='bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700'>
                Buy Now
            </button>
        </div>

    </div>
  )
}
export default Card