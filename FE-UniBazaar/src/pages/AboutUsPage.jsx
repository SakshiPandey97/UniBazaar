import React from "react";
import { motion } from "framer-motion";
import TeamMemberCard from "@/customComponents/TeamMemberCard";
import SAKSHI_ICON from "@/assets/imgs/sakshi_icon.png";
import SHUBHAM_ICON from "@/assets/imgs/shubham_icon.jpg";
import AVANEESH_ICON from "@/assets/imgs/avaneesh_icon.webp";
import TANMAY_ICON from "@/assets/imgs/tanmay_icon.webp";

const AboutUsPage = () => {
  const team = [
    { name: "Sakshi Pandey", image: SAKSHI_ICON },
    { name: "Avaneesh Khandekar", image: AVANEESH_ICON },
    { name: "Tanmay Saxena", image: TANMAY_ICON },
    { name: "Shubham Singh", image: SHUBHAM_ICON },
  ];

  return (
    <div className="bg-white text-gray-900">
      {/* Hero */}
      <section className="relative bg-gradient-to-br from-[#F58B00] to-[#FFC67D] px-6 py-28 text-white text-center overflow-hidden">
        <motion.h1
          initial={{ opacity: 0, y: -40 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1 }}
          className="text-5xl font-extrabold tracking-tight"
        >
          About UniBazaar
        </motion.h1>
        <p className="mt-6 text-lg max-w-3xl mx-auto text-white/90">
          Redefining student commerce with a sustainable and local-first approach.
        </p>
      </section>

      {/* Our Mission Section */}
      <section className="py-20 px-6 bg-gray-50">
        <div className="max-w-6xl mx-auto flex flex-col md:flex-row items-center gap-10">
          <motion.div
            initial={{ opacity: 0, y: 40 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="text-center md:text-left flex-1"
          >
            <h2 className="text-4xl font-bold text-gray-800">Our Mission</h2>
            <p className="mt-4 text-lg text-gray-600 mb-10">
              UniBazaar empowers university students to engage in sustainable commerce by
              providing a localized platform for buying and selling essentials.
            </p>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
              <div className="bg-white shadow-md p-6 rounded-xl">
                <h3 className="text-xl font-semibold mb-2">🌱 Sustainability First</h3>
                <p className="text-gray-600 text-sm">Encouraging reuse and reducing campus waste.</p>
              </div>
              <div className="bg-white shadow-md p-6 rounded-xl">
                <h3 className="text-xl font-semibold mb-2">🔗 Community Driven</h3>
                <p className="text-gray-600 text-sm">Built to connect students through trust and purpose.</p>
              </div>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Why UniBazaar Section */}
      <section className="py-20 px-6 bg-white">
        <div className="max-w-6xl mx-auto flex flex-col md:flex-row-reverse items-center gap-10">
          <motion.div
            initial={{ opacity: 0, y: 40 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="text-center md:text-left flex-1"
          >
            <h2 className="text-4xl font-bold text-gray-800">Why UniBazaar</h2>
            <p className="mt-4 text-lg text-gray-600 mb-10">
              Fast, safe, and student-focused. Everything you need—just a few steps away.
            </p>
            <div className="grid grid-cols-2 gap-6 text-center">
              <div>
                <p className="text-4xl font-bold text-[#F58B00]">2K+</p>
                <p className="text-gray-600 mt-1 text-sm">Items Exchanged</p>
              </div>
              <div>
                <p className="text-4xl font-bold text-[#F58B00]">15+</p>
                <p className="text-gray-600 mt-1 text-sm">Campuses Onboarded</p>
              </div>
              <div>
                <p className="text-4xl font-bold text-[#F58B00]">80%</p>
                <p className="text-gray-600 mt-1 text-sm">Lower Carbon Impact</p>
              </div>
              <div>
                <p className="text-4xl font-bold text-[#F58B00]">100%</p>
                <p className="text-gray-600 mt-1 text-sm">Student-Run</p>
              </div>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Our Vision Section */}
      <section className="py-20 px-6 bg-gray-50">
        <div className="max-w-6xl mx-auto flex flex-col md:flex-row items-center gap-10">
          <motion.div
            initial={{ opacity: 0, y: 40 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="text-center md:text-left flex-1"
          >
            <h2 className="text-4xl font-bold text-gray-800">Our Vision</h2>
            <p className="mt-4 text-lg text-gray-600">
              To make the circular economy a campus culture.
            </p>
            <div className="mt-6 bg-white shadow-lg p-6 rounded-xl max-w-xl">
              <p className="italic text-gray-700">
                “UniBazaar changed how I shop on campus. It's sustainable, local, and just works.”
              </p>
              <p className="mt-4 font-semibold text-right text-gray-800">– Vivek, UF Student</p>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Team Section */}
{/* Team Section */}
<section className="py-24 bg-white">
  <div className="text-center mb-16">
    <h2 className="text-4xl font-bold text-gray-800">Meet the Team</h2>
    <p className="mt-2 text-gray-600">
      The dreamers, builders, and doers behind UniBazaar
    </p>
  </div>

  <div className="max-w-6xl mx-auto grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-10 px-6">
    {team.map((member, i) => (
      <motion.div
        key={i}
        initial={{ opacity: 0, y: 30 }}
        whileInView={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: i * 0.1 }}
        whileHover={{ scale: 1.03 }}
        className="bg-gray-50 shadow-lg rounded-xl overflow-hidden text-center hover:shadow-xl transition"
      >
        <div className="p-6 flex flex-col items-center">
          <img
            src={member.image}
            alt={member.name}
            className="w-28 h-28 object-cover rounded-full border-4 border-white shadow-md mb-4"
          />
          <h3 className="text-xl font-semibold text-gray-800">{member.name}</h3>
          <p className="text-sm text-gray-500 mt-1">Full-Stack Developer</p>
          
        </div>
      </motion.div>
    ))}
  </div>
</section>

    </div>
  );
};

export default AboutUsPage;
