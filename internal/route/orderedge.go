package route

type Order struct {
	Price    float64
	Quantity float64
}

type OrderEdge struct {
	BaseToken  string
	QuoteToken string
	BidOrders  []Order
	AskOrders  []Order
}

func (e OrderEdge) From() string { return e.BaseToken }
func (e OrderEdge) To() string   { return e.QuoteToken }

// SimulateSell mô phỏng việc bán amount base token qua OrderEdge này.
// Đối với OrderEdge, thực hiện walk qua bid orders xem có bán được
// hết amount hay không?
// Kết quả trả về:
//   - acquiredQuote: Trả về lượng quote token thu được nếu bán được hết
//     amount. Trả về 0 order book không đủ depth để fill hết amount.
//   - isFeasible: true nếu order book không đủ depth để fill hết amount
//     và ngược lại.
func (e OrderEdge) SimulateSell(amount float64) (float64, bool) {
	acquiredQuoteTotal := 0.0
	for _, order := range e.BidOrders {
		if order.Quantity < amount {
			acquiredQuoteTotal += order.Price * order.Quantity
			amount -= order.Quantity
		} else {
			acquiredQuoteTotal += order.Price * amount
			amount = 0
			break
		}
	}

	if amount > 0 {
		return 0.0, false
	}

	return acquiredQuoteTotal, true
}

// SimulateBuy mô phỏng việc mua amount base token qua OrderEdge này.
// Đối với OrderEdge, thực hiện walk qua ask orders để kiểm tra có đủ
// thanh khoản (liquidity) để mua hết amount hay không.
// Kết quả trả về:
//   - requiredQuote: Lượng quote token cần thiết để mua được hết amount base token.
//     Trả về 0 nếu order book không đủ depth để fill hết amount.
//   - isFeasible: true nếu order book đủ depth để fill hết amount, false nếu không.
func (e OrderEdge) SimulateBuy(amount float64) (float64, bool) {
	requiredQuoteTotal := 0.0
	for _, order := range e.AskOrders {
		if order.Quantity < amount {
			amount -= order.Quantity
			requiredQuoteTotal += order.Price * order.Quantity
		} else {
			requiredQuoteTotal += order.Price * amount
			amount = 0
			break
		}
	}

	if amount > 0 {
		return 0.0, false
	}

	return requiredQuoteTotal, true
}

// GetReverseEdge trả về một cạnh OrderEdge đảo ngược chiều giao dịch so với
// cạnh hiện tại. Các lệnh ask và bid của cạnh đảo ngược được tính toán lại:
//   - AskOrders mới được tạo từ BidOrders cũ, với giá và khối lượng đảo nghịch:
//   - Giá mới = 1 / Giá cũ
//   - Khối lượng mới = Giá cũ * Khối lượng cũ
//   - BidOrders mới được tạo từ AskOrders cũ, với công thức tương tự.
//
// Điều này đảm bảo khi đảo chiều, order book vẫn phản ánh đúng thanh khoản
// và giá trị chuyển đổi giữa hai token.
func (e OrderEdge) GetReverseEdge() Edge {
	reverseEdge := OrderEdge{
		BaseToken:  e.QuoteToken,
		QuoteToken: e.BaseToken,
		AskOrders:  make([]Order, 0, len(e.BidOrders)),
		BidOrders:  make([]Order, 0, len(e.AskOrders)),
	}

	for _, order := range e.BidOrders {
		reverseEdge.AskOrders = append(reverseEdge.AskOrders, Order{
			Price:    1.0 / order.Price,
			Quantity: order.Price * order.Quantity,
		})
	}

	for _, order := range e.AskOrders {
		reverseEdge.BidOrders = append(reverseEdge.BidOrders, Order{
			Price:    1.0 / order.Price,
			Quantity: order.Price * order.Quantity,
		})
	}

	return reverseEdge
}
