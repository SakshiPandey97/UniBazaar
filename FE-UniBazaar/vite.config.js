import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// Detect if running inside Docker
const isDocker = process.env.DOCKER === "true";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": "/src",
    },
  },
  server: {
    allowedHosts: [
      "unibazaar-4cjp.onrender.com",
      "localhost", // add this if you need to support local development as well
    ],
    port: 3000,
    open: !isDocker, // Open in local dev, not in Docker
    host: "0.0.0.0", // Ensures Docker works properly
  },
  test: {    // Add this section
    globals: true,        // Enables global API like `describe`, `it`
    environment: "jsdom", // Use JSDOM to simulate the browser environment
    setupFiles: "./src/__tests__/setupTests.js", // Global setup file
  },
});
