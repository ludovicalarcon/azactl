package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "azactl",
	Short: "azactl launch a container dev environment based on profile",
	Long: `azactl launch a container dev environment based on profile using nerdctl behind the scene
- Jekyll
- Golang`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
	err := exec.Command("docker", "compose", "version").Run()
	if err != nil {
		log.Fatalln("nerdctl is not present, please install it first")
	}
}
