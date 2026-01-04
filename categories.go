package coinmarketcal

import (
	"context"
)

type Category struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Etc  string `json:"etc,omitempty"`
}

func (c *Client) GetCategories(ctx context.Context) ([]Category, error) {
	response := DefaultBody[[]Category]{}
	_, err := c.GetCall(ctx, "/categories", nil, &response)
	if err != nil {
		return response.Body, err
	}
	return response.Body, nil
}
