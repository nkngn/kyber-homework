package route

import (
	"testing"
)

func Test_SimulateSell(t *testing.T) {
	tests := []struct {
		name   string
		orders []Order
		amount float64
		want   float64
		wantOk bool
	}{
		{
			name: "Sell less than first order",
			orders: []Order{
				{Price: 2, Quantity: 10},
				{Price: 1.5, Quantity: 5},
			},
			amount: 5,
			want:   10, // 5 * 2
			wantOk: true,
		},
		{
			name: "Sell exactly first order",
			orders: []Order{
				{Price: 2, Quantity: 10},
				{Price: 1.5, Quantity: 5},
			},
			amount: 10,
			want:   20, // 10 * 2
			wantOk: true,
		},
		{
			name: "Sell across multiple orders",
			orders: []Order{
				{Price: 2, Quantity: 10},
				{Price: 1.5, Quantity: 5},
			},
			amount: 12,   // needs both orders
			want:   23.0, // 10 * 2 + 2 * 1.5
			wantOk: true,
		},
		{
			name: "Sell more than available",
			orders: []Order{
				{Price: 2, Quantity: 10},
				{Price: 1.5, Quantity: 5},
			},
			amount: 30, // only 15 available
			want:   0.0,
			wantOk: false,
		},
		{
			name:   "No orders",
			orders: []Order{},
			amount: 10,
			want:   0.0,
			wantOk: false,
		},
		{
			name: "Zero amount",
			orders: []Order{
				{Price: 2, Quantity: 10},
			},
			amount: 0,
			want:   0,
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := OrderEdge{
				BidOrders: tt.orders,
			}
			got, ok := edge.SimulateSell(tt.amount)
			if got != tt.want || ok != tt.wantOk {
				t.Errorf("SimulateSell(%v) = (%v, %v), want (%v, %v)", tt.amount, got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func Test_SimulateBuy(t *testing.T) {
	tests := []struct {
		name   string
		orders []Order
		amount float64
		want   float64
		wantOk bool
	}{
		{
			name: "Buy less than first order",
			orders: []Order{
				{Price: 2, Quantity: 10},
				{Price: 2.5, Quantity: 5},
			},
			amount: 5,
			want:   10, // 5 * 2
			wantOk: true,
		},
		{
			name: "Buy exactly first order",
			orders: []Order{
				{Price: 2, Quantity: 10},
				{Price: 2.5, Quantity: 5},
			},
			amount: 10,
			want:   20, // 10 * 2
			wantOk: true,
		},
		{
			name: "Buy across multiple orders",
			orders: []Order{
				{Price: 2, Quantity: 10},
				{Price: 2.5, Quantity: 5},
			},
			amount: 12,           // needs both orders
			want:   10*2 + 2*2.5, // 20 + 5 = 25
			wantOk: true,
		},
		{
			name: "Buy more than available",
			orders: []Order{
				{Price: 2, Quantity: 10},
				{Price: 2.5, Quantity: 5},
			},
			amount: 20, // only 15 available
			want:   0.0,
			wantOk: false,
		},
		{
			name:   "No orders",
			orders: []Order{},
			amount: 10,
			want:   0.0,
			wantOk: false,
		},
		{
			name: "Zero amount",
			orders: []Order{
				{Price: 2, Quantity: 10},
			},
			amount: 0,
			want:   0,
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := OrderEdge{
				AskOrders: tt.orders,
			}
			got, ok := edge.SimulateBuy(tt.amount)
			if got != tt.want || ok != tt.wantOk {
				t.Errorf("SimulateBuy(%v) = (%v, %v), want (%v, %v)", tt.amount, got, ok, tt.want, tt.wantOk)
			}
		})
	}
}
