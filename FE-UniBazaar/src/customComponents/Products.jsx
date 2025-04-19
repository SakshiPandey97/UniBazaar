import React from "react";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
import useStartChat from "../hooks/useStartChat";
import Spinner from "./Spinner";
import useFetchProducts from "../hooks/useFetchProducts";
import ProductCard from "./ProductCard";

const Products = () => {
  const handleStartChat = useStartChat();
  const { products, loading, error } = useFetchProducts();
  if (loading) return <Spinner />;
  if (error) return <div>{error}</div>;

  return (
    <div className="flex flex-col justify-center w-full py-16 bg-gradient-to-b from-gray-50 to-gray-200">
      <div className="flex flex-col m-auto max-w-[400px] sm:max-w-[500px] md:max-w-[750px] lg:max-w-[1000px] xl:max-w-[1250px] justify-center">
        <h1 className="py-5 text-4xl text-gray-900 font-extrabold text-center uppercase tracking-wide">
          Browse Products
        </h1>
        <Carousel className="relative flex bg-white border border-gray-200 rounded-xl p-6 shadow-lg transition-all hover:shadow-2xl">
          
          <CarouselPrevious className="absolute z-10 left-4 top-1/2 -translate-y-1/2 bg-gradient-to-r from-gray-800 to-gray-600 text-white p-3 rounded-full shadow-md hover:scale-110 transition-all flex items-center justify-center w-10 h-10">
            <span className="text-lg font-bold">&#8592;</span>
          </CarouselPrevious>

          <CarouselContent className="flex justify-between gap-x-6 px-4">
            {products.map((product) => (
              <CarouselItem
                key={product.productId}
                className="basis-full sm:basis-1/2 md:basis-1/3 lg:basis-1/4 xl:basis-1/5 transform transition-transform hover:scale-105 hover:shadow-lg"
              >
                <ProductCard product={product} onStartChat={handleStartChat} />
              </CarouselItem>
            ))}
          </CarouselContent>

          <CarouselNext className="absolute right-4 top-1/2 -translate-y-1/2 bg-gradient-to-l from-gray-800 to-gray-600 text-white p-3 rounded-full shadow-md hover:scale-110 transition-all flex items-center justify-center w-10 h-10">
            <span className="text-lg font-bold">&#8594;</span>
          </CarouselNext>

        </Carousel>
      </div>
    </div>
  );
};

export default Products;
