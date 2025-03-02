import { Card, CardHeader, CardDescription, CardFooter } from "@/components/ui/card";
import { reverseProductConditionMapping } from "../utils/productMappings"; 


const ProductCard = ({ title, price, condition, image, description  }) => {

  console.log(title, price, condition, image, description);
  const conditionString = reverseProductConditionMapping[condition] || "Unknown";

  return (
<Card className="h-100 w-full bg-white shadow-xl rounded-lg flex flex-col justify-between bg-gray-100">
  <CardHeader>
    <img
      src={image}
      alt={title}
      className="w-full h-40 object-cover rounded-lg"
    />
  </CardHeader>
  <CardDescription className="flex flex-col justify-center items-center text-center py-4">
    <h2 className="text-xl font-semibold">{title}</h2>
    <div className="flex flex-row"><p className="text-teal-600 mt-2 text-lg font-bold">$ </p> <p className="text-bg mt-2 text-lg font-bold">{price}</p></div>
    <p>Condition: {conditionString}</p>
  </CardDescription>
  <CardFooter>
    <button className="w-full px-4 py-2 bg-teal-600 text-white rounded-lg hover:bg-teal-700">
      Buy Now
    </button>
  </CardFooter>
</Card>

  );
};

export default ProductCard;
