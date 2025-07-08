package route

type Edge interface {
	From() string
	To() string

	// SimulateSell mô phỏng việc bán amount base token qua cạnh này.
	SimulateSell(amount float64) (float64, bool)

	// SimulateBuy mô phỏng việc mua amount base token qua cạnh này.
	SimulateBuy(amount float64) (float64, bool)

	// GetReverseEdge trả về một cạnh đảo ngược chiều giao dịch so với
	// cạnh hiện tại.
	GetReverseEdge() Edge
}
