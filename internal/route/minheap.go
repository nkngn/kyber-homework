package route

type TokenInfo struct {
	Token       string
	MinRequired float64
}

// Heap (Min Heap dựa trên MinRequired)
// Sử dụng cho giải thuật Uniform Cost Search
type TokenMinHeap []TokenInfo

func (h TokenMinHeap) Len() int {
	return len(h)
}

func (h TokenMinHeap) Less(i, j int) bool {
	// Min Heap: phần tử có MinRequired nhỏ hơn sẽ lên đầu
	return h[i].MinRequired < h[j].MinRequired
}

func (h TokenMinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *TokenMinHeap) Push(x any) {
	*h = append(*h, x.(TokenInfo))
}

func (h *TokenMinHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}
