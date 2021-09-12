package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/treethought/tipfs/ui"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pifs",
	Short: "ipfs tui",
	Run: func(cmd *cobra.Command, args []string) {
		app := ui.New()
		app.Start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.masto.yaml)")

}
