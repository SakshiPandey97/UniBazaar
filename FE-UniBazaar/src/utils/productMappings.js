export const productConditionMapping = {
  "Excellent": 5,
  "Very Good": 4,
  "Good": 3,
  "Fair": 2,
  "Poor": 1,
};

export const reverseProductConditionMapping = Object.fromEntries(
  Object.entries(productConditionMapping).map(([key, value]) => [value, key])
);
export const PRODUCT_CONDITIONS = Object.keys(productConditionMapping);

export const conditionColorMap = {
  5: "bg-green-600 text-green-2000", // Excellent
  4: "bg-green-300 text-green-800",  // Very Good
  3: "bg-yellow-300 text-yellow-800", // Good
  2: "bg-orange-200 text-orange-800", // Fair
  1: "bg-red-200 text-red-800",  // Poor
};
