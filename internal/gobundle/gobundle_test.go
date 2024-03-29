package gobundle

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoBundle(t *testing.T) {
	fmt.Println(os.Getwd())
	testcases := []struct {
		name                   string
		inputFilePath          string
		expectedOutputFilePath string
	}{
		{
			name:                   "test_project0",
			inputFilePath:          "../../test_files/test_project0/main.go",
			expectedOutputFilePath: "../../test_files/expected_output0/main.go",
		},
		{
			name:                   "test_project1",
			inputFilePath:          "../../test_files/test_project1/main.go",
			expectedOutputFilePath: "../../test_files/expected_output1/main.go",
		},
		{
			name:                   "test_project2",
			inputFilePath:          "../../test_files/test_project2/main.go",
			expectedOutputFilePath: "../../test_files/expected_output2/main.go",
		},
		{
			name:                   "test_project3",
			inputFilePath:          "../../test_files/test_project3/main.go",
			expectedOutputFilePath: "../../test_files/expected_output3/main.go",
		},
		{
			name:                   "test_project4",
			inputFilePath:          "../../test_files/test_project4/main.go",
			expectedOutputFilePath: "../../test_files/expected_output4/main.go",
		},
		{
			name:                   "test_project5",
			inputFilePath:          "../../test_files/test_project5/main.go",
			expectedOutputFilePath: "../../test_files/expected_output5/main.go",
		},
		{
			name:                   "test_project6",
			inputFilePath:          "../../test_files/test_project6/main.go",
			expectedOutputFilePath: "../../test_files/expected_output6/main.go",
		},
		{
			name:                   "test_project7",
			inputFilePath:          "../../test_files/test_project7/main.go",
			expectedOutputFilePath: "../../test_files/expected_output7/main.go",
		},
		{
			name:                   "test_project8",
			inputFilePath:          "../../test_files/test_project8/main.go",
			expectedOutputFilePath: "../../test_files/expected_output8/main.go",
		},
		{
			name:                   "test_project9",
			inputFilePath:          "../../test_files/test_project9/main.go",
			expectedOutputFilePath: "../../test_files/expected_output9/main.go",
		},
		{
			name:                   "test_project10",
			inputFilePath:          "../../test_files/test_project10/main.go",
			expectedOutputFilePath: "../../test_files/expected_output10/main.go",
		},
		{
			name:                   "test_project11",
			inputFilePath:          "../../test_files/test_project11/main.go",
			expectedOutputFilePath: "../../test_files/expected_output11/main.go",
		},
		{
			name:                   "test_project12",
			inputFilePath:          "../../test_files/test_project12/main.go",
			expectedOutputFilePath: "../../test_files/expected_output12/main.go",
		},
		{
			name:                   "test_project13",
			inputFilePath:          "../../test_files/test_project13/main.go",
			expectedOutputFilePath: "../../test_files/expected_output13/main.go",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			output, err := GoBundle(testcase.inputFilePath)
			fmt.Println(string(output))
			assert.NoError(t, err)

			expectedOutput, err := os.ReadFile(testcase.expectedOutputFilePath)
			assert.NoError(t, err)

			assert.Equal(t, string(expectedOutput), string(output))
		})
	}
}
