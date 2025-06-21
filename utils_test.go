package main

import (
	"testing"
	"time"
)

func TestNormalizeRating(t *testing.T) {
	tests := map[string]string{
		"Buy":            "buy",
		"Strong Buy":     "buy",
		"Outperform":     "buy",
		"Neutral":        "hold",
		"Sector Perform": "hold",
		"Underperform":   "sell",
		"Unknown Term":   "unknown",
		"market perform": "hold",
	}

	for input, expected := range tests {
		result := normalizeRating(input)
		if result != expected {
			t.Errorf("normalizeRating(%q) = %q, expected %q", input, result, expected)
		}
	}
}

func TestParsePrice(t *testing.T) {
	price, err := parsePrice("$10.50")
	if err != nil || price != 10.50 {
		t.Errorf("parsePrice failed, got %f, err %v", price, err)
	}

	price, err = parsePrice("1,200.75")
	if err == nil {
		t.Errorf("Expected error for malformed price, got %f", price)
	}
}

func TestRecomendBuy(t *testing.T) {
	r := Rating{
		TargetFrom: "$10.00",
		TargetTo:   "$11.00",
		RatingTo:   "Buy",
		Action:     "upgraded by",
		Time:       time.Now(),
	}
	if !recomendBuy(r) {
		t.Error("Expected recomendBuy to return true")
	}

	r2 := Rating{
		TargetFrom: "$10.00",
		TargetTo:   "$10.20", // only 2% increase
		RatingTo:   "Buy",
		Action:     "upgraded by",
		Time:       time.Now(),
	}
	if recomendBuy(r2) {
		t.Error("Expected recomendBuy to return false (low increase)")
	}

	r3 := Rating{
		TargetFrom: "$10.00",
		TargetTo:   "$11.00",
		RatingTo:   "Neutral",
		Action:     "upgraded by",
		Time:       time.Now(),
	}
	if recomendBuy(r3) {
		t.Error("Expected recomendBuy to return false (not buy)")
	}
}
