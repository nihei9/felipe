package dot

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dot",
		Short: "dot generate .dot files.",
		Long:  "dot generate .dot files.",
		RunE:  run,
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	fmt.Println("sorry, not yet implemented")

	return nil
}
