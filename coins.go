package coinmarketcal

import (
	"context"
)

type Coin struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Symbol      string `json:"symbol,omitempty"`
	Fullname    string `json:"fullname,omitempty"`
	Upcoming    int    `json:"upcoming,omitempty"`
	Popular     int    `json:"popular,omitempty"`
	Trending    int    `json:"trending,omitempty"`
	Influential int    `json:"influential,omitempty"`
	Catalyst    int    `json:"catalyst,omitempty"`
	Etc         string `json:"etc,omitempty"`
}

func (c *Client) GetCoins(ctx context.Context) ([]Coin, error) {
	response := DefaultBody[[]Coin]{}
	_, err := c.GetCall(ctx, "/coins", nil, &response)
	if err != nil {
		return response.Body, err
	}
	return response.Body, nil
}
