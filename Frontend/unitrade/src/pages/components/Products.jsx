import React, { useEffect, useState } from 'react';
import Card from './Card';
import Slider from "react-slick";
import axios from 'axios';
import { useSearchParams } from 'react-router-dom';

const Products = () => {
    const [products, setProducts] = useState([]);
    const [isLeftArrowDisabled, setIsLeftArrowDisabled] = useState(true);

    useEffect(() => {
        axios.get("http://127.0.0.1:8080/products")
            .then(response => {
                // No need to flatten, response is already flat
                const sortedProducts = response.data.sort((a, b) => 
                    new Date(b.productPostDate) - new Date(a.productPostDate)
                );

                setProducts(sortedProducts);
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
                settings: { slidesToShow: 2, slidesToScroll: 2 }
            },
            {
                breakpoint: 480,
                settings: { slidesToShow: 1, slidesToScroll: 1 }
            }
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
                            key={product.productId}
                            title={product.productTitle}
                            price={product.productPrice}
                            image={product.productImage}
                        />
                    )}
                </Slider>
            </div>
        </div>
    );
};

export default Products;
