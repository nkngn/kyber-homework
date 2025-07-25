package route

import (
	"container/heap"
	"errors"
	"math"
	"slices"
)

var (
	ErrNoRoute       = errors.New("no feasible route found")
	ErrArbitrageLoop = errors.New("arbitrage loop detected")
)

type Graph interface {
	AddEdge(e Edge)
	Neighbors(token string) []Edge
	BestBidPrice(base, quote string, amount float64) (float64, []string, error)
	BestAskPrice(base, quote string, amount float64) (float64, []string, error)
}

type graph struct {
	// map có key là tên token, value là các trading pairs (symbols) xuất phát
	// từ token này
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
	// hàm append tự tạo slice nếu g.edges[e.From()] truyền vào là nil, khá hay
	g.edges[e.From()] = append(g.edges[e.From()], e)
}

// Neighbors trả về danh sách các cạnh xuất phát từ token truyền vào
func (g graph) Neighbors(token string) []Edge {
	return g.edges[token]
}

// BestBidPrice tìm giá bán tốt nhất (tối đa hóa lượng quote token thu được)
// khi bán amount base token, xuất phát từ token base và kết thúc ở token quote.
//
// Kết quả trả về:
//   - float64: tỷ lệ quote/base tốt nhất (maxAcquired[quote] / amount)
//   - []string: đường đi (route) từ base đến quote
//   - err: trường hợp không tìm được đường đi hoặc xuất hiện arbitrage loop
func (g *graph) BestBidPrice(base, quote string, amount float64) (
	float64, []string, error) {
	maxAcquired, prevs, err := g.propagateBellmanFord(base, quote, amount)
	if err != nil {
		return 0, nil, err
	}

	return maxAcquired[quote] / amount, getPath(prevs, base, quote), nil
}

// BestAskPrice tìm giá mua tốt nhất (tối thiểu hóa lượng quote token cần thiết)
// để mua amount base token, xuất phát từ token base và kết thúc ở token quote.
//
// Kết quả trả về:
//   - float64: tỷ lệ quote/base tốt nhất (minRequired[quote] / amount)
//   - []string: đường đi (route) từ base đến quote
//   - err: trường hợp không tìm được đường đi hoặc xuất hiện arbitrage loop
func (g graph) BestAskPrice(base, quote string, amount float64) (
	float64, []string, error) {
	minRequired, prevs, err := g.bellmanFord(base, quote, amount)
	if err != nil {
		return 0, nil, err
	}

	path := getPath(prevs, base, quote)
	slices.Reverse(path)
	return minRequired[quote] / amount, path, nil

	// minRequired, prev, isFeasible := g.ucs(base, quote, amount)
	// if isFeasible {
	// 	path := getPath(prev, base, quote)
	// 	slices.Reverse(path)
	// 	return minRequired[quote] / amount, path, nil
	// }

	// return 0, nil, ErrNoRoute
}

// getPath sử dụng để trả về danh sách token trên đường đi từ base đến quote,
// có thể là ask route hoặc bid route. Hàm này truy vết ngược từ quote về
// base theo prevs, sau đó đảo ngược kết quả để trả về đúng thứ tự từ base
// đến quote.
// Tham số:
//   - prevs: map để truy vết đường đi tối ưu (key là đỉnh, value là đỉnh liền
//     trước)
//
// Kết quả trả về:
//   - path: slice lưu danh sách token trên đường đi từ base đến quote, bao
//     gồm cả base lẫn quote
func getPath(prevs map[string]string, base, quote string) []string {
	path := []string{}
	path = append(path, quote)
	current, ok := prevs[quote]
	for {
		if !ok {
			break
		}
		path = append(path, current)

		if current == base {
			break
		}
		current, ok = prevs[current]
	}
	slices.Reverse(path)
	return path
	// return strings.Join(path, "->")
}

