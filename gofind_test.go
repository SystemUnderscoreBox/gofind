package main

import (
	"os"
	"os/exec"
	"testing"
)

func Test_All(t *testing.T) {
	preTest(t)

	Test_ValidInputs(t)
	Test_InvalidInputs(t)
}

// Used only to run shell commands before testing.
// That is, compile and copy to executable file path.
func preTest(t *testing.T) {
	cmd := exec.Command("go", "build")
	_, err := cmd.Output()
	if err != nil {
		t.Errorf("could not run command: %v", err)
	}

	cmd = exec.Command("cp", "gofind", os.Getenv("HOME")+"/bin/gofind")
	_, err = cmd.Output()
	if err != nil {
		t.Errorf("could not run command: %v", err)
	}
}

// By default, the tests are run from $HOME, so that proper
// testing can be carried out.
func Test_ValidInputs(t *testing.T) {
}

// By default, the tests are run from $HOME, so that proper
// testing can be carried out.
func Test_InvalidInputs(t *testing.T) {
}
