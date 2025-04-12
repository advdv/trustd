// graph-utils.ts
import {
  RandomGraphResponse,
  Node as RpcNode,
  Edge as RpcEdge,
} from "./proto/internal/rpc/v1/rpc_pb";
import { Node as RFNode, Edge as RFEdge } from "@xyflow/react";

/**
 * Convert a RandomGraphResponse from the server
 * into arrays of React Flow-compatible nodes and edges.
 */
export function convertRandomGraphResponse(response: RandomGraphResponse): {
  nodes: RFNode[];
  edges: RFEdge[];
} {
  // Convert each RPC Node to a React Flow Node
  const flowNodes: RFNode[] = response.nodes.map((node: RpcNode) => {
    return {
      id: node.id,
      type: node.type,
      position: {
        // Convert from bigint (or number) to a normal number
        x: Number(node.position?.x ?? 0),
        y: Number(node.position?.y ?? 0),
      },
      // If you store custom labels or other data, pass them here
      data: {
        label: node.data?.label ?? node.id,
      },
    };
  });

  // Convert each RPC Edge to a React Flow Edge
  const flowEdges: RFEdge[] = response.edges.map((edge: RpcEdge) => {
    return {
      id: edge.id,
      source: edge.source,
      target: edge.target,
      type: edge.type,
    };
  });

  return { nodes: flowNodes, edges: flowEdges };
}
