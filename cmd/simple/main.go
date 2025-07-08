package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nkngn/kyber-homework/internal/route"
)

func main() {
	// read input from file, build graph
	base, quote, graph, err := ReadSimpleInput("test/input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi đọc file input.txt: %v\n", err)
		os.Exit(1)
	}

	// Đối với simple problem, lượng base token cần bán/mua luôn là 1 đơn vị
	// find best ask price
	bestAskPrice, bestAskRoute, _ := graph.BestAskPrice(base, quote, 1.0)
	fmt.Println(strings.Join(bestAskRoute, "->"))
	fmt.Printf("%.6f\n", bestAskPrice)

	// find best bid price
	bestBidPrice, bestBidRoute, _ := graph.BestBidPrice(base, quote, 1.0)
	fmt.Println(strings.Join(bestBidRoute, "->"))
	fmt.Printf("%.6f\n", bestBidPrice)
}

func ReadSimpleInput(filePath string) (string, string, route.Graph, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Dòng 1: base_currency quote_currency
	// KNC ETH
	scanner.Scan()
	line1 := scanner.Text()
	parts := strings.Fields(line1)
	baseCurrency := parts[0]
	quoteCurrency := parts[1]

	// Dòng 2: n - số cặp giao dịch
	// 2
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	// n dòng tiếp theo: các cặp giao dịch, mỗi cặp tương ứng hai cạnh trong đồ thị
	// KNC USDT 1.1 0.9
	// ETH USDT 360 355
	edges := make([]route.Edge, 0, n*2)
	for range n {
		scanner.Scan()
		fields := strings.Fields(scanner.Text())
		ask, _ := strconv.ParseFloat(fields[2], 64)
		bid, _ := strconv.ParseFloat(fields[3], 64)
		edge := route.SimpleEdge{
			BaseToken:  fields[0],
			QuoteToken: fields[1],
			AskPrice:   ask,
			BidPrice:   bid,
		}
		edges = append(edges, edge)
		edges = append(edges, edge.GetReverseEdge())
	}

	// Tạo đồ thị từ danh sách các cạnh
	g := route.NewGraphWithEdges(edges)

	return baseCurrency, quoteCurrency, g, nil
}
