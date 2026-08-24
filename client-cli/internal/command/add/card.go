package add

import (
	"context"

	"github.com/spf13/cobra"
)

// newCardCmd создаёт команду add card.
func newCardCmd(cl Client) *cobra.Command {
	var brand, bank, number, holder, expiry, cvv string

	cmd := &cobra.Command{
		Use:   "card",
		Short: "Добавить банковскую карту",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.AddCard(context.Background(), brand, bank, number, holder, expiry, cvv)
		},
	}

	cmd.Flags().StringVarP(&brand, "brand", "b", "", "Card brand (visa, mastercard)")
	cmd.Flags().StringVar(&bank, "bank", "", "Bank name")
	cmd.Flags().StringVarP(&number, "number", "n", "", "Card number")
	cmd.Flags().StringVar(&holder, "holder", "", "Card holder name")
	cmd.Flags().StringVarP(&expiry, "expiry", "e", "", "Expiry date (MM/YY)")
	cmd.Flags().StringVar(&cvv, "cvv", "", "CVV code")

	cobra.CheckErr(cmd.MarkFlagRequired("brand"))
	cobra.CheckErr(cmd.MarkFlagRequired("number"))
	cobra.CheckErr(cmd.MarkFlagRequired("holder"))
	cobra.CheckErr(cmd.MarkFlagRequired("expiry"))
	cobra.CheckErr(cmd.MarkFlagRequired("cvv"))

	return cmd
}
