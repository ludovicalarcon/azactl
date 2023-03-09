package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const composeFileName = "compose.yaml"

var Profile string
var Version string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a container based environment of the desired profile",
	Long: `Run a container based environment of the desired profile.
Available profile:
- go
- helm
- jekyll
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := exec.Command("docker", "compose", "version").Run()
		if err != nil {
			log.Println("docker compose is not present, please install it first")
			return err
		}

		switch Profile {
		case "jekyll":
			err = jekyllProfile()
		case "go":
			err = goProfile()
		case "helm":
			err = helmProfile()
		default:
			fmt.Println("Unknown profile! list of profile are:")
			fmt.Println("jekyll | go | helm")
		}

		if err != nil {
			cleanComposeFile(composeFileName)
		}

		return err
	},
	SilenceUsage: true,
}

func createComposeFile(fileName string, content string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(content)

	if err != nil {
		return err
	}
	return nil
}

func composeFileUp() error {

	var cmd *exec.Cmd
	cmd = exec.Command("docker", "compose", "-f", composeFileName, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	// Interractive has to be killed with ctrl+c/ctrl+d
	if err != nil && err.(*exec.ExitError).ExitCode() != 130 {
		log.Println("[ERROR] docker compose up:", err)
		return err
	}
	return nil
}

func composeFileDown() error {
	var cmd *exec.Cmd
	cmd = exec.Command("docker", "compose", "-f", composeFileName, "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("[ERROR] docker compose down:", err)
		return err
	}
	return nil
}

func runComposeFile(profile string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("docker", "compose", "-f", composeFileName, "run", "--rm", profile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()

	// Interractive has to be killed with ctrl+c/ctrl+d
	if err != nil && err.(*exec.ExitError).ExitCode() != 130 {
		log.Println("[ERROR] docker compose run", err)
		return err
	}
	return nil
}

func cleanComposeFile(fileName string) error {
	err := os.Remove(fileName)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&Profile, "profile", "p", "", "Profile tu use [jekyll/go/helm] (required)")
	runCmd.Flags().StringVarP(&Version, "version", "v", "", "Image version to use (default latest)")
	runCmd.MarkFlagRequired("profile")
}
