import React from "react";
import Card from './Card';
import Slider from "react-slick";
import Spinner from './Spinner';
import useProducts from "../../hooks/useProducts";

const Products = () => {
  const { products, loading, error } = useProducts();  

  const settings = {
    dots: true,
    infinite: true,
    speed: 500,
    slidesToShow: 5,
    slidesToScroll: 1,
    initialSlide: 0,
    lazyLoad: 'ondemand',
    responsive: [
      { breakpoint: 1024, settings: { slidesToShow: 3, slidesToScroll: 3, infinite: true, dots: true } },
      { breakpoint: 600, settings: { slidesToShow: 2, slidesToScroll: 2 } },
      { breakpoint: 480, settings: { slidesToShow: 1, slidesToScroll: 1 } }
    ]
  };

  if (loading) return <Spinner />; 
  if (error) return <div>{error}</div>; 

  return (
    <div className='flex flex-col justify-center w-full py-24'>
      <div className='flex flex-col md:max-w-[1250px] m-auto max-w-[600px]'>
        <h1 className='py-5 text-3xl text-teal-600 font-bold'>Browse products</h1>
        <Slider {...settings}>
          {products.map(product => 
            <div key={product.productId} >  {/* Added margin right to each card */}
              <Card
                title={product.productTitle}
                price={product.productPrice}
                image={product.productImage}
                description={product.productDescription}
              />
            </div>
          )}
        </Slider>
      </div>
    </div>
  );
};

export default Products;
