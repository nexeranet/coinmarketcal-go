package coinmarketcal

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestGetCategories(t *testing.T) {
	key := os.Getenv("KEY")
	c := NewClient(DefaultURL, key)
	ctx := context.Background()
	result, err := c.GetCategories(ctx)
	if err != nil {
		t.Errorf("Error getting coins: %v", err)
	}
	t.Log(result)

	file, err := os.Create("categories_list.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	// Ensure the file is closed after the function completes
	defer file.Close()

	// Write a string to the file
	for _, item := range result {

		_, err = file.WriteString(
			fmt.Sprintf("Category: %d, %s, %s\n", item.ID, item.Name, item.Etc),
		)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

}