// propagateBellmanFord là một biến thể của thuật toán Bellman-Ford dùng để
// lan truyền số lượng token tối đa có thể thu được tại mỗi đỉnh.
// Xuất phát từ một lượng token ban đầu ở đỉnh base, thuật toán lặp n-1 lần
// (với n là số đỉnh), mỗi lần sẽ thử "bán" toàn bộ lượng token hiện có ở mỗi
// đỉnh qua các cạnh, cập nhật số lượng token tối đa thu được ở các đỉnh kề nếu
// có thể.
// Kết quả trả về:
//   - maxAcquired: map từ tên token đến số lượng token tối đa có thể thu được
//     tại đỉnh đó
//   - prevs: map để truy vết đường đi tối ưu (key là đỉnh, value là đỉnh liền
//     trước)
//   - err: trường hợp không tìm được đường đi hoặc xuất hiện arbitrage loop
func (g *graph) propagateBellmanFord(base, quote string, amount float64) (
	map[string]float64, map[string]string, error) {
	_, ok := g.edges[base]
	if !ok {
		return nil, nil, ErrNoRoute
	}

	_, ok = g.edges[quote]
	if !ok {
		return nil, nil, ErrNoRoute
	}

	// Khởi tạo số token tối đa có thể thu được cho các đỉnh, đỉnh khởi đầu
	// bằng lượng token cần bán, các đỉnh khác bằng 0
	maxAcquired := make(map[string]float64, len(g.edges))
	for token := range g.edges {
		maxAcquired[token] = 0
	}
	maxAcquired[base] = amount

	// prevs là một map có key là đỉnh, value là đỉnh liền trước của nó
	// Dùng để xây dựng route sau này
	prevs := make(map[string]string, len(g.edges))

	// Lặp n-1 lần theo tư tưởng Bellman-Ford, với n là số đỉnh
	for range len(g.edges) - 1 {
		for baseToken, edges := range g.edges {
			if maxAcquired[baseToken] == 0 {
				continue
			}

			// Đối với mỗi cạnh, thực hiện bán thử xem có được không?
			// Nếu được thì thu về bao nhiêu quote token?
			for _, edge := range edges {
				acquiredQuote, isFeasible := edge.SimulateSell(
					maxAcquired[baseToken],
				)

				// Order book does not have enough depth to fill
				if !isFeasible {
					continue
				}

				// Cập nhật của đỉnh quote nếu bán được nhiều token hơn
				if acquiredQuote > maxAcquired[edge.To()] {
					maxAcquired[edge.To()] = acquiredQuote
					prevs[edge.To()] = edge.From()
				}
			}
		}
	}

	// Kiểm tra đỉnh nguồn có bị cập nhật không, do thuật toán khởi đầu
	// từ một lượng amount thay vì 0
	if maxAcquired[base] > amount {
		return nil, nil, ErrArbitrageLoop
	}

	// Lặp qua tất cả các cạnh một lần nữa để kiểm tra arbitrage loop
	for baseToken, edges := range g.edges {
		if maxAcquired[baseToken] == 0 {
			continue
		}

		// Đối với mỗi cạnh, thực hiện bán thử xem có được không?
		// Nếu được thì thu về bao nhiêu quote token?
		for _, edge := range edges {
			acquiredQuote, isFeasible := edge.SimulateSell(
				maxAcquired[baseToken],
			)

			// Order book does not have enough depth to fill
			if !isFeasible {
				continue
			}

			// Lượng token vẫn tăng, arbitrage loop tồn tại
			if acquiredQuote > maxAcquired[edge.To()] {
				return nil, nil, ErrArbitrageLoop
			}
		}
	}

	if maxAcquired[quote] == 0.0 {
		return nil, nil, ErrNoRoute
	}

	return maxAcquired, prevs, nil
}

