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

	if err = os.Chdir(os.Getenv("HOME")); err != nil {
		t.Errorf("could not run command: %v", err)
	}

	err = os.Mkdir("bin", 0750)
	if err != nil && !os.IsExist(err) {
		t.Errorf("mkdir failed: %v", err)
	}

	cmd = exec.Command("cp", "gofind/gofind/gofind", os.Getenv("HOME")+"/bin")
	_, err = cmd.Output()
	if err != nil {
		t.Errorf("could not run command: %v", err)
	}

	if err = os.Chdir(os.Getenv("HOME")); err != nil {
		t.Errorf("could not run command: %v", err)
	}
}

// By default, the tests are run from $HOME, so that proper
// testing can be carried out.
func Test_ValidInputs(t *testing.T) {
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
