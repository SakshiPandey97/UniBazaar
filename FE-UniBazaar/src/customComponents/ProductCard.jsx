import { useState, useRef, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { generateStars } from "@/utils/generateStar";
import { useLocation } from "react-router-dom";
import { FiMoreVertical, FiCheck, FiX } from "react-icons/fi";
import {
  productConditionMapping,
  reverseProductConditionMapping,
  PRODUCT_CONDITIONS,
  conditionColorMap,
} from "../utils/productMappings";
import useUpdateProduct from "../hooks/useUpdateProduct";
import useDeleteProduct from "../hooks/useDeleteProduct";
import { getCurrentUserId } from "@/utils/getUserId";

const ProductCard = ({
  product,
  onClick,
  onEdit,
  isEditing: propIsEditing,
  onCancel: propOnCancel,
  onProductUpdated,
  onProductDeleted
}) => {
  const [isHovered, setIsHovered] = useState(false);
  const [menuVisible, setMenuVisible] = useState(false);
  const location = useLocation();
  const cardRef = useRef(null);
  const [editableProduct, setEditableProduct] = useState(product);
  const [imagePreview, setImagePreview] = useState(product.productImage);
  const [newImageFile, setNewImageFile] = useState(null);
  const [isEditing, setIsEditing] = useState(propIsEditing);

  const { updateProduct, isSaving, error } = useUpdateProduct();
  const { deleteProduct, isLoading: isDeleting } = useDeleteProduct();

  useEffect(() => {
    setIsEditing(propIsEditing);
    if (!propIsEditing) {
      setEditableProduct(product);
      setImagePreview(product.productImage);
      setNewImageFile(null);
    } else {
      setEditableProduct(product);
      setImagePreview(product.productImage);
    }
  }, [propIsEditing, product]);

  useEffect(() => {
    const handleClick = (event) => {
      if (cardRef.current && !cardRef.current.contains(event.target)) {
        setMenuVisible(false);
        if (isEditing) {
          setIsEditing(false);
          setEditableProduct(product);
          setImagePreview(product.productImage);
          setNewImageFile(null);
          if (propOnCancel) {
            propOnCancel();
          }
        }
      }
    };

    document.addEventListener("mousedown", handleClick);
    return () => {
      document.removeEventListener("mousedown", handleClick);
    };
  }, [isEditing, product, propOnCancel]);

  const handleMenuToggle = () => {
    setMenuVisible(!menuVisible);
  };

  const handleEdit = () => {
    setIsEditing(true);
    if (onEdit) {
      onEdit(product.productId);
    }
    setMenuVisible(false);
  };

  const handleChange = (field, value) => {
    setEditableProduct((prev) => ({ ...prev, [field]: value }));
  };

  const handleImageUpload = (e) => {
    const file = e.target.files[0];
    if (file) {
      setNewImageFile(file);
      setImagePreview(URL.createObjectURL(file));
    } else {
      setImagePreview(product.productImage);
      setNewImageFile(null);
    }
  };

  const handleSave = async () => {
    const updatedProductData = { ...editableProduct };
    const condition = updatedProductData.productCondition;

    const updatedProduct = await updateProduct(updatedProductData, newImageFile, condition);

    if (updatedProduct) {
      if (onProductUpdated) {
        onProductUpdated(updatedProduct);
      }

      setEditableProduct(prev => ({
        ...prev,
        productImage: updatedProduct.productImage
      }));
      setImagePreview(updatedProduct.productImage);

      setIsEditing(false);
      setNewImageFile(null);
    } else {
      console.error("Failed to update product.");
    }
  };
  const handleCancel = () => {
    setEditableProduct(product);
    setImagePreview(product.productImage);
    setNewImageFile(null);
    setIsEditing(false);
    setMenuVisible(false);
    if (propOnCancel) {
      propOnCancel();
    }
  };

  const handleDelete = async () => {
    const userId = getCurrentUserId();
    if (!userId) {
      console.error("No user ID found.");
      return;
    }

    await deleteProduct(userId, product.productId);

    if (typeof onProductDeleted === "function") {
      onProductDeleted(product.productId);
    }
  };

  return (
    <div
      ref={cardRef}
      className={`relative flex flex-col w-full max-w-sm ${location.pathname === "/" ? "h-[22rem]" : "h-auto"
        } border border-gray-300 rounded-xl shadow-lg overflow-hidden bg-gray-900 transition-transform transform hover:scale-105 hover:shadow-2xl`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {location.pathname === "/userproducts" && (
        <div className="absolute top-3 right-3 z-20">
          <button
            onClick={handleMenuToggle}
            className="text-white hover:text-gray-400 cursor-pointer"
          >
            <FiMoreVertical size={20} />
          </button>
          {menuVisible && (
            <div className="absolute right-0 mt-2 w-24 bg-gray-800 text-white rounded-md shadow-lg z-30">
              {!isEditing ? (
                <div
                  onClick={handleEdit}
                  className="px-4 py-2 hover:bg-gray-700 cursor-pointer"
                >
                  Edit
                </div>
              ) : (
                <div className="px-4 py-2 text-gray-500 cursor-default">Edit</div>
              )}
              <div
                onClick={handleDelete}
                className={`px-4 py-2 hover:bg-gray-700 cursor-pointer text-red-400 ${isDeleting ? "opacity-50 pointer-events-none" : ""}`}
              >
                {isDeleting ? "Deleting..." : "Delete"}
              </div>
            </div>
          )}
        </div>
      )}

      {isEditing ? (
        <div className="p-4 space-y-2">
          <div className="relative h-64 w-full overflow-hidden rounded-md">
            <img
              className="w-full h-full object-cover"
              src={imagePreview}
              alt={editableProduct.productTitle}
            />
            <div className="absolute bottom-0 left-0 w-full bg-black bg-opacity-60 p-2 text-white text-sm">
              <label className="flex flex-col items-center cursor-pointer">
                <span className="mb-1">Upload New Image</span>
                <input
                  type="file"
                  accept="image/*"
                  onChange={handleImageUpload}
                  className="hidden"
                />
                <span className="bg-gray-700 px-3 py-1 mt-1 rounded-md">Choose File</span>
              </label>
            </div>
          </div>
          <input
            type="text"
            value={editableProduct.productTitle}
            onChange={(e) => handleChange("productTitle", e.target.value)}
            className="w-full bg-gray-800 p-2 rounded-md text-white font-semibold"
            placeholder="Product Title"
          />
          <textarea
            value={editableProduct.productDescription}
            onChange={(e) => handleChange("productDescription", e.target.value)}
            className="w-full bg-gray-800 p-2 rounded-md text-white"
            placeholder="Description"
          />
          <input
            type="number"
            value={editableProduct.productPrice}
            onChange={(e) => handleChange("productPrice", e.target.value)}
            className="w-full bg-gray-800 p-2 rounded-md text-white"
            placeholder="Price"
          />
          <select
            value={editableProduct.productCondition}
            onChange={(e) => handleChange("productCondition", parseInt(e.target.value))}
            className="w-full bg-gray-800 p-2 rounded-md text-white"
          >
            <option value="" disabled>Select Condition</option>
            {PRODUCT_CONDITIONS.map((condition) => (
              <option key={productConditionMapping[condition]} value={productConditionMapping[condition]}>
                {condition}
              </option>
            ))}
          </select>
          <div className="flex justify-end space-x-2">
            <button onClick={handleSave} disabled={isSaving} className={`bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-400 ${isSaving ? 'opacity-50 cursor-not-allowed' : ''}`}>
              <FiCheck className="inline-block mr-1" /> {isSaving ? 'Saving...' : 'Save'}
            </button>
            <button onClick={handleCancel} disabled={isSaving} className="bg-gray-500 text-white px-4 py-2 rounded-md hover:bg-gray-400">
              <FiX className="inline-block mr-1" /> Cancel
            </button>
          </div>
          {error && <p className="text-red-500 text-sm">{error}</p>}
        </div>
      ) : (
        <>
          <div
            className={`relative w-full overflow-hidden ${location.pathname === "/" ? "h-full" : "h-64"}`}
          >
            <img
              className="w-full h-full object-cover transition-all duration-500"
              src={editableProduct.productImage}
              alt={editableProduct.productTitle}
            />
            <div
              className={`absolute bottom-0 left-0 w-full bg-gradient-to-t from-black to-transparent p-4 text-white transition-opacity ${isHovered ? "opacity-0" : "opacity-100"
                }`}
            >
              <h5 className="text-lg font-semibold tracking-tight">{editableProduct.productTitle}</h5>
            </div>
          </div>

          <div
            className={`absolute top-0 left-0 w-full h-full bg-black/40 backdrop-blur-md text-white p-5 flex flex-col justify-between opacity-0 transition-opacity duration-500 ${isHovered ? "opacity-100" : "opacity-0 pointer-events-none"
              }`}
          >
            <div>
              <h5 className="text-xl font-semibold">{editableProduct.productTitle}</h5>
              <p className="text-white text-sm mt-2 font-medium">{editableProduct.productDescription}</p>
            </div>

            <div>
              <div className="flex items-center">
                <div className="flex space-x-1">{generateStars(editableProduct.productCondition)}</div>
                {reverseProductConditionMapping[editableProduct.productCondition] && (
                  <span className={`ml-2 text-xs font-semibold px-2.5 py-0.5 rounded-md ${conditionColorMap[editableProduct.productCondition]}`}>
                    {reverseProductConditionMapping[editableProduct.productCondition]}
                  </span>
                )}
              </div>

              <div className="flex items-center justify-between mt-4">
                <span className="text-2xl font-bold text-[#F58B00]">${editableProduct.productPrice}</span>

                {(location.pathname === "/" || location.pathname === "/products") && (
                  <Button
                    className="bg-[#F58B00] hover:bg-[#FFC67D] text-white font-bold py-2 px-4 rounded-lg shadow-md transition-all hover:shadow-lg hover:text-black cursor-pointer"
                    onClick={onClick}
                  >
                    Message
                  </Button>
                )}
              </div>
            </div>
          </div>
        </>
      )}
    </div>
  );
};

export default ProductCard;