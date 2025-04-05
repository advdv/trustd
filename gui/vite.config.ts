import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [TanStackRouterVite({}), react(), tailwindcss()],
  build: {
    rollupOptions: {
      output: {
        // to reduce the size of the largest chunk we split out react.
        manualChunks(id) {
          if (id.includes("node_modules")) {
            if (id.includes("react") && !id.includes("@tanstack")) {
              return "react";
            }
          }
        },
      },
    },
  },
});
