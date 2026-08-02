package main

import (
	"os"
	"os/exec"
	"testing"
)

func Test_All(t *testing.T) {
	Test_ValidInputs(t)
	Test_InvalidInputs(t)
}

// By default, the tests are run from $HOME, so that proper
// testing can be carried out.
func Test_ValidInputs(t *testing.T) {

	dir, _ := os.Getwd()
	t.Log(os.ReadDir(dir))

	err := exec.Command("gofind", "gofind.go").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "--help").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "-h").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "-l", "test-file.txt").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "--literal", "test-file.txt").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "*.mp4").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "-s", os.Getenv("HOME")+"/bin", "gofind").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "--source", os.Getenv("HOME")+"/bin", "gofind").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "-s", "-d", "/home", "Downloads").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "-s", "-d", "/home", `"*loads"`).Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "--source", "--directory", "/home", `"*loads"`).Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "-t", "gofind.go").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "--time-elapsed", "gofind.go").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "-l", "*.mp4").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "--literal", "*.mp4").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}

	err = exec.Command("gofind", "*find*").Run()
	if err != nil {
		t.Errorf("TEST FAILED: %v", err)
	}
}

// By default, the tests are run from $HOME, so that proper
// testing can be carried out.
func Test_InvalidInputs(t *testing.T) {

}
