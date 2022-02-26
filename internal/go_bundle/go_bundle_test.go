package go_bundle_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/njkevlani/go_bundle/internal/go_bundle"
	"github.com/stretchr/testify/assert"
)

func TestGoBundle(t *testing.T) {
  fmt.Println(os.Getwd())
	testcases := []struct {
		inputFilePath          string
		expectedOutputFilePath string
	}{
		{
			inputFilePath:          "./test_files/test_project1//main.go",
			expectedOutputFilePath: "./test_files/expected_output1.go",
		},
	}

	for _, testcase := range testcases {
		output, err := go_bundle.GoBundle(testcase.inputFilePath)
		assert.NoError(t, err)

		expectedOutput, err := os.ReadFile(testcase.expectedOutputFilePath)
		assert.NoError(t, err)

    assert.Equal(t, string(expectedOutput), string(output))
	}
}
