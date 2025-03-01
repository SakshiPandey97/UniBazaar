import { Button } from "@/components/ui/button";
import { generateStars } from "@/utils/generateStar";

const ProductCard = ({ title, price, condition, image, description }) => {
  console.log(title, price, condition, image, description);

  return (
    <div className="flex flex-col w-full h-85 max-w-sm border border-black rounded-lg shadow-sm bg-[#032B54]">
      <div className="flex flex-col h-1/2 overflow-hidden">
        <img
          className="rounded-t-lg w-full h-full object-cover"
          src={image}
          alt={title}
        />
      </div>
      <div className="flex flex-col h-1/2">
        <div className="px-5 pb-5">
          <h5 className="text-xl font-semibold tracking-tight text-gray-900 dark:text-white">
            {title}
          </h5>
          <div className="flex items-center mt-2.5 mb-5">
            <div className="flex items-center space-x-1 rtl:space-x-reverse">
              {generateStars(condition)}
            </div>
            <span className="bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-0.5 rounded-sm dark:bg-blue-200 dark:text-blue-800 ms-3">
              5.0
            </span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-3xl font-bold text-white ">
              ${price}
            </span>
            <Button className="hover:border-[#F58B00] border-2 p-2 bg-[#F58B00] hover:bg-[#FFC67D] text-balck font-bold py-2 px-4">
              Read More
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProductCard;
