package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golang-api",
	Short: "A brief description of your application",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
func init() {
	cobra.OnInitialize()
}
