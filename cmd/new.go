package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	resourceType string
	resourceName string
	outputPath   string
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [resource-type] [resource-name]",
	Short: "Create a new ArgoCD resource",
	Long: `Create a new ArgoCD resource with an opinionated template.
Currently supported resource types:
- applicationset: Create a new ApplicationSet manifest

The resources will be created in the templates/apps/ directory by default.`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    runNew,
	Example: "  argo-helper new applicationset my-apps",
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Local flags
	newCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output path (default is templates/apps/)")
}

// SetNewFlags sets the flags for the new command
// This is used by the TUI to set the flags programmatically
func SetNewFlags(cmd *cobra.Command, rType string, rName string, output string) {
	// Set global variables directly since this will be used from TUI as well
	resourceType = rType
	resourceName = rName
	outputPath = output
}

// RunNew is exported for use by the TUI
func RunNew(cmd *cobra.Command, args []string) error {
	return runNew(cmd, args)
}

func runNew(cmd *cobra.Command, args []string) error {
	// Parse arguments
	resourceType = args[0]
	if len(args) > 1 {
		resourceName = args[1]
	}

	// Validate resource type
	if resourceType != "applicationset" {
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	// If resourceName is not provided via argument or flag, prompt for it
	if resourceName == "" {
		return fmt.Errorf("resource name is required")
	}

	// Set default output path if not provided
	if outputPath == "" {
		outputPath = "templates/apps"
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// If dry run is enabled, just print what would be created
	if viper.GetBool("dry-run") {
		return printNewDryRun()
	}

	// Create the resource
	if err := createResource(); err != nil {
		return err
	}

	// Print success message
	filename := fmt.Sprintf("%s-%s.yaml", resourceType, resourceName)
	fmt.Printf("\n✅ %s '%s' successfully created at %s\n\n",
		capitalizeFirstLetter(resourceType),
		resourceName,
		filepath.Join(outputPath, filename))

	fmt.Println("Next steps:")
	fmt.Println("1. Review and customize the generated resource")
	fmt.Println("2. Apply to your ArgoCD instance or commit to your repository")

	return nil
}

func createResource() error {
	// Define the content based on the resource type
	var content string

	switch resourceType {
	case "applicationset":
		content = generateApplicationSetTemplate()
	}

	// Write the file
	filename := fmt.Sprintf("%s-%s.yaml", resourceType, resourceName)
	filePath := filepath.Join(outputPath, filename)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}

	fmt.Printf("Created file: %s\n", filePath)
	return nil
}

func generateApplicationSetTemplate() string {
	return fmt.Sprintf(`apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: %s
  namespace: argocd
spec:
  generators:
    - git:
        repoURL: {{ .Values.global.repoURL }}
        revision: {{ .Values.global.targetRevision }}
        directories:
          - path: apps/*
  template:
    metadata:
      name: '{{ "{{ path.basename }}" }}'
      labels:
        {{- include "common.labels" . | nindent 8 }}
    spec:
      project: {{ include "common.projectName" . }}
      source:
        repoURL: {{ .Values.global.repoURL }}
        targetRevision: {{ .Values.global.targetRevision }}
        path: '{{ "{{ path }}" }}'
      destination:
        server: {{ .Values.destination.server | default "https://kubernetes.default.svc" }}
        namespace: '{{ "{{ path.basename }}" }}'
      syncPolicy:
        {{- toYaml .Values.applications.defaults.syncPolicy | nindent 8 }}
`, resourceName)
}

func printNewDryRun() error {
	fmt.Println("Dry run: The following resource would be created:")
	fmt.Printf("\nResource Type: %s\n", resourceType)
	fmt.Printf("Resource Name: %s\n", resourceName)
	fmt.Printf("Output Path: %s\n\n", outputPath)

	filename := fmt.Sprintf("%s-%s.yaml", resourceType, resourceName)
	fmt.Printf("File: %s\n\n", filepath.Join(outputPath, filename))

	fmt.Println("Template content:")
	fmt.Println("---")

	switch resourceType {
	case "applicationset":
		fmt.Println(generateApplicationSetTemplate())
	}

	fmt.Println("---")
	fmt.Printf("\nTo create this resource, run again without the --dry-run flag\n")

	return nil
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	if 'a' <= s[0] && s[0] <= 'z' {
		return string(s[0]-32) + s[1:]
	}
	return s
}
