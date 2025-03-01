import React from "react";
import uploadImg from "@/assets/imgs/correct-upload.svg"
function ShowProductPreview({ productData }) {
  console.log(productData.productImage);
  return (
    <div className="modal-overlay">
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="w-full max-w-sm bg-white border border-gray-200 rounded-lg shadow-sm dark:bg-gray-800 dark:border-gray-700">
          <img
            className="p-8 rounded-t-lg"
            src={uploadImg}
            alt="product image"
          />
          <div className="px-5 pb-5">
            <h5 className="text-xl font-bold tracking-tight text-white">
              Product Uploaded Successfully
            </h5>
            <h5 className="text-sm font-semibold tracking-tight text-[#FFC67D]">
              Redirecting to Home Page
            </h5>
          </div>
        </div>{" "}
      </div>
    </div>
  );
}

export default ShowProductPreview;
