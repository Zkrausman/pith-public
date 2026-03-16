package parser

import (
	"strings"
	"testing"
)

func TestEnvParser(t *testing.T) {
	p := &EnvParser{}
	input := `PATH=/usr/bin
USER=zkrau
GITHUB_TOKEN=secret_token
PASSWORD=12345
SHELL=/bin/bash
`
	output := p.Parse(input)
	if strings.Contains(output, "GITHUB_TOKEN") || strings.Contains(output, "PASSWORD") {
		t.Error("EnvParser should redact secrets")
	}
	if !strings.Contains(output, "USER=zkrau") {
		t.Error("EnvParser should keep normal variables")
	}
}

func TestDockerPsParser(t *testing.T) {
	p := &DockerPsParser{}
	input := `CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
1234567890ab   nginx     "entry"   1h ago    Up        80/tcp    my-nginx
`
	output := p.Parse(input)
	if !strings.Contains(output, "1234567890ab") || !strings.Contains(output, "my-nginx") {
		t.Errorf("DockerPsParser missing essential info: %s", output)
	}
}

func TestDependencyParser(t *testing.T) {
	p := &DependencyParser{}
	input := `diet@0.3.0 E:\Repos\Diet
├── github.com/spf13/cobra@v1.8.0
└── github.com/AlecAivazis/survey/v2@v2.3.7
`
	output := p.Parse(input)
	if strings.Contains(output, "├──") {
		t.Error("DependencyParser should strip visual tree characters")
	}
	if !strings.Contains(output, "github.com/spf13/cobra") {
		t.Error("DependencyParser missing package name")
	}
}

func TestTestParser(t *testing.T) {
	p := &TestParser{}
	input := `Running tests...
PASS: TestOne
FAIL: TestTwo
    error: expected true, got false
DONE: 1 passed, 1 failed
`
	output := p.Parse(input)
	if !strings.Contains(output, "FAIL: TestTwo") || !strings.Contains(output, "1 passed") {
		t.Errorf("TestParser missing summary or failure: %s", output)
	}
}
