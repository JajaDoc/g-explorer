package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/JajaDoc/g-explorer/utils"
)

var cfgFile string
var exportFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gls",
	Short: "TUI Directory Explorer",
	Long: `This tool is directory explorer based text user interface.`,
	Args: cobra.MaximumNArgs(0),
	Run:  doGlsCmd,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		utils.Exit(1)
	}
	utils.Cleanup()
}
