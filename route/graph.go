package route

import (
	"container/heap"
	"math"
	"slices"
)

type Graph interface {
	AddEdge(e Edge)
	Neighbors(token string) []Edge
	BestBidPrice(base, quote string, amount float64) (float64, []string, bool)
	BestAskPrice(base, quote string, amount float64) (float64, []string, bool)
}

type graph struct {
	edges map[string][]Edge
}

func NewGraph() Graph {
	return &graph{
		edges: make(map[string][]Edge),
	}
}

func NewGraphWithEdges(edgeList []Edge) Graph {
	g := &graph{edges: make(map[string][]Edge)}
	for _, e := range edgeList {
		g.AddEdge(e)
	}
	return g
}

func (g *graph) AddEdge(e Edge) {
	// _, ok := g.edges[e.From()]
	// if ok {
	// 	g.edges[e.From()] = append(g.edges[e.From()], e)
	// } else {
	// 	g.edges[e.From()] = []Edge{e}
	// }
	g.edges[e.From()] = append(g.edges[e.From()], e)
}

func (g graph) Neighbors(token string) []Edge {
	return g.edges[token]
}

func (g *graph) BestBidPrice(base, quote string, amount float64) (float64, []string, bool) {
	distances, predecessors := g.bellmanFord(base, amount)
	return distances[quote], GetPath(predecessors, base, quote), true
}

func (g *graph) bellmanFord(source string, amount float64) (map[string]float64, map[string]string) {
	dist := make(map[string]float64) // distances
	for token := range g.edges {
		dist[token] = 0
	}
	dist[source] = amount
	predecessors := make(map[string]string)

	for i := 0; i < len(g.edges)-1; i++ {
		for u, edges := range g.edges {
			if dist[u] == 0 {
				continue
			}
			for _, edge := range edges {
				// amtOut := edge.BidPrice * dist[u]
				amtOut, isFeasible := edge.SimulateSell(dist[u])

				if !isFeasible {
					continue
				}

				if amtOut > dist[edge.To()] {
					dist[edge.To()] = amtOut
					predecessors[edge.To()] = edge.From()
				}
			}
		}
	}

	return dist, predecessors
}

func GetPath(predecessors map[string]string, base, quote string) []string {
	path := []string{}
	path = append(path, quote)
	current, ok := predecessors[quote]
	for {
		if !ok {
			break
		}
		path = append(path, current)

		if current == base {
			break
		}
		current, ok = predecessors[current]
	}
	slices.Reverse(path)
	return path
	// return strings.Join(path, "->")
}

func (g graph) BestAskPrice(base, quote string, amount float64) (float64, []string, bool) {
	minRequired, prev := g.uniCostSearch(base, quote, 1.0)
	return minRequired[quote], GetPath(prev, quote, base), true
	// fmt.Println(GetPath(prev, quote, base))
	// fmt.Printf("%.6f\n", minRequired[quote])
}

func (g graph) uniCostSearch(source, dest string, amount float64) (map[string]float64, map[string]string) {
	minRequired := map[string]float64{}
	for base := range g.edges {
		minRequired[base] = math.Inf(1)
	}
	minRequired[source] = amount        // minReq[u] = lượng base token cần
	predecessors := map[string]string{} // để reconstruct path

	minHeap := &TokenMinHeap{}
	heap.Init(minHeap)
	minHeap.Push(TokenInfo{Token: source, MinRequired: amount})

	visited := make(map[string]bool, len(g.edges))

	for {
		if minHeap.Len() == 0 {
			break
		}

		tokenInfo := minHeap.Pop()
		token := tokenInfo.(TokenInfo).Token
		required := tokenInfo.(TokenInfo).MinRequired

		if visited[token] {
			continue
		}
		visited[token] = true

		for _, edge := range g.edges[token] {
			quoteRequired, feasible := edge.SimulateBuy(required)
			// if required*edge.AskPrice < minRequired[edge.To()] {
			if feasible && quoteRequired < minRequired[edge.To()] {
				// minRequired[edge.To()] = required * edge.AskPrice
				minRequired[edge.To()] = quoteRequired
				minHeap.Push(TokenInfo{Token: edge.To(), MinRequired: minRequired[edge.To()]})
				predecessors[edge.From()] = edge.To()
			}
		}
	}

	if !visited[dest] {
		return nil, nil // không có route khả thi
	}

	return minRequired, predecessors
}
