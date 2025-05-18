package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewCommand(t *testing.T) {
	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "argo-helper-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up test cases
	testCases := []struct {
		name         string
		resourceType string
		resourceName string
		outputPath   string
		shouldError  bool
	}{
		{
			name:         "Valid ApplicationSet",
			resourceType: "applicationset",
			resourceName: "test-apps",
			outputPath:   tempDir,
			shouldError:  false,
		},
		{
			name:         "Invalid Resource Type",
			resourceType: "invalid-type",
			resourceName: "test-invalid",
			outputPath:   tempDir,
			shouldError:  true,
		},
		{
			name:         "Empty Resource Name",
			resourceType: "applicationset",
			resourceName: "",
			outputPath:   tempDir,
			shouldError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new command for testing
			cmd := &cobra.Command{}
			
			// Set the flags
			SetNewFlags(cmd, tc.resourceType, tc.resourceName, tc.outputPath)
			
			// Run the command
			args := []string{tc.resourceType}
			if tc.resourceName != "" {
				args = append(args, tc.resourceName)
			}
			
			err := runNew(cmd, args)
			
			// Check if error occurred as expected
			if tc.shouldError && err == nil {
				t.Errorf("Expected error but got none")
			}
			
			if !tc.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			// If no error, check that the file was created
			if !tc.shouldError {
				expectedFilePath := filepath.Join(tc.outputPath, 
					tc.resourceType+"-"+tc.resourceName+".yaml")
				
				if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
					t.Errorf("Expected file was not created: %s", expectedFilePath)
				}
				
				// Clean up the generated file
				os.Remove(expectedFilePath)
			}
		})
	}
}