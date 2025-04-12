package rpc

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"

	"connectrpc.com/connect"
	rpcv1 "github.com/advdv/trustd/internal/rpc/v1"
)

// GenerateWattsStrogatzGraph creates a small-world network using the classic
// Watts–Strogatz model. It takes:
//   - src: an explicit random source (so you can control seeding)
//   - n: number of nodes
//   - k: each node is initially connected to k nearest neighbors (k/2 on each side in a ring)
//   - beta: rewiring probability in [0,1]
//
// It returns a RandomGraphResponse containing Nodes and Edges that match
// your protobuf definitions in this package.
//
//nolint:gocognit,varnamelen
func GenerateWattsStrogatzGraph(r *rand.Rand, n, k int, beta float64) (*rpcv1.RandomGraphResponse, string, string) {
	// adjacency[i] will be a set of neighbors of node i
	adjacency := make([]map[int]bool, n)
	for i := range adjacency {
		adjacency[i] = make(map[int]bool)
	}

	// 1. Create initial ring of edges:
	//    each node i connects to k/2 neighbors to the right (mod n).
	//    We'll store edges in adjacency to represent undirected connections.
	for i := range n {
		for j := 1; j <= k/2; j++ {
			neighbor := (i + j) % n
			adjacency[i][neighbor] = true
			adjacency[neighbor][i] = true
		}
	}

	// 2. Rewire edges with probability beta.
	//    Only consider edges where i < neighbor to avoid duplicating undirected edges.
	for i := range n {
		for j := 1; j <= k/2; j++ {
			oldNeighbor := (i + j) % n
			if i < oldNeighbor {
				if r.Float64() < beta {
					// Remove old edge
					adjacency[i][oldNeighbor] = false
					adjacency[oldNeighbor][i] = false

					// Rewire to a new neighbor that is neither i nor already a neighbor
					for {
						newNeighbor := r.IntN(n)
						if newNeighbor != i && !adjacency[i][newNeighbor] {
							adjacency[i][newNeighbor] = true
							adjacency[newNeighbor][i] = true
							break
						}
					}
				}
			}
		}
	}

	// 3. Create Nodes with positions on a circle
	//    (this is just for an example layout—positions are optional or can be changed).
	nodes := make([]*rpcv1.Node, 0, n)
	radius := 300.0
	for i := range n {
		angle := 2.0 * math.Pi * float64(i) / float64(n)
		x := int64(radius * math.Cos(angle))
		y := int64(radius * math.Sin(angle))

		node := &rpcv1.Node{}
		node.SetId(fmt.Sprintf("%d", i))
		node.SetType("labelNode")

		pos := &rpcv1.Position{}
		pos.SetX(x)
		pos.SetY(y)
		node.SetPosition(pos)

		nodes = append(nodes, node)
	}

	// ---- NEW CODE: Assign exactly one bobNode and one aliceNode randomly ----
	var bobID, adaID string
	if n >= 2 {
		bobIndex := r.IntN(n)
		adaIndex := r.IntN(n)
		for adaIndex == bobIndex {
			adaIndex = r.IntN(n)
		}
		nodes[bobIndex].SetType("bobNode")
		nodes[adaIndex].SetType("adaNode")
		bobID = nodes[bobIndex].GetId()
		adaID = nodes[adaIndex].GetId()
	}
	// -------------------------------------------------------------------------

	// 4. Convert adjacency into a list of Edges
	//    We only add an edge once (i -> j) for i < j to avoid duplicates.
	var edges []*rpcv1.Edge
	edgeCount := 0
	for i := range n {
		for j := i + 1; j < n; j++ {
			if adjacency[i][j] {
				e := &rpcv1.Edge{}
				e.SetId(fmt.Sprintf("e-%d", edgeCount))
				e.SetSource(fmt.Sprintf("%d", i))
				e.SetTarget(fmt.Sprintf("%d", j))
				edges = append(edges, e)
				edgeCount++
			}
		}
	}

	// 5. Build and return the RandomGraphResponse
	resp := &rpcv1.RandomGraphResponse{}
	resp.SetNodes(nodes)
	resp.SetEdges(edges)
	return resp, bobID, adaID
}

