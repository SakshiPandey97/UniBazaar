import React, { Suspense } from "react";
import { Routes, Route, useLocation } from "react-router-dom";
import { AnimatePresence } from "framer-motion";
import Spinner from "@/customComponents/Spinner";
import PageWrapper from "@/customComponents/PageWrapper";
import SellProductPage from "@/pages/SellProductPage";
import Messaging from "@/pages/Messaging";
import AboutUsPage from "@/pages/AboutUsPage";
import ProductsPage from "@/pages/ProductsPage";
import LandingPage from "@/pages/LandingPage";
import MyProductsPage from "@/pages/MyProductsPage";
import PrivateRoute from "@/customComponents/PrivateRoute";

function AnimatedRoutes() {
  const location = useLocation();

  return (
    <AnimatePresence mode="wait">
      <Routes location={location} key={location.pathname}>
        <Route path="/" element={<LandingPage />} />

        <Route
          path="/sell"
          element={
            <PrivateRoute>
              <Suspense fallback={<Spinner />}>
                <PageWrapper direction="right">
                  <SellProductPage />
                </PageWrapper>
              </Suspense>
            </PrivateRoute>
          }
        />

        <Route
          path="/messaging"
          element={
            <PrivateRoute>
              <Suspense fallback={<Spinner />}>
                <PageWrapper direction="right">
                  <Messaging />
                </PageWrapper>
              </Suspense>
            </PrivateRoute>
          }
        />

        <Route
          path="/products"
          element={
            <PageWrapper direction="right">
              <Suspense fallback={<Spinner />}>
                <ProductsPage />
              </Suspense>
            </PageWrapper>
          }
        />
        <Route
          path="/userproducts"
          element={
            <PrivateRoute>
              <Suspense fallback={<Spinner />}>
                <PageWrapper direction="right">
                  <MyProductsPage />
                </PageWrapper>
              </Suspense>
            </PrivateRoute>
          }
        />
        <Route
          path="/about"
          element={
            <Suspense fallback={<Spinner />}>
              <PageWrapper direction="right">
                <AboutUsPage />
              </PageWrapper>
            </Suspense>
          }
        />
      </Routes>
    </AnimatePresence>
  );
}

export default AnimatedRoutes;