// bellmanFord là một biến thể của thuật toán Bellman-Ford dùng để tìm số lượng
// quote token tối thiểu cần thiết để mua được một lượng amount base token, xuất
// phát từ đỉnh base và kết thúc ở đỉnh quote.
//
// Ý tưởng:
//   - Gán minRequired[base] = amount (lượng base token cần mua ở đỉnh xuất phát),
//     các đỉnh còn lại là +Inf.
//   - Lặp n-1 lần (với n là số đỉnh), mỗi lần thử mua base token qua các cạnh
//     (SimulateBuy), cập nhật minRequired nếu tìm được giá trị nhỏ hơn.
//   - Sau n-1 lần, lặp thêm 1 lần để kiểm tra arbitrage loop: nếu còn cập nhật
//     được minRequired thì tồn tại chu trình lợi nhuận vô hạn.
//
// Kết quả trả về:
//   - minRequired: map từ tên token đến số lượng quote token tối thiểu cần thiết
//     để mua được amount base token tại đỉnh đó.
//   - prevs: map để truy vết đường đi tối ưu (key là đỉnh, value là đỉnh liền trước).
//   - err: trả về ErrNoRoute nếu không tìm được đường đi, ErrArbitrageLoop nếu phát hiện chu trình lợi nhuận.
//
// Lưu ý: Hàm này chỉ cho kết quả hợp lý khi đồ thị không có arbitrage loop.
func (g *graph) bellmanFord(base, quote string, amount float64) (
	map[string]float64, map[string]string, error) {
	_, ok := g.edges[base]
	if !ok {
		return nil, nil, ErrNoRoute
	}

	_, ok = g.edges[quote]
	if !ok {
		return nil, nil, ErrNoRoute
	}

	// Khởi tạo số token tối đa tối thiểu để mua một lượng amount base token
	// cho các đỉnh, đỉnh khởi đầu bằng lượng token cần mua, các đỉnh khác
	// bằng 0
	minRequired := map[string]float64{}
	for base := range g.edges {
		minRequired[base] = math.Inf(1)
	}
	minRequired[base] = amount

	// prevs là một map có key là đỉnh, value là đỉnh liền trước của nó
	// Dùng để xây dựng route sau này
	prevs := make(map[string]string, len(g.edges))

	// Lặp n-1 lần theo tư tưởng Bellman-Ford, với n là số đỉnh
	for range len(g.edges) - 1 {
		for baseToken, edges := range g.edges {
			if math.IsInf(minRequired[baseToken], 1) {
				continue
			}

			// Đối với mỗi cạnh, thực hiện mua thử xem có được không?
			// Nếu được thì cần bao nhiêu quote token?
			for _, edge := range edges {
				quoteRequired, isFeasible := edge.SimulateBuy(
					minRequired[baseToken],
				)

				// Order book does not have enough depth to fill
				if !isFeasible {
					continue
				}

				// Cập nhật của đỉnh quote cần ít token hơn
				if quoteRequired < minRequired[edge.To()] {
					minRequired[edge.To()] = quoteRequired
					prevs[edge.To()] = edge.From()
				}
			}
		}
	}

	// Lặp qua tất cả các cạnh một lần nữa để kiểm tra arbitrage loop
	for baseToken, edges := range g.edges {
		if math.IsInf(minRequired[baseToken], 1) {
			continue
		}

		// Đối với mỗi cạnh, thực hiện mua thử xem có được không?
		// Nếu được thì cần bao nhiêu quote token?
		for _, edge := range edges {
			quoteRequired, isFeasible := edge.SimulateBuy(
				minRequired[baseToken],
			)

			// Order book does not have enough depth to fill
			if !isFeasible {
				continue
			}

			// Lượng token vẫn tăng, arbitrage loop tồn tại
			if quoteRequired < minRequired[edge.To()] {
				return nil, nil, ErrArbitrageLoop
			}
		}
	}

	if math.IsInf(minRequired[quote], 1) {
		return nil, nil, ErrNoRoute
	}

	return minRequired, prevs, nil
}

