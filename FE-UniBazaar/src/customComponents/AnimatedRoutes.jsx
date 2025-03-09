import React, { Suspense } from "react";
import { Routes, Route, useLocation } from "react-router-dom";
import { AnimatePresence } from "framer-motion";
import Spinner from "@/customComponents/Spinner";
import PageWrapper from "@/customComponents/PageWrapper";
import SellProductPage from "@/pages/SellProductPage";
import AboutUsPage from "@/pages/AboutUsPage";
import LandingPage from "@/pages/LandingPage";

function AnimatedRoutes() {
  const location = useLocation();

  return (
    <AnimatePresence mode="wait">
      <Routes location={location} key={location.pathname}>
        <Route path="/" element={<LandingPage />} />
        <Route
          path="/sell"
          element={
            <Suspense fallback={<Spinner />}>
              <PageWrapper>
                <SellProductPage />
              </PageWrapper>
            </Suspense>
          }
        />
        <Route
          path="/about"
          element={
            <PageWrapper direction="right">
              <AboutUsPage />
            </PageWrapper>
          }
        />
      </Routes>
    </AnimatePresence>
  );
}

export default AnimatedRoutes;