// ForceDirectedLayout applies a simple force-directed layout to the given RandomGraphResponse.
// It modifies and returns the same response with updated node X/Y positions.
//
//nolint:gocognit
func ForceDirectedLayout(
	rng *rand.Rand,
	iterations int,
	area float64,
	resp *rpcv1.RandomGraphResponse,
) *rpcv1.RandomGraphResponse {
	nodes := resp.GetNodes()
	edges := resp.GetEdges()

	//nolint:varnamelen
	n := len(nodes)
	if n == 0 {
		return resp
	}

	// --------------------------------------------------------------------------
	// 1) Parameters & Setup
	// --------------------------------------------------------------------------

	// Number of iterations to run. Increase if you want a more "settled" layout.
	// const iterations = 500

	// Arbitrary "area" for the layout: larger area -> more spread out.
	// In real usage, you might want to parameterize or dynamically scale this.
	// const area = 10000.0
	//nolint:varnamelen
	k := math.Sqrt(area / float64(n)) // Ideal distance between nodes (Fruchterman–Reingold concept)

	// We'll store floating-point positions internally, then cast to int64 at the end.
	positions := make([][2]float64, n)
	// We'll use displacement vectors to accumulate net force on each node per iteration.
	disp := make([][2]float64, n)

	// --------------------------------------------------------------------------
	// 2) Randomly initialize positions
	// --------------------------------------------------------------------------
	for i := range n {
		// Just pick random points in a box of size sqrt(area) x sqrt(area).
		// You could also seed them all at (0,0), but random init often helps avoid local minima.
		positions[i][0] = rng.Float64() * math.Sqrt(area)
		positions[i][1] = rng.Float64() * math.Sqrt(area)
	}

	// Helper functions for repulsive/attractive forces:
	// Fruchterman–Reingold standard:
	//   repulsiveForce(d) ~ k^2 / d
	//   attractiveForce(d) ~ d^2 / k

	repulsive := func(dist float64) float64 {
		return (k * k) / dist
	}
	attractive := func(dist float64) float64 {
		return (dist * dist) / k
	}

	// --------------------------------------------------------------------------
	// 3) Main iteration loop
	// --------------------------------------------------------------------------
	for range iterations {
		// Reset displacement
		for i := range n {
			disp[i][0] = 0
			disp[i][1] = 0
		}

		// -----------------------------
		// 3a) REPULSIVE FORCES
		// -----------------------------
		for iidx := range n {
			for jidx := iidx + 1; jidx < n; jidx++ {
				dx := positions[jidx][0] - positions[iidx][0]
				dy := positions[jidx][1] - positions[iidx][1]
				dist := math.Hypot(dx, dy)
				if dist < 1e-9 {
					// Avoid division by zero; nudge slightly
					dx = (rng.Float64() - 0.5) * 0.01
					dy = (rng.Float64() - 0.5) * 0.01
					dist = math.Hypot(dx, dy)
				}

				// Repulsive force magnitude
				force := repulsive(dist)

				// Normalize & apply
				fx := (dx / dist) * force
				fy := (dy / dist) * force

				disp[iidx][0] -= fx
				disp[iidx][1] -= fy
				disp[jidx][0] += fx
				disp[jidx][1] += fy
			}
		}

		// -----------------------------
		// 3b) ATTRACTIVE FORCES (Edges)
		// -----------------------------
		for _, e := range edges {
			srcID := e.GetSource()
			tgtID := e.GetTarget()

			// We must find the indices of the source & target nodes in "nodes".
			// For large graphs, you'd typically store a map of nodeID->index for efficiency.
			srcIndex := findNodeIndex(nodes, srcID)
			tgtIndex := findNodeIndex(nodes, tgtID)
			if srcIndex < 0 || tgtIndex < 0 {
				continue // skip if invalid
			}

			dx := positions[tgtIndex][0] - positions[srcIndex][0]
			dy := positions[tgtIndex][1] - positions[srcIndex][1]
			dist := math.Hypot(dx, dy)
			if dist < 1e-9 {
				// Avoid division by zero
				dx = (rng.Float64() - 0.5) * 0.01
				dy = (rng.Float64() - 0.5) * 0.01
				dist = math.Hypot(dx, dy)
			}

			// Attractive force magnitude
			force := attractive(dist)

			// Normalize & apply
			fx := (dx / dist) * force
			fy := (dy / dist) * force

			disp[srcIndex][0] += fx
			disp[srcIndex][1] += fy
			disp[tgtIndex][0] -= fx
			disp[tgtIndex][1] -= fy
		}

		// -----------------------------
		// 3c) Update positions
		// -----------------------------
		// Simple approach with a constant "step" or small limiting factor
		// (In real FR, you'd reduce step as 'temperature' each iteration.)
		for i := range n {
			// Move node by disp, but limit maximum movement per iteration
			dx := disp[i][0]
			dy := disp[i][1]
			maxStep := 10.0 // e.g. clamp movement so nodes don't "shoot off" too far

			dist := math.Hypot(dx, dy)
			if dist > maxStep {
				dx = dx / dist * maxStep
				dy = dy / dist * maxStep
			}
			positions[i][0] += dx
			positions[i][1] += dy
		}
	}

	// --------------------------------------------------------------------------
	// 4) Write final positions back into the response (as int64)
	// --------------------------------------------------------------------------
	for i, node := range nodes {
		xPos := int64(math.Round(positions[i][0]))
		yPos := int64(math.Round(positions[i][1]))

		// If the node doesn't have a position, we instantiate one.
		if node.GetPosition() == nil {
			node.SetPosition(&rpcv1.Position{})
		}
		node.GetPosition().SetX(xPos)
		node.GetPosition().SetY(yPos)
	}

	return resp
}

