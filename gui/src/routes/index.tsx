import { createFileRoute } from "@tanstack/react-router";
import { GraphService } from "../proto/internal/rpc/v1/rpc_pb";

import "@xyflow/react/dist/style.css";
import { z } from "zod";
import {
  ReactFlow,
  MiniMap,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  Handle,
  Position,
  BaseEdge,
  getSmoothStepPath,
} from "@xyflow/react";
import {
  callUnaryMethod,
  createConnectQueryKey,
} from "@connectrpc/connect-query";
import { convertRandomGraphResponse } from "../graph-utils";

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

// A node that is reserved for bob
function BobNode({ data }: { data: { label: string } }) {
  return (
    <>
      <Handle type="target" position={Position.Top} />
      <div style={{ backgroundColor: "blue", padding: "1em" }}>
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

// A node that is reserved for bob
function AdaNode({ data }: { data: { label: string } }) {
  return (
    <>
      <Handle type="target" position={Position.Top} />
      <div style={{ backgroundColor: "red", padding: "1em" }}>{data.label}</div>
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

// A node that is reserved for bob
function BobWalkNode({ data }: { data: { label: string } }) {
  return (
    <>
      <Handle type="target" position={Position.Top} />
      <div style={{ backgroundColor: "blue", padding: "0.1em" }}>
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

// A node that is reserved for bob
function AdaWalkNode({ data }: { data: { label: string } }) {
  return (
    <>
      <Handle type="target" position={Position.Top} />
      <div style={{ backgroundColor: "red", padding: "0.1em" }}>
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

export function BobWalkEdge({
  sourceX,
  sourceY,
  targetX,
  targetY,
  ...props
}: {
  sourceX: number;
  sourceY: number;
  targetX: number;
  targetY: number;
}) {
  const [edgePath] = getSmoothStepPath({
    sourceX,
    sourceY,
    targetX,
    targetY,
  });

  return (
    <BaseEdge
      path={edgePath}
      {...props}
      style={{ strokeWidth: 3, stroke: "blue" }}
    />
  );
}

export function AdaWalkEdge({
  sourceX,
  sourceY,
  targetX,
  targetY,
  ...props
}: {
  sourceX: number;
  sourceY: number;
  targetX: number;
  targetY: number;
}) {
  const [edgePath] = getSmoothStepPath({
    sourceX,
    sourceY,
    targetX,
    targetY,
  });

  return (
    <BaseEdge
      path={edgePath}
      {...props}
      style={{ strokeWidth: 5, stroke: "red" }}
    />
  );
}

export function UnwalkedEdge({
  sourceX,
  sourceY,
  targetX,
  targetY,
  ...props
}: {
  sourceX: number;
  sourceY: number;
  targetX: number;
  targetY: number;
}) {
  const [edgePath] = getSmoothStepPath({
    sourceX,
    sourceY,
    targetX,
    targetY,
  });

  return (
    <BaseEdge
      path={edgePath}
      {...props}
      style={{ strokeWidth: 1, stroke: "black" }}
    />
  );
}

// custom edge types.
const edgeTypes = {
  bobWalkEdge: BobWalkEdge,
  adaWalkEdge: AdaWalkEdge,
  unwalkedEdge: UnwalkedEdge,
};

// Register custom node types
const nodeTypes = {
  labelNode: LabelNode,
  bobNode: BobNode,
  bobWalkNode: BobWalkNode,
  adaNode: AdaNode,
  adaWalkNode: AdaWalkNode,
};

// declare the route for this page.
export const Route = createFileRoute("/")({
  validateSearch: z.object({
    seed3: z.coerce.bigint(),
    seed4: z.coerce.bigint(),
  }),
  loader: ({ context: { queryClient, crpcTransport } }) => {
    return queryClient.ensureQueryData({
      staleTime: 0,
      gcTime: 0,
      queryFn: () =>
        callUnaryMethod(crpcTransport, GraphService.method.randomGraph, {
          seed1: BigInt(5),
          seed2: BigInt(3),
          numNodes: BigInt(820), // 1/10.000.000th
          initialConnected: BigInt(2),
          rewiringProbability: 0.9,

          walkLength: BigInt(50),
          numWalks: BigInt(4),

          layoutIterations: BigInt(300),
          layoutArea: 10000000,

          seed3: BigInt(5),
          seed4: BigInt(3),
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
        edgeTypes={edgeTypes}
        fitView
      >
        <Controls />
        <MiniMap />
        <Background gap={12} size={1} />
      </ReactFlow>
    </div>
  );
}
