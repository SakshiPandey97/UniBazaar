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

  if (loading) return <Spinner />;
  if (error) return <div>{error}</div>;

  return (
    <div className="flex flex-col justify-center w-full py-24">
      <div className="flex flex-col m-auto max-w-[400px] sm:max-w-[500px] md:max-w-[750px] lg:max-w-[1000px] xl:max-w-[1250px] justify-center">
        <h1 className="py-5 text-3xl text-[#320B54] font-bold">
          Browse products
        </h1>
        <Carousel className="relative flex bg-[#FFC67D] border rounded-lg p-4">
          <CarouselPrevious className="absolute left-2 top-1/2 -translate-y-1/2 bg-white z-10" />
          <CarouselContent className="flex justify-between gap-x-4">
            {products.map((product) => (
              <CarouselItem
                key={product.productId}
                className="basis-full sm:basis-1/2 md:basis-1/3 lg:basis-1/4 xl:basis-1/5"
              >
                <ProductCard
                  title={product.productTitle}
                  price={product.productPrice}
                  condition={product.productCondition}
                  image={product.productImage}
                  description={product.productDescription}
                />
              </CarouselItem>
            ))}
          </CarouselContent>
          <CarouselNext className="absolute right-2 top-1/2 -translate-y-1/2 bg-white" />
        </Carousel>
      </div>
    </div>
  );
};

export default Products;
