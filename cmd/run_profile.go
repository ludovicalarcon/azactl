package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ProfileInfo struct {
	name     string
	image    string
	version  string
	home     string
	bindPort string
	hostPath string
}

func (info *ProfileInfo) init() {
	info.hostPath = "$PWD"

	if Version != "" {
		info.version = ":" + Version
	}
	if HostPath != "" {
		info.hostPath = HostPath
	}
}

func generateComposeFile(info ProfileInfo) string {
	compose := `version: '3.3'
services:
  %s:
    container_name: %s
    image: %s%s
    volumes:
      - '%s:%s'
    tty: true
    stdin_open: true
`
	if info.bindPort != "" {
		compose += `
    ports:
      - '%s'
`
		return fmt.Sprintf(compose, info.name, info.name, info.image, info.version, info.hostPath, info.home, info.bindPort)
	}

	return fmt.Sprintf(compose, info.name, info.name, info.image, info.version, info.hostPath, info.home)
}

func helmProfile() error {
	info := ProfileInfo{
		name:  "helm",
		image: "azalax/helm",
		home:  "/home/helm/workspace",
	}
	info.init()
	compose := generateComposeFile(info)

	if err := createComposeFile(composeFileName, compose); err != nil {
		return err
	}
	if err := runComposeFile(Profile); err != nil {
		return err
	}

	return cleanComposeFile(composeFileName)
}

func goProfile() error {
	info := ProfileInfo{
		name:  "go",
		image: "azalax/golang",
		home:  "/home/go/workspace",
	}
	info.init()
	compose := generateComposeFile(info)

	if err := createComposeFile(composeFileName, compose); err != nil {
		return err
	}
	if err := runComposeFile(Profile); err != nil {
		return err
	}
	return cleanComposeFile(composeFileName)
}

func jekyllProfile() error {
	info := ProfileInfo{
		name:     "jekyll",
		image:    "azalax/jekyll",
		home:     "/home/jekyll/workspace",
		bindPort: "127.0.0.1:8282:4000/tcp",
	}
	info.init()
	compose := generateComposeFile(info)

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
