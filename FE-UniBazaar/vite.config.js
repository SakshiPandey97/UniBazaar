import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

const isDocker = process.env.DOCKER === "true";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": "/src",
    },
  },
  server: {
    allowedHosts: ["unibazaar-4cjp.onrender.com", "localhost"],
    port: 3000,
    open: !isDocker,
    host: "0.0.0.0",
  },
  test: {
    // Add this section
    globals: true, // Enables global API like `describe`, `it`
    environment: "jsdom", // Use JSDOM to simulate the browser environment
    setupFiles: "./src/__tests__/setupTests.js", // Global setup file
  }
});
