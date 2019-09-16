package cmd

import (
	"github.com/nihei9/felipe/pkg/felipe/cmd/dot"
	"github.com/nihei9/felipe/pkg/felipe/cmd/query"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "felipe",
		Short:         "felipe visualizes dependencies.",
		Long:          "felipe visualizes dependencies.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(query.NewCmd())
	cmd.AddCommand(dot.NewCmd())

	return cmd
}
