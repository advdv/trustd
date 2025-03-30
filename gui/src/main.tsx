import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider, createRouter } from "@tanstack/react-router";

// import the generated route tree, and create the router.
import { routeTree } from "./routeTree.gen";
const router = createRouter({ routeTree });

// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

// global css
import "./index.css";

// render the react app.
const rootEl = document.getElementById("root");
if (rootEl !== null && rootEl.innerHTML == "") {
  createRoot(rootEl).render(
    <StrictMode>
      <RouterProvider router={router} />
    </StrictMode>,
  );
}
