package command

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// newListCmd создаёт команду list.
func newListCmd(cl Client) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Список элементов",
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := cl.ListItems(context.Background())
			if err != nil {
				return err
			}

			if len(items) == 0 {
				fmt.Println("No items found")
				return nil
			}

			fmt.Println("ID\t\t\t\t\tTYPE\t\tCREATED")
			for _, item := range items {
				fmt.Printf("%s\t%s\t%s\n",
					item.Id,
					item.Type.String(),
					item.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
				)
			}

			return nil
		},
	}
}
