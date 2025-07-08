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
	base, quote, amount, graph, err := ReadExpandedInput("test/expanded_input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi đọc file expanded_input.txt: %v\n", err)
		os.Exit(1)
	}

	// find best ask price
	bestAskPrice, bestAskRoute, isFeasible := graph.BestAskPrice(base, quote, amount)
	if isFeasible {
		fmt.Println(strings.Join(bestAskRoute, "->"))
		fmt.Printf("%.6f\n", bestAskPrice)
	} else {
		fmt.Printf("%s->%s ask route is not feasible\n", quote, base)
	}

	// find best bid price
	bestBidPrice, bestBidRoute, isFeasible := graph.BestBidPrice(base, quote, amount)
	if isFeasible {
		fmt.Println(strings.Join(bestBidRoute, "->"))
		fmt.Printf("%.6f\n", bestBidPrice)
	} else {
		fmt.Printf("%s->%s bid route is not feasible\n", quote, base)
	}
}

func ReadExpandedInput(filePath string) (string, string, float64, route.Graph, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", 0, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Dòng 1: base_currency quote_currency amount
	// KNC ETH 100
	scanner.Scan()
	line1 := scanner.Text()
	parts := strings.Fields(line1)
	baseCurrency := parts[0]
	quoteCurrency := parts[1]
	amount, _ := strconv.ParseFloat(parts[2], 64)

	// Dòng 2: n - số cặp giao dịch
	// 2
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	// n block tiếp theo: mỗi block thể hiện một cặp giao dịch, đi cùng order book tương ứng
	// KNC USDT
	// 2
	// 1.1 150
	// 1.2 200
	// 2
	// 0.9 100
	// 0.8 300
	edges := make([]route.Edge, 0, n*2)
	for range n {
		// Tên cặp
		scanner.Scan()
		pairLine := scanner.Text()
		pairFields := strings.Fields(pairLine)
		base := pairFields[0]
		quote := pairFields[1]

		// Số lượng ask orders
		scanner.Scan()
		nAsk, _ := strconv.Atoi(scanner.Text())
		askOrders := make([]route.Order, 0, nAsk)
		for range nAsk {
			scanner.Scan()
			askFields := strings.Fields(scanner.Text())
			price, _ := strconv.ParseFloat(askFields[0], 64)
			qty, _ := strconv.ParseFloat(askFields[1], 64)
			askOrders = append(askOrders, route.Order{Price: price, Quantity: qty})
		}

		// Số lượng bid orders
		scanner.Scan()
		nBid, _ := strconv.Atoi(scanner.Text())
		bidOrders := make([]route.Order, 0, nBid)
		for range nBid {
			scanner.Scan()
			bidFields := strings.Fields(scanner.Text())
			price, _ := strconv.ParseFloat(bidFields[0], 64)
			qty, _ := strconv.ParseFloat(bidFields[1], 64)
			bidOrders = append(bidOrders, route.Order{Price: price, Quantity: qty})
		}

		// Tạo OrderEdge và reverse edge
		edge := route.OrderEdge{
			BaseToken:  base,
			QuoteToken: quote,
			AskOrders:  askOrders,
			BidOrders:  bidOrders,
		}
		edges = append(edges, edge)
		edges = append(edges, edge.GetReverseEdge())
	}

	g := route.NewGraphWithEdges(edges)
	return baseCurrency, quoteCurrency, amount, g, nil
}
