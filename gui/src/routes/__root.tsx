import * as React from "react";
import { Outlet, createRootRouteWithContext } from "@tanstack/react-router";
import { QueryClient } from "@tanstack/react-query";
import { Transport } from "@connectrpc/connect";

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
  crpcTransport: Transport;
}>()({
  component: RootComponent,
});

function RootComponent() {
  return (
    <React.Fragment>
      <Outlet />
    </React.Fragment>
  );
}
