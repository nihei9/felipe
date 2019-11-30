package main

import (
	"os"

	"github.com/nihei9/felipe/cmd/felipe/dot"
	"github.com/nihei9/felipe/cmd/felipe/query"
	"github.com/spf13/cobra"
)

func main() {
	os.Exit(doMain())
}

func doMain() int {
	cmd := newCmd()
	cmd.SetOutput(os.Stdout)
	err := cmd.Execute()
	if err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		return 1
	}

	return 0
}

func newCmd() *cobra.Command {
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
