import React from "react";
import bannerImg from "../assets/imgs/banner_image.svg";
import { motion } from "framer-motion";
import { useSearchContext } from "../context/SearchContext";
import { useNavigate } from "react-router-dom";
import { FaSearch } from "react-icons/fa";


const Banner = () => {
  const { searchTerm, setSearchTerm } = useSearchContext();
  const navigate = useNavigate();

  const handleSearchChange = (e) => {
    setSearchTerm(e.target.value);
  };

  const handleKeyDown = (e) => {
    if (e.key === "Enter") {
      navigate("/products");
    }
  };

  return (
    <div className="w-100vh h-[70vh] relative py-24">
      <img
        src={bannerImg}
        className="absolute inset-0 w-full h-full object-cover"
        alt="Banner"
      />

      <div className="w-full m-auto relative z-10 py-4">
        <div className="flex flex-col items-center justify-center h-full gap-4">
          <h1 className="text-[96px] font-serif font-Playfair font-semibold py-4">
            <motion.span
              initial={{ x: -100, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              transition={{ duration: 1, ease: "easeInOut" }}
              className="text-[#FA4616] inline-block"
            >
              Uni
            </motion.span>
            <motion.span
              initial={{ x: 100, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              transition={{ duration: 1, ease: "easeInOut" }}
              className="text-[#0021A5] inline-block"
            >
              Bazaar
            </motion.span>
          </h1>
          <p className="text-[32px] text-[#D9D9D9] font-Raleway py-4">
            Connecting students for buying/selling
          </p>

          <div className="relative flex items-center w-full max-w-[600px] bg-white/50 rounded-lg">
            <div className="absolute left-3 text-gray-500">
              <FaSearch size={20} />
            </div>
            <input
              type="text"
              placeholder="Search for items..."
              value={searchTerm}
              onChange={handleSearchChange}
              onKeyDown={handleKeyDown}
              className="w-full pl-10 pr-12 py-2 bg-transparent rounded-lg focus:outline-none text-gray-800 placeholder-gray-600"
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Banner;
