package route

type SimpleEdge struct {
	BaseToken  string
	QuoteToken string
	BidPrice   float64
	AskPrice   float64
}

func (e SimpleEdge) From() string { return e.BaseToken }
func (e SimpleEdge) To() string   { return e.QuoteToken }

// SimulateSell mô phỏng việc bán amount base token qua SimpleEdge này.
// Đối với SimpleEdge, giả định thanh khoản (liquidity) là vô hạn nên luôn
// bán được bất kỳ amount nào, không cần kiểm tra order book.
// Trả về lượng quote token thu được và luôn trả về true.
func (e SimpleEdge) SimulateSell(amount float64) (float64, bool) {
	return amount * e.BidPrice, true
}

// SimulateBuy mô phỏng việc mua amount base token qua SimpleEdge này.
// Đối với SimpleEdge, giả định thanh khoản (liquidity) là vô hạn nên luôn
// mua được bất kỳ amount nào, không cần kiểm tra order book.
// Trả về lượng quote token cần thiết để mua amount base token và luôn trả về true.
func (e SimpleEdge) SimulateBuy(amount float64) (float64, bool) {
	return amount * e.AskPrice, true
}

// GetReverseEdge trả về một cạnh SimpleEdge đảo ngược chiều giao dịch so với
// cạnh hiện tại. Giá Bid/Ask của cạnh đảo ngược sẽ là nghịch đảo của Ask/Bid
// của cạnh gốc.
// Ví dụ: Nếu cạnh gốc là A->B với BidPrice/AskPrice thì cạnh đảo ngược là
// B->A với BidPrice = 1/AskPrice, AskPrice = 1/BidPrice.
func (e SimpleEdge) GetReverseEdge() Edge {
	return &SimpleEdge{
		BaseToken:  e.QuoteToken,
		QuoteToken: e.BaseToken,
		BidPrice:   1.0 / e.AskPrice,
		AskPrice:   1.0 / e.BidPrice,
	}
}
