import React from "react";
import { motion } from "framer-motion";
import TeamMemberCard from "@/customComponents/TeamMemberCard";
import SAKSHI_ICON from"@/assets/imgs/sakshi_icon.png"
import SHUBHAM_ICON from"@/assets/imgs/shubham_icon.jpg"
import AVANEESH_ICON from"@/assets/imgs/avaneesh_icon.webp"
import TANMAY_ICON from"@/assets/imgs/tanmay_icon.webp"

const AboutUsPage = () => {
  return (
    <div className="bg-gray-50 text-gray-900">
      {/* Hero Section */}
      <section className="relative bg-gradient-to-r from-[#F58B00] to-[#FFC67D] py-20 px-6 text-center text-white">
        <motion.h1
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1 }}
          className="text-5xl font-extrabold tracking-wide"
        >
          About Us
        </motion.h1>
        <p className="mt-4 text-lg max-w-2xl mx-auto">
          Empowering businesses with innovative solutions and cutting-edge technology.
        </p>
      </section>

      {/* Mission Section */}
      <section className="py-16 px-6 text-center">
        <motion.h2
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          className="text-4xl font-bold text-gray-800"
        >
          Our Mission
        </motion.h2>
        <p className="mt-4 text-lg max-w-3xl mx-auto text-gray-600">
          We aim to revolutionize the industry by providing top-notch services that drive success
          and bring efficiency to businesses worldwide.
        </p>
      </section>

      {/* Team Section */}
      <section className="py-16 bg-gray-100">
        <div className="text-center">
          <h2 className="text-4xl font-bold text-gray-800">Meet Our Team</h2>
          <p className="text-gray-600 mt-2">The talented individuals behind our success</p>
        </div>
        <div className="flex flex-wrap justify-center mt-10 gap-6">
          <TeamMemberCard
            name="Sakshi Pandey"
            position="CEO & Founder"
            image={SAKSHI_ICON}
          />
          <TeamMemberCard
            name="Avaneesh Khandekar"
            position="Head of Engineering"
            image={AVANEESH_ICON}
          />
          <TeamMemberCard
            name="Tanmay Saxena"
            position="Marketing Director"
            image={TANMAY_ICON}
          />
          <TeamMemberCard
            name="Shubham Singh"
            position="Lead Designer"
            image={SHUBHAM_ICON}
          />
        </div>
      </section>

      {/* Call to Action Section */}
      <section className="py-16 bg-gradient-to-r from-gray-900 to-gray-700 text-white text-center">
        <motion.h2
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          className="text-4xl font-bold"
        >
          Join Our Journey
        </motion.h2>
        <p className="mt-4 text-lg max-w-2xl mx-auto">
          Ready to be a part of something bigger? Get in touch with us today!
        </p>
        <button className="mt-6 px-6 py-3 bg-[#F58B00] text-gray-900 font-bold rounded-lg hover:bg-[#FFC67D] transition-all">
          Contact Us
        </button>
      </section>
    </div>
  );
};

export default AboutUsPage;
