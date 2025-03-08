import React, { lazy, Suspense }  from "react";
import Banner from "@/customComponents/Banner";
import Spinner from "@/customComponents/Spinner";
const Products = lazy(() => import("@/customComponents/Products"));

function LandingPage() {
  return (
    <div className="bg-[#D6D2D2]">
      <Banner />
      <Suspense fallback={<Spinner />}>
        <Products />
      </Suspense>
    </div>
  );
}

export default LandingPage;
