import React from "react";
import bannerImg from "../assets/imgs/banner_image.jpg";
import { motion } from "framer-motion";
import { useSearchContext } from "../context/SearchContext";
import { useNavigate } from "react-router-dom";
import { FaSearch } from "react-icons/fa";

const Banner = () => {
  const { searchTerm, setSearchTerm } = useSearchContext();
  const navigate = useNavigate();

  const handleSearchChange = (e) => setSearchTerm(e.target.value);
  const handleKeyDown = (e) => {
    if (e.key === "Enter") navigate("/products");
  };

  return (
    <div className="relative h-[70vh] w-full overflow-hidden">
      {/* Background illustration */}
      <img
        loading="lazy"
        src={bannerImg}
        alt="Banner illustration"
        className="absolute inset-0 h-full w-full object-cover"
      />

      {/* Darkening overlay for better text contrast */}
      <div className="absolute inset-0 bg-gradient-to-b from-black/40 via-black/20 to-transparent" />

      {/* Content */}
      <div className="relative z-10 flex h-full flex-col items-center justify-center px-4 text-center">
        {/* UniBazaar word‑mark with fine white outline */}
        <h1 className="drop-shadow-[0_4px_6px_rgba(0,0,0,0.35)] font-serif font-Playfair font-semibold leading-tight text-6xl sm:text-7xl lg:text-8xl [-webkit-text-stroke:0.75px_white]">
          <motion.span
            initial={{ x: -100, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ duration: 1 }}
            className="text-[#033460]"
          >
            Uni
          </motion.span>
          <motion.span
            initial={{ x: 100, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ duration: 1 }}
            className="text-[#FA4616]"
          >
            Bazaar
          </motion.span>
        </h1>

        {/* Tagline in Playfair italic & bold with blue outline */}
        <p className="mb-6 mt-4 text-2xl sm:text-3xl font-serif font-Playfair font-bold text-[#FA4616] [-webkit-text-stroke:0.5px_#033460]">
          Declutter, Discover, Deal&nbsp;
        </p>

        {/* Search bar */}
        <div className="relative w-full max-w-[600px] rounded-lg bg-white/20 backdrop-blur-md shadow-md">
          <FaSearch
            size={20}
            className="absolute left-3 top-1/2 -translate-y-1/2 text-white/80"
          />
          <input
            type="text"
            placeholder="Search for items..."
            value={searchTerm}
            onChange={handleSearchChange}
            onKeyDown={handleKeyDown}
            className="w-full bg-transparent py-2 pl-10 pr-4 text-white placeholder-white/70 focus:outline-none"
          />
        </div>
      </div>
    </div>
  );
};

export default Banner;