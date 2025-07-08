package route

import (
	"container/heap"
	"testing"
)

func TestTokenMinHeap_PushAndPop(t *testing.T) {
	h := &TokenMinHeap{}
	heap.Init(h)

	tokens := []TokenInfo{
		{"A", 5.0},
		{"B", 2.0},
		{"C", 8.0},
		{"D", 1.0},
	}

	for _, token := range tokens {
		heap.Push(h, token)
	}

	expectedOrder := []TokenInfo{
		{"D", 1.0},
		{"B", 2.0},
		{"A", 5.0},
		{"C", 8.0},
	}

	for i, expected := range expectedOrder {
		item := heap.Pop(h).(TokenInfo)
		if item != expected {
			t.Errorf("Pop %d: got %+v, want %+v", i, item, expected)
		}
	}
}

func TestTokenMinHeap_Len(t *testing.T) {
	h := &TokenMinHeap{}
	heap.Init(h)
	if h.Len() != 0 {
		t.Errorf("Expected length 0, got %d", h.Len())
	}
	heap.Push(h, TokenInfo{"A", 1.0})
	heap.Push(h, TokenInfo{"B", 2.0})
	if h.Len() != 2 {
		t.Errorf("Expected length 2, got %d", h.Len())
	}
}
