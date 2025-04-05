// global css
import "./index.css";

// basic react imports
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

// tanstack query and router
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { QueryClient } from "@tanstack/react-query";

// connect rpc dependencies.
import { createConnectTransport } from "@connectrpc/connect-web";

// The transport defines what endpoint the application gets its data from.
const crpcTransport = createConnectTransport({
  baseUrl: "https://demo.connectrpc.com",
});

// tanstack query for loading data from our Go code.
const queryClient = new QueryClient();

// import the generated route tree, and create the router.
import { routeTree } from "./routeTree.gen";
const router = createRouter({
  routeTree,
  // provide context for all routes, ses: https://tanstack.com/router/v1/docs/framework/react/guide/router-context
  context: {
    queryClient,
    crpcTransport,
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

// render the react app. Notice that we DO NOT add providers for tanstack query, or connect rpc
// client. All data is loaded through route loaders.
const rootEl = document.getElementById("root");
if (rootEl !== null && rootEl.innerHTML == "") {
  createRoot(rootEl).render(
    <StrictMode>
      <RouterProvider router={router} />
    </StrictMode>,
  );
}
