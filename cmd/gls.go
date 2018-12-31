package cmd

import (
	"github.com/spf13/cobra"
	"github.com/JajaDoc/g-explorer/ui"
)

func doGlsCmd(cmd *cobra.Command, args []string) {
	ui.Run()
}
