package route

type SimpleEdge struct {
	BaseToken  string
	QuoteToken string
	BidPrice   float64
	AskPrice   float64
}

// Implement Edge interface
func (e SimpleEdge) From() string { return e.BaseToken }
func (e SimpleEdge) To() string   { return e.QuoteToken }

// Luôn luôn bán được vì không cần kiểm tra order book có fill đủ amount hay không
// func (e SimpleEdge) IsFeasible(amount float64) bool {
// 	return true
// }

// Luôn luôn bán được vì không cần kiểm tra order book có fill đủ amount hay không
func (e SimpleEdge) SimulateSell(amount float64) (float64, bool) {
	return amount * e.BidPrice, true
}

// Luôn luôn mua được vì không cần kiểm tra order book có fill đủ amount hay không
// amount là lượng base token cần mua
// Trả về lượng quote token cần thiết để mua và có khả thi để mua theo edge này hay không?
func (e SimpleEdge) SimulateBuy(amount float64) (float64, bool) {
	return amount * e.AskPrice, true
}

func (e SimpleEdge) GetReverseEdge() Edge {
	return &SimpleEdge{
		BaseToken:  e.QuoteToken,
		QuoteToken: e.BaseToken,
		BidPrice:   1.0 / e.AskPrice,
		AskPrice:   1.0 / e.BidPrice,
	}
}
