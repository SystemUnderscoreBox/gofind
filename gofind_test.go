package main

import (
	"fmt"
	"testing"
)

func testAll(t *testing.T) {
	if err := testValidInputs(t); err != nil {
		fmt.Printf("testValidInputs failed, got: %v", err)
	}

	if err := testInvalidInputs(t); err != nil {
		fmt.Printf("testInvalidInputs failed, got: %v", err)
	}
}

func testValidInputs(t *testing.T) error {
	return nil
}

func testInvalidInputs(t *testing.T) error {
	return nil
}

func miscellaneous(t *testing.T) {

}
