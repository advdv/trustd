import { createFileRoute } from "@tanstack/react-router";
import { GraphService } from "../proto/internal/rpc/v1/rpc_pb";

import "@xyflow/react/dist/style.css";
import {
  ReactFlow,
  MiniMap,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  Handle,
  Position,
} from "@xyflow/react";
import {
  callUnaryMethod,
  createConnectQueryKey,
} from "@connectrpc/connect-query";
import { convertRandomGraphResponse } from "../graph-utils";

// declare the route for this page.
export const Route = createFileRoute("/")({
  loader: ({ context: { queryClient, crpcTransport } }) => {
    return queryClient.ensureQueryData({
      queryFn: () =>
        callUnaryMethod(crpcTransport, GraphService.method.randomGraph, {
          seed1: BigInt(1),
          seed2: BigInt(30),
          numNodes: BigInt(820), // 1/10.000.000th
          initialConnected: BigInt(2),
          rewiringProbability: 0.9,

          layoutIterations: BigInt(300),
          layoutArea: 10000000,
        }),
      queryKey: createConnectQueryKey({
        transport: crpcTransport,
        schema: GraphService.method.randomGraph,
        cardinality: "finite",
      }),
    });
  },
  component: RouteComponent,
});

// A minimal custom node that only shows text
function LabelNode({ data }: { data: { label: string } }) {
  return (
    <>
      <Handle type="target" position={Position.Top} />
      <div style={{ background: "transparent", border: "none" }}>
        {data.label}
      </div>
      <Handle type="source" position={Position.Bottom} id="a" />
      <Handle
        type="source"
        position={Position.Bottom}
        id="b"
        style={{ left: 10 }}
      />
    </>
  );
}

// Register custom node types
const nodeTypes = { labelNode: LabelNode };

// render the route.
function RouteComponent() {
  const nodesAndEdges = Route.useLoaderData();
  const { nodes: initialNodes, edges: initialEdges } =
    convertRandomGraphResponse(nodesAndEdges);

  const [nodes, , onNodesChange] = useNodesState(initialNodes);
  const [edges, , onEdgesChange] = useEdgesState(initialEdges);

  return (
    <div style={{ width: "100vw", height: "100vh" }}>
      <ReactFlow
        minZoom={0.001}
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        fitView
        defaultEdgeOptions={{ type: "step" }}
      >
        <Controls />
        <MiniMap />
        <Background gap={12} size={1} />
      </ReactFlow>
    </div>
  );
}
