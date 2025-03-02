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
