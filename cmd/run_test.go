package cmd

import (
	"os"
	"testing"
)

func TestCreateComposeFile_Should_CreateTheFile(t *testing.T) {
	path := t.TempDir()
	filePath := path + "/tmp.yaml"

	err := createComposeFile(filePath, "foo: bar")

	if err != nil {
		t.Fatal("Should not have return error but got", err)
	}

	if _, err := os.Stat(filePath); err != nil {
		t.Error("Compose file should have been created but it is not present")
	}
}

func TestCreateComposeFile_Should_ReturnErrInvalidFileName(t *testing.T) {
	path := t.TempDir()
	err := createComposeFile(path, "foo: bar")

	if err == nil {
		t.Error("It should have return an error as it is a directory")
	}
}

func TestCleanComposeFile_Should_DeleteTheFile(t *testing.T) {
	path := t.TempDir()
	filePath := path + "/tmp.yaml"

	_, err := os.Create(filePath)

	if err != nil {
		t.Fatal(err)
	}

	err = cleanComposeFile(filePath)

	if err != nil {
		t.Fatal("Should not have return error but got", err)
	}

	if _, err = os.Stat(filePath); err == nil {
		t.Error("Compose file should have beeen deleted but it is still present")
	}
}

func TestCleanComposeFile_Should_ReturnErrInvalidFileName(t *testing.T) {
	err := cleanComposeFile("foo.yaml")

	if err == nil {
		t.Error("It should have return an error as file does not exist")
	}
}

func TestRunCmd_Should_ReturnErrDockerNotInstalled(t *testing.T) {
	err := runCmd.RunE(runCmd, []string{})

	if err == nil {
		t.Error("It should have return an error as docker is not in PATH")
	}
}
