import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import { RouterProvider, createRouter } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

// tanstack query for loading data from our Go code.
const queryClient = new QueryClient();

// import the generated route tree, and create the router.
import { routeTree } from "./routeTree.gen";
const router = createRouter({
  routeTree,
  // provide context for all routes, ses: https://tanstack.com/router/v1/docs/framework/react/guide/router-context
  context: {
    queryClient,
  },
  // preload data already when user hovers over a link.
  defaultPreload: "intent",
  // Since we're using React Query, we don't want loader calls to ever be stale
  // This will ensure that the loader is always called when the route is preloaded or visited
  defaultPreloadStaleTime: 0,
  // see: https://tanstack.com/router/v1/docs/framework/react/guide/scroll-restoration
  scrollRestoration: true,
});

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
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </StrictMode>,
  );
}
