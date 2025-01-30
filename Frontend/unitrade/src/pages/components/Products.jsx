import React from 'react'
import Card from './Card'
import Slider from "react-slick"

const Products = () => {
    var settings = {
        dots: true,
        infinite: true,
        speed: 500,
        slidesToShow: 5,
        slidesToScroll: 1,
        responsive: [
            {
            breakpoint: 1024,
            settings: {
                slidesToShow: 3,
                slidesToScroll: 3,
                infinite: true,
                dots: true
            }
            },
            {
            breakpoint: 600,
            settings: {
                slidesToShow: 2,
                slidesToScroll: 2
            }
            },
            {
            breakpoint: 480,
            settings: {
                slidesToShow: 1,
                slidesToScroll: 1
            }
            }
            // You can unslick at a given breakpoint now by adding:
            // settings: "unslick"
            // instead of a settings object
        ]
        
      };
  return (
    <div className='w-full bg-white py-24'>
        <div className='md:max-w-[2100px] m-auto max-w-[600px]'>
            <h1 className='py-5 text-3xl font-bold'> Browse products</h1>
            <Slider {...settings}>
                <Card/>
                <Card/>
                <Card/>
                <Card/>
            </Slider>
        </div>

    </div>
  )
}

export default Products