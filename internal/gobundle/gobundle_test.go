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
