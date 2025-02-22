import React from "react";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
import Spinner from "./Spinner";
import useProducts from "../hooks/useProducts";
import ProductCard from "./ProductCard";

const Products = () => {
  const { products, loading, error } = useProducts();

  // Log the products array to the console for debugging
  console.log("Products: ", products);

  if (loading) return <Spinner />;
  if (error) return <div>{error}</div>;

  return (
    <div className="flex flex-col justify-center w-full py-24">
      <div className="flex flex-col md:max-w-[1250px] m-auto max-w-[600px] justify-center">
        <h1 className="py-5 text-3xl text-teal-600 font-bold">
          Browse products
        </h1>
        <Carousel>
          <CarouselPrevious />
          <CarouselContent className="flex justify-between gap-x-4">
            {products.map((product) => {
              console.log("Product: ", product); // Log each individual product
              return (
                <CarouselItem
                  key={product.productId}
                  className="basis-1/3 md:basis-1/4"
                >
                  <ProductCard
                    title={product.productTitle}
                    price={product.productPrice}
                    image={product.productImage}
                    description={product.productDescription}
                  />
                </CarouselItem>
              );
            })}
          </CarouselContent>
          <CarouselNext />
        </Carousel>
      </div>
    </div>
  );
};

export default Products;
