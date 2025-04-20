import { render, screen } from "@testing-library/react";
import { SearchProvider } from "../../context/SearchContext";
import { MemoryRouter } from "react-router-dom";
import Banner from "../../customComponents/Banner";
import { describe, it, expect } from "vitest";

describe("Banner Component", () => {
  it("renders the banner text", () => {
    render(
      <MemoryRouter>
        <SearchProvider>
          <Banner />
        </SearchProvider>
      </MemoryRouter>
    );

    expect(screen.getByText("Uni")).toBeInTheDocument();
    expect(screen.getByText("Bazaar")).toBeInTheDocument();
    expect(screen.getByText("Declutter, Discover, Deal")).toBeInTheDocument();
  });
});
