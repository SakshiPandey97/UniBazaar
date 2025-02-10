import React, { useEffect, useState } from 'react'
import Card from './Card'
import Slider from "react-slick"
import axios from 'axios'
import { useSearchParams } from 'react-router-dom'

// Fetch the product data from backend
// Pass this product data to the card component
// We will use axios to transform fetched json data, instead of using fetch() which requires manual handling of JSON parsing and headers.

const Products = () => {
    const [products, setProducts] = useState([]);
    const [isLeftArrowDisabled, setIsLeftArrowDisabled] = useState(true);


    useEffect(() => {
        axios.get("http://192.168.0.203:8080/products")
            .then(response => {
                
                // here we will flatten the nested product array
                const allProducts = response.data.flatMap(user => user.Products);
                
                // Sort products by most recent post date
                const sortedProducts = allProducts.sort((a, b) => 
                    new Date(b.ProductPostDate) - new Date(a.ProductPostDate)
                );

                setProducts(allProducts);
            })
            .catch(error => {
                console.error('Error fetching products:', error);
            });

    }, []);

    const settings = {
        dots: true,
        infinite: true,
        speed: 500,
        slidesToShow: 5,
        slidesToScroll: 1,
        initialSlide: 0,
        responsive: [
            {
            breakpoint: 1024,
            settings: { slidesToShow: 3, slidesToScroll: 3, infinite: true, dots: true },
            },
            {
            breakpoint: 600,
            settings: { slidesToShow: 2, slidesToScroll: 2}
            },
            {
            breakpoint: 480,
            settings: {  slidesToShow: 1,  slidesToScroll: 1 }
            }
            // You can unslick at a given breakpoint now by adding:
            // settings: "unslick"
            // instead of a settings object
        ],
        beforeChange: () => {
            setIsLeftArrowDisabled(false); // Enable left arrow when user starts sliding
        },
        afterChange: (current) => {
            if (current === 0) {
                setIsLeftArrowDisabled(true); // Disable left arrow when at the first slide
            }
        }
        
      };


  return (
    <div className='w-full bg-white py-24'>
        <div className='md:max-w-[2100px] m-auto max-w-[600px]'>
            <h1 className='py-5 text-3xl text-teal-600 font-bold'> Browse products</h1>
            <Slider {...settings}>
                {products.map(product => 
                    <Card
                        key={product.ProductId}
                        title={product.ProductTitle}
                        price={product.ProductPrice}
                        image={product.ProductImage}
                    />
                )}
            </Slider>
        </div>

    </div>
  )
}

export default Products