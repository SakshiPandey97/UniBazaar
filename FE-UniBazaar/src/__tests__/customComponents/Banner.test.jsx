import { render, screen } from "@testing-library/react";
import Banner from "../../customComponents/Banner";
import { describe, it, expect } from "vitest";

describe("Banner Component", () => {
  it("renders the banner text", () => {
    render(<Banner />);
    expect(screen.getByText("Uni")).toBeInTheDocument();
    expect(screen.getByText("Bazaar")).toBeInTheDocument();
    expect(screen.getByText("Connecting students for buying/selling")).toBeInTheDocument();
  });
});