// findNodeIndex looks up the index of a node (by ID) in a slice of nodes.
// Returns -1 if not found.
func findNodeIndex(nodes []*rpcv1.Node, nodeID string) int {
	for i, n := range nodes {
		if n.GetId() == nodeID {
			return i
		}
	}
	return -1
}

// NonWeightedRandomWalk performs a random walk of `walkLength` steps starting
// from the given node ID in the provided graph, treating edges as undirected
// and picking neighbors uniformly at random.
//
// - The *source node* (where the walk starts) keeps its original type.
// - Every *other node* visited is updated to newNodeType.
// - Every *edge* traversed is updated to newEdgeType.
func NonWeightedRandomWalk(
	rng *rand.Rand,
	resp *rpcv1.RandomGraphResponse,
	walkLength int,
	startNodeID string,
	newNodeType string,
	newEdgeType string,
) []string {
	if resp == nil {
		return nil
	}
	nodes := resp.GetNodes()
	edges := resp.GetEdges()
	if len(nodes) == 0 {
		return nil
	}

	// Ensure the start node is valid. Otherwise, fallback to nodes[0].
	validStart := false
	for _, node := range nodes {
		if node.GetId() == startNodeID {
			validStart = true
			break
		}
	}
	if !validStart {
		startNodeID = nodes[0].GetId() // fallback
	}

	// Build adjacency (undirected).
	adjacency := make(map[string][]string, len(nodes))
	for _, edge := range edges {
		s := edge.GetSource()
		t := edge.GetTarget()
		adjacency[s] = append(adjacency[s], t)
		adjacency[t] = append(adjacency[t], s)
	}

	// Build quick lookups for nodes and edges.
	nodeMap := make(map[string]*rpcv1.Node, len(nodes))
	for _, nd := range nodes {
		nodeMap[nd.GetId()] = nd
	}

	// Unordered edge-key helper to ensure we treat edges as undirected.
	minMax := func(a, b string) (string, string) {
		if a < b {
			return a, b
		}
		return b, a
	}
	edgeMap := make(map[[2]string]*rpcv1.Edge, len(edges))
	for _, e := range edges {
		s, t := e.GetSource(), e.GetTarget()
		key := [2]string{}
		key[0], key[1] = minMax(s, t)
		edgeMap[key] = e
	}

	path := make([]string, 0, walkLength+1)
	current := startNodeID

	// *** DO NOT change the type of the starting node. ***
	// The user wants the starting node to keep its original style/type.
	// So we do NOT do: nodeMap[current].SetType(newNodeType)
	path = append(path, current)

	// Walk
	for i := 0; i < walkLength; i++ {
		neighbors := adjacency[current]
		if len(neighbors) == 0 {
			break
		}
		next := neighbors[rng.IntN(len(neighbors))]

		// -- 1) Update the node's type (for the newly visited node) --
		nodeMap[next].SetType(newNodeType)

		// -- 2) Update the edge's type for the walked edge --
		s, t := minMax(current, next)
		if ePtr, ok := edgeMap[[2]string{s, t}]; ok {
			ePtr.SetType(newEdgeType)
		}

		path = append(path, next)
		current = next
	}

	return path
}

func (g) RandomGraph(
	_ context.Context, req *connect.Request[rpcv1.RandomGraphRequest],
) (*connect.Response[rpcv1.RandomGraphResponse], error) {
	//nolint:gosec
	graphRng := rand.New(rand.NewPCG(
		req.Msg.GetSeed1(), req.Msg.GetSeed2(),
	))
	//nolint:gosec
	walkRng := rand.New(rand.NewPCG(
		req.Msg.GetSeed3(), req.Msg.GetSeed4(),
	))

	_ = walkRng

	graph, bobID, adaID := GenerateWattsStrogatzGraph(graphRng,
		int(req.Msg.GetNumNodes()),
		int(req.Msg.GetInitialConnected()),
		req.Msg.GetRewiringProbability())

	graph = ForceDirectedLayout(graphRng,
		int(req.Msg.GetLayoutIterations()), req.Msg.GetLayoutArea(), graph)

	// @TODO figure out why the starting node does't keep its original type
	// @TODO make sure the walk edges use the same bezier edges, or make the default smooth edgeagain.
	// @TODO Get multiple random walks working
	// @TODO make sure the start/end nodes keep their original (non walked) style
	NonWeightedRandomWalk(walkRng, graph, int(req.Msg.GetWalkLength()), bobID, "bobWalkNode", "bobWalkEdge")
	NonWeightedRandomWalk(walkRng, graph, int(req.Msg.GetWalkLength()), adaID, "adaWalkNode", "adaWalkEdge")

	// set the type to a base edge if it's not walked.
	for _, edge := range graph.GetEdges() {
		if edge.GetType() == "" {
			edge.SetType("unwalkedEdge")
		}
	}

	return connect.NewResponse(graph), nil
}
