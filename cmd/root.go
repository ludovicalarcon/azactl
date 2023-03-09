package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "azactl",
	Short: "azactl launch a container dev environment based on profile",
	Long: `azactl launch a container dev environment based on profile using nerdctl behind the scene
- Jekyll
- Golang
- Helm`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}