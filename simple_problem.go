package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nkngn/kyber-homework/route"
)

type TradingPair struct {
	Base     string
	Quote    string
	AskPrice float64
	BidPrice float64
}

func main() {
	// read input
	base, quote, _, pairs := ReadInput("input.txt")

	fmt.Println("Base:", base, "Quote:", quote)
	for _, p := range pairs {
		fmt.Printf("%s/%s Ask: %.2f Bid: %.2f\n", p.Base, p.Quote, p.AskPrice, p.BidPrice)
	}

	// build graph
	g := BuildGraph(pairs)

	// find best ask price
	required, path, _ := g.BestAskPrice(base, quote, 1.0)
	fmt.Println(strings.Join(path, "->"))
	fmt.Printf("%.6f\n", required)

	// find best bid price
	distances, path, _ := g.BestBidPrice(base, quote, 1.0)
	fmt.Println(strings.Join(path, "->"))
	fmt.Printf("%.6f\n", distances)

	// test min heap
	// h := &TokenMinHeap{}
	// heap.Init(h)

	// heap.Push(h, TokenInfo{Token: "BTC", MinRequired: 2.5})
	// heap.Push(h, TokenInfo{Token: "ETH", MinRequired: 1.1})
	// heap.Push(h, TokenInfo{Token: "BNB", MinRequired: 1.8})

	// fmt.Println("Min token:", heap.Pop(h).(TokenInfo)) // ETH
	// fmt.Println("Next min:", heap.Pop(h).(TokenInfo))  // BNB
	// fmt.Println("Next min:", heap.Pop(h).(TokenInfo))  // BTC
}

func ReadInput(filePath string) (string, string, int, []TradingPair) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Dòng 1: base_currency quote_currency
	scanner.Scan()
	line1 := scanner.Text()
	parts := strings.Fields(line1)
	baseCurrency := parts[0]
	quoteCurrency := parts[1]

	// Dòng 2: n
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	// n dòng tiếp theo: các cặp giao dịch
	pairs := make([]TradingPair, 0, n)
	for range n {
		scanner.Scan()
		fields := strings.Fields(scanner.Text())
		ask, _ := strconv.ParseFloat(fields[2], 64)
		bid, _ := strconv.ParseFloat(fields[3], 64)
		pair := TradingPair{
			Base:     fields[0],
			Quote:    fields[1],
			AskPrice: ask,
			BidPrice: bid,
		}
		pairs = append(pairs, pair)
	}

	return baseCurrency, quoteCurrency, n, pairs
}

func BuildGraph(pairs []TradingPair) route.Graph {
	g := route.NewGraph()
	for _, p := range pairs {
		// create two edges for each trading pair
		e1 := route.SimpleEdge{
			BaseToken:  p.Base,
			QuoteToken: p.Quote,
			BidPrice:   p.BidPrice,
			AskPrice:   p.AskPrice,
		}
		g.AddEdge(e1)

		e2 := e1.GetReverseEdge()
		g.AddEdge(e2)
	}

	return g
}
