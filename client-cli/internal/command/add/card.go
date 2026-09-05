package add

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// newCardCmd создаёт команду add card.
func newCardCmd(cl Client) *cobra.Command {
	var brand, bank, number, holder, expiry string

	cmd := &cobra.Command{
		Use:   "card",
		Short: "Добавить банковскую карту",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print("CVV: ")
			cvv, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("failed to read cvv: %w", err)
			}
			fmt.Println()

			return cl.AddCard(context.Background(), brand, bank, number, holder, expiry, string(cvv))
		},
	}

	cmd.Flags().StringVarP(&brand, "brand", "b", "", "Card brand (visa, mastercard)")
	cmd.Flags().StringVar(&bank, "bank", "", "Bank name")
	cmd.Flags().StringVarP(&number, "number", "n", "", "Card number")
	cmd.Flags().StringVar(&holder, "holder", "", "Card holder name")
	cmd.Flags().StringVarP(&expiry, "expiry", "e", "", "Expiry date (MM/YY)")

	cobra.CheckErr(cmd.MarkFlagRequired("brand"))
	cobra.CheckErr(cmd.MarkFlagRequired("number"))
	cobra.CheckErr(cmd.MarkFlagRequired("holder"))
	cobra.CheckErr(cmd.MarkFlagRequired("expiry"))

	return cmd
}
