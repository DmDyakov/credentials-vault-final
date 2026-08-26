package validate

import (
	"fmt"
	"strings"

	creditcard "github.com/durango/go-credit-card"
)

// Card проверяет данные банковской карты.
func Card(brand, bank, number, holder, expiry, cvv string) error {
	if brand == "" {
		return fmt.Errorf("brand is required")
	}

	if bank == "" {
		return fmt.Errorf("bank is required")
	}

	if holder == "" {
		return fmt.Errorf("holder is required")
	}

	parts := strings.Split(expiry, "/")
	if len(parts) != 2 {
		return fmt.Errorf("expiry must be in MM/YY format")
	}

	card := creditcard.Card{
		Number: number,
		Cvv:    cvv,
		Month:  parts[0],
		Year:   "20" + parts[1],
	}

	if err := card.Validate(); err != nil {
		return fmt.Errorf("invalid card: %w", err)
	}

	return nil
}
