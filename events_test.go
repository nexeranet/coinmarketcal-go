package coinmarketcal

import (
	"context"
	"os"
	"testing"
)

func TestGetEvents(t *testing.T) {
	key := os.Getenv("KEY")
	c := NewClient(DefaultURL, key)
	ctx := context.Background()
	max := 50
	categories := "14,8,4,3,11,18,1,17,13"
	opts := EventsRequest{
		Max:        &max,
		Categories: &categories,
	}
	result, err := c.GetEvents(ctx, opts)
	if err != nil {
		t.Errorf("Error getting coins: %v", err)
	}
	// Write a string to the file
	t.Log(len(result.Body))
	for _, event := range result.Body {
		t.Log(event.Categories)
	}

}
