package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
)

const composeFileName = "compose.yaml"

var Profile string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a container based environment of the desired profile",
	Long:  "run a container based environment of the desired profile",
	RunE: func(cmd *cobra.Command, args []string) error {

		var err error = nil

		switch Profile {
		case "jekyll":
			err = jekyllProfile()
		case "go":
			err = goProfile()
		default:
			fmt.Println("Unknown profile! list of profile are:")
			fmt.Println("jekyll | go")
		}

		if err != nil {
			cleanComposeFile(composeFileName)
		}

		return err
	},
}

func goProfile() error {
	const compose = `version: '3.3'
services:
  go:
    container_name: go
    image: azalax/golang:1.20
    volumes:
      - '$PWD:/home/go/workspace'
    tty: true
    stdin_open: true
`
	if err := createComposeFile(composeFileName, compose); err != nil {
		return err
	}
	if err := runComposeFile(Profile); err != nil {
		return err
	}
	return cleanComposeFile(composeFileName)
}

func jekyllProfile() error {
	const compose = `version: '3.3'
services:
  jekyll:
    container_name: jekyll
    image: jekyll/jekyll
    volumes:
      - '$PWD:/workspace'
    ports:
      - '127.0.0.1:8282:4000/tcp'
    working_dir: /workspace
    command: sh -c "bundle install && bundle exec jekyll server --host 0.0.0.0"
    tty: true
    stdin_open: true
`
	var wg sync.WaitGroup
	wg.Add(1) // add a goroutine to wait

	// Jekyll server has to be killed with ctrl+c
	// catch ctrl+c and clean
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		defer wg.Done()

		composeFileDown()
		cleanComposeFile(composeFileName)
	}()

	if err := createComposeFile(composeFileName, compose); err != nil {
		return err
	}
	if err := composeFileUp(); err != nil {
		return err
	}

	wg.Wait() // wait goroutine to complete
	return nil
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

	runCmd.Flags().StringVarP(&Profile, "profile", "p", "", "Profile tu use [jekyll/go] (required)")
	runCmd.MarkFlagRequired("profile")
}
