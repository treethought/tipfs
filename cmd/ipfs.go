package cmd

import (
	"github.com/spf13/cobra"
)

var ipfsCmd = &cobra.Command{
	Use:   "add",
	Short: "add file",
	Run: func(cmd *cobra.Command, args []string) {
		// ipfs.ListFiles()
	},
}

func init() {
	// cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(ipfsCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.masto.yaml)")

}
