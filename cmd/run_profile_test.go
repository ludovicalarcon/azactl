package cmd

import (
	"testing"
)

func TestGenerateComposeFile_Should_handle_basic(t *testing.T) {
	info := ProfileInfo{
		name:  "foo",
		image: "bar",
		home:  "/foo/bar",
	}
	want := `version: '3.3'
services:
  foo:
    container_name: foo
    image: bar
    volumes:
      - '$PWD:/foo/bar'
    tty: true
    stdin_open: true
`

	got := generateComposeFile(info)
	if got != want {
		t.Errorf("got = %s; want %s", got, want)
	}
}

func TestGenerateComposeFile_Should_handle_versionOverride(t *testing.T) {
	Version = "X.Y.Z"
	info := ProfileInfo{
		name:  "foo",
		image: "bar",
		home:  "/foo/bar",
	}
	want := `version: '3.3'
services:
  foo:
    container_name: foo
    image: bar:X.Y.Z
    volumes:
      - '$PWD:/foo/bar'
    tty: true
    stdin_open: true
`
	info.init()
	got := generateComposeFile(info)
	if got != want {
		t.Errorf("got = %s; want %s", got, want)
	}
}

func TestGenerateComposeFile_Should_handle_bindingPort(t *testing.T) {
	info := ProfileInfo{
		name:     "foo",
		image:    "bar",
		home:     "/foo/bar",
		bindPort: "XXX:YYY",
	}
	want := `version: '3.3'
services:
  foo:
    container_name: foo
    image: bar
    volumes:
      - '$PWD:/foo/bar'
    tty: true
    stdin_open: true
`
	want += `
    ports:
      - 'XXX:YYY'
`

	got := generateComposeFile(info)
	if got != want {
		t.Errorf("got = %s; want %s", got, want)
	}
}
