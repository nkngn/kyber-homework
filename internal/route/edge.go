package route

type Edge interface {
	From() string
	To() string
	SimulateSell(amount float64) (float64, bool)
	SimulateBuy(amount float64) (float64, bool)
	GetReverseEdge() Edge
}
