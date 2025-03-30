import { createFileRoute } from "@tanstack/react-router";

import "@xyflow/react/dist/style.css";
import {
  ReactFlow,
  MiniMap,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
} from "@xyflow/react";

// declare the route for this page.
export const Route = createFileRoute("/")({
  component: RouteComponent,
});

// generateRandomGraph will generate a random graph.
function generateRandomGraph(nodeCount: number, edgeCount: number) {
  const nodes = Array.from({ length: nodeCount }, (_, i) => ({
    id: (i + 1).toString(),
    position: { x: Math.random() * 1000, y: Math.random() * 1000 },
    data: { label: `Node ${(i + 1).toString()}` },
  }));

  const edges = Array.from({ length: edgeCount }, (_, i) => {
    let source, target;
    do {
      source = Math.ceil(Math.random() * nodeCount).toString();
      target = Math.ceil(Math.random() * nodeCount).toString();
    } while (source === target); // ensure source != target

    return {
      id: `e${i.toString()}-${source}-${target}`,
      source,
      target,
      label: Math.random().toFixed(4), // random decimal
    };
  });

  return { nodes, edges };
}

// declare initial nodes and edges.
const { nodes: initialNodes, edges: initialEdges } = generateRandomGraph(
  100,
  50,
);

// render the route.
function RouteComponent() {
  const [nodes, , onNodesChange] = useNodesState(initialNodes);
  const [edges, , onEdgesChange] = useEdgesState(initialEdges);

  return (
    <div style={{ width: "100vw", height: "100vh" }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
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
