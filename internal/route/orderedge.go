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

// Bán một lượng `amount` base token, trả về lượng quote token thu được,
// trường hợp order book không đủ thì trả về false tương đương với không bán được
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

// Mua một lượng `amount` base token, trả về lượng quote token cần thiết,
// trường hợp order book không đủ lượng base token để mua thì trả về false
// tương đương với không mua được
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