// ucs (Uniform Cost Search) là một biến thể của thuật toán Dijkstra dùng để
// tìm số lượng quote token tối thiểu cần thiết để mua được một lượng amount
// base token, xuất phát từ đỉnh base và kết thúc ở đỉnh quote. ucs nhanh hơn
// Dijkstra do nó sử dụng min heap để chọn ra đỉnh có min required ở mỗi lần
// lặp. Độ phức tạp O(E log(V)).
//
// Ý tưởng:
//   - Gán minRequired[base] = amount (lượng base token cần mua ở đỉnh xuất phát),
//     các đỉnh còn lại là +Inf.
//   - Mỗi bước, lấy ra đỉnh có minRequired nhỏ nhất chưa visited, giả lập việc
//     mua base token qua các cạnh (SimulateBuy). Nếu khả thi và lượng quote cần
//     nhỏ hơn giá trị hiện tại ở đỉnh kề, thì cập nhật.
//   - Mỗi đỉnh chỉ được visited một lần, đảm bảo không đi vào chu trình lợi
//     nhuận vô hạn.
//
// Kết quả trả về:
//   - minRequired: map từ tên token đến số lượng quote token tối thiểu cần thiết
//     để mua được amount base token tại đỉnh đó.
//   - prevs: map để truy vết đường đi tối ưu (key là đỉnh, value là đỉnh liền trước).
//   - err: nếu không tìm được route khả thi.
//
// Lưu ý: Hàm này chỉ cho kết quả hợp lý khi đồ thị không có arbitrage loop.
func (g graph) ucs(base, quote string, amount float64) (
	map[string]float64, map[string]string, error) {
	_, ok := g.edges[base]
	if !ok {
		return nil, nil, ErrNoRoute
	}

	_, ok = g.edges[quote]
	if !ok {
		return nil, nil, ErrNoRoute
	}

	// Khởi tạo số token tối đa tối thiểu để mua một lượng amount base token
	// cho các đỉnh, đỉnh khởi đầu bằng lượng token cần mua, các đỉnh khác
	// bằng 0
	minRequired := map[string]float64{}
	for base := range g.edges {
		minRequired[base] = math.Inf(1)
	}
	minRequired[base] = amount

	// prevs là một map có key là đỉnh, value là đỉnh liền trước của nó
	// Dùng để xây dựng route sau này
	prevs := map[string]string{}

	// Khởi tạo min heap cho thuật toán Dijkstra, để lấy ra đỉnh có số token
	// nhỏ nhất tại mỗi bước
	minHeap := &TokenMinHeap{}
	heap.Init(minHeap)
	minHeap.Push(TokenInfo{Token: base, MinRequired: amount})

	visited := make(map[string]bool, len(g.edges))

	for {
		if minHeap.Len() == 0 {
			break
		}

		// Lấy ra đỉnh có số lượng tối thiểu
		tokenInfo := minHeap.Pop()
		token := tokenInfo.(TokenInfo).Token
		required := tokenInfo.(TokenInfo).MinRequired

		if visited[token] {
			continue
		}
		visited[token] = true

		// Đối với mỗi cạnh xuất phát từ đỉnh vừa lấy, giả lập việc mua base
		// token. Nếu khả thi (order book fill đủ) và lượng quote cần nhỏ
		// hơn lượng hiện tại thì cập nhật
		for _, edge := range g.edges[token] {
			quoteRequired, feasible := edge.SimulateBuy(required)
			if feasible && quoteRequired < minRequired[edge.To()] {
				minRequired[edge.To()] = quoteRequired
				minHeap.Push(TokenInfo{
					Token: edge.To(), MinRequired: minRequired[edge.To()],
				})
				prevs[edge.To()] = edge.From()
			}
		}
	}

	if !visited[quote] {
		return nil, nil, ErrNoRoute // không có route khả thi
	}

	return minRequired, prevs, nil
}
