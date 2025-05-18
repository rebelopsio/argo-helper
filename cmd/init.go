package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	repoPath     string
	projectName  string
	withExamples bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new ArgoCD repository structure",
	Long: `Initialize a new ArgoCD repository with an opinionated structure.
This will create the necessary directories and files for managing your
applications with ArgoCD, following a Helm-like structure including:

- Custom resources directory for CRDs
- Values directory for environment-specific values
- Templates for ArgoCD applications and projects
- Helper templates for common functions`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Local flags
	initCmd.Flags().StringVarP(&projectName, "project", "p", "", "name of the ArgoCD project (required)")
	initCmd.Flags().BoolVarP(&withExamples, "examples", "e", false, "include example applications and ApplicationSet")
	if err := initCmd.MarkFlagRequired("project"); err != nil {
		fmt.Println("Error marking flag as required:", err)
	}
}

// SetInitFlags sets the flags for the init command
// This is used by the TUI to set the flags programmatically
func SetInitFlags(cmd *cobra.Command, project string, examples bool) {
	projectName = project
	withExamples = examples
}

// RunInit is exported for use by the TUI
func RunInit(cmd *cobra.Command, args []string) error {
	return runInit(cmd, args)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Get the repository path from args or use current directory
	if len(args) > 0 {
		repoPath = args[0]
	} else {
		var err error
		repoPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// If dry run is enabled, just print what would be created
	if viper.GetBool("dry-run") {
		return printDryRun()
	}

	// Create the directory structure
	if err := createRepoStructure(); err != nil {
		return err
	}

	// Print success message and next steps
	fmt.Printf("\n🎉 ArgoCD repository structure successfully created at %s\n\n", repoPath)
	fmt.Println("Next steps:")
	fmt.Println("1. Update the values.yaml file with your repository URL and other settings")
	fmt.Println("2. Create your application templates in templates/apps/")
	fmt.Println("3. Add environment-specific values in values/")

	if withExamples {
		fmt.Println("\nExample files have been created to help you get started:")
		fmt.Println("- templates/apps/example-app.yaml - Example application template")
		fmt.Println("- examples/applicationset.yaml - Example ApplicationSet")
		fmt.Println("- values/dev/values.yaml and values/prod/values.yaml - Environment-specific values")
	}

	return nil
}

func createRepoStructure() error {
	// Define the directory structure
	dirs := []string{
		"custom-resources",
		"values",
		"templates/apps",
		"templates/projects",
	}

	// Add examples directories if enabled
	if withExamples {
		dirs = append(dirs, "examples", "values/dev", "values/prod")
	}

	// Create directories
	for _, dir := range dirs {
		path := filepath.Join(repoPath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		fmt.Printf("Created directory: %s\n", path)
	}

	// Create initial files
	files := map[string]string{
		".helmignore": `# Patterns to ignore when building packages.
*.tgz
*.lock
.DS_Store
.git/
.gitignore
.vscode/
*.swp
*.bak
`,
		"Chart.yaml": fmt.Sprintf(`apiVersion: v2
name: %s
description: ArgoCD applications and projects for %s
type: application
version: 0.1.0
appVersion: "1.0.0"
maintainers:
  - name: %s Team
created: %s
`, projectName, projectName, projectName, time.Now().Format("2006-01-02")),
		"values.yaml": fmt.Sprintf(`# Default values for %s ArgoCD applications

# Global settings
global:
  environment: dev
  project: %s
  repoURL: ""  # Set this to your Git repository URL
  targetRevision: HEAD

# ArgoCD Project settings
project:
  description: "%s ArgoCD Project"
  sourceRepos:
    - "*"  # Adjust based on your security requirements
  destinations:
    - namespace: "*"
      server: "https://kubernetes.default.svc"
  clusterResourceWhitelist:
    - group: "*"
      kind: "*"

# Application defaults
applications:
  defaults:
    syncPolicy:
      automated:
        prune: true
        selfHeal: true
      syncOptions:
        - CreateNamespace=true
`, projectName, projectName, projectName),
		"templates/_helpers.tpl": `{{/*
Common labels
*/}}
{{- define "common.labels" -}}
app.kubernetes.io/managed-by: argocd
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/part-of: {{ .Values.global.project }}
{{- end }}

{{/*
Generate application name
*/}}
{{- define "common.appName" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Values.global.project $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Generate project name
*/}}
{{- define "common.projectName" -}}
{{- printf "%s" .Values.global.project -}}
{{- end -}}
`,
		"templates/projects/project.yaml": `{{- $projectName := include "common.projectName" . -}}
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: {{ $projectName }}
  namespace: argocd
  labels:
    {{- include "common.labels" . | nindent 4 }}
spec:
  description: {{ .Values.project.description }}
  sourceRepos:
  {{- range .Values.project.sourceRepos }}
    - {{ . }}
  {{- end }}
  destinations:
  {{- range .Values.project.destinations }}
    - namespace: {{ .namespace }}
      server: {{ .server }}
  {{- end }}
  clusterResourceWhitelist:
  {{- range .Values.project.clusterResourceWhitelist }}
    - group: {{ .group }}
      kind: {{ .kind }}
  {{- end }}
`,
		"README.md": fmt.Sprintf(`# %s ArgoCD Repository

This repository contains the ArgoCD applications and projects for the %s project, structured as a Helm chart.

## Structure

- `+"`custom-resources/`"+`: Contains Custom Resource Definitions (CRDs) if needed
- `+"`values/`"+`: Contains environment-specific values files
- `+"`templates/`"+`:
  - `+"`apps/`"+`: Application templates
  - `+"`projects/`"+`: Project templates
  - `+"`_helpers.tpl`"+`: Common template helpers
- `+"`values.yaml`"+`: Default values
- `+"`Chart.yaml`"+`: Chart metadata

## Usage

1. Update the `+"`values.yaml`"+` file with your repository URL and other settings
2. Add your application templates in `+"`templates/apps/`"+`
3. Add environment-specific values in `+"`values/`"+`
4. Use `+"`helm template`"+` to generate manifests or commit to your ArgoCD repository

## Adding New Applications

Create a new application template in `+"`templates/apps/`"+` following this pattern:

`+"```yaml"+`
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{ include "common.appName" . }}
  namespace: argocd
spec:
  project: {{ include "common.projectName" . }}
  source:
    repoURL: {{ .Values.global.repoURL }}
    targetRevision: {{ .Values.global.targetRevision }}
    path: apps/your-app
  destination:
    server: {{ .Values.destination.server }}
    namespace: {{ .Values.destination.namespace }}
  syncPolicy:
    {{- toYaml .Values.applications.defaults.syncPolicy | nindent 4 }}
`+"```"+`
`, projectName, projectName),
	}

	// Add example files if enabled
	if withExamples {
		exampleFiles := map[string]string{
			"templates/apps/example-app.yaml": fmt.Sprintf(`apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{ include "common.appName" . }}-example
  namespace: argocd
  labels:
    {{- include "common.labels" . | nindent 4 }}
spec:
  project: {{ include "common.projectName" . }}
  source:
    repoURL: {{ .Values.global.repoURL }}
    targetRevision: {{ .Values.global.targetRevision }}
    path: apps/example-app
  destination:
    server: "{{ .Values.destination.server | default "https://kubernetes.default.svc" }}"
    namespace: example
  syncPolicy:
    {{- toYaml .Values.applications.defaults.syncPolicy | nindent 4 }}
`),
			"examples/applicationset.yaml": `apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: {{ include "common.projectName" . }}-apps
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
      name: '{{ "{{path.basename}}" }}'
      labels:
        {{- include "common.labels" . | nindent 8 }}
    spec:
      project: {{ include "common.projectName" . }}
      source:
        repoURL: {{ .Values.global.repoURL }}
        targetRevision: {{ .Values.global.targetRevision }}
        path: '{{ "{{path}}" }}'
      destination:
        server: {{ .Values.destination.server | default "https://kubernetes.default.svc" }}
        namespace: '{{ "{{path.basename}}" }}'
      syncPolicy:
        {{- toYaml .Values.applications.defaults.syncPolicy | nindent 8 }}
`,
			"values/dev/values.yaml": fmt.Sprintf(`# Development environment values for %s

global:
  environment: dev

# Override specific application values for development
`, projectName),
			"values/prod/values.yaml": fmt.Sprintf(`# Production environment values for %s

global:
  environment: prod

# Override specific application values for production
# Make sure to be careful with production configurations
`, projectName),
		}

		// Merge the example files with the main files map
		for filename, content := range exampleFiles {
			files[filename] = content
		}
	}

	// Write all files
	for filename, content := range files {
		path := filepath.Join(repoPath, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", filename, err)
		}
		fmt.Printf("Created file: %s\n", path)
	}

	return nil
}

func printDryRun() error {
	fmt.Println("Dry run: The following structure would be created:")
	fmt.Printf("\nRoot directory: %s\n\n", repoPath)

	// Print directory structure
	items := []string{
		".helmignore",
		"Chart.yaml",
		"values.yaml",
		"custom-resources/",
		"values/",
		"templates/",
		"templates/apps/",
		"templates/projects/",
		"templates/_helpers.tpl",
		"templates/projects/project.yaml",
		"README.md",
	}

	// Add example items if enabled
	if withExamples {
		exampleItems := []string{
			"examples/",
			"examples/applicationset.yaml",
			"values/dev/",
			"values/dev/values.yaml",
			"values/prod/",
			"values/prod/values.yaml",
			"templates/apps/example-app.yaml",
		}
		items = append(items, exampleItems...)
	}

	for _, item := range items {
		fmt.Printf("  %s\n", item)
	}

	// Print completion message
	fmt.Printf("\nTo create this structure, run again without the --dry-run flag\n")
	if !withExamples {
		fmt.Printf("Add --examples or -e flag to include example applications and values\n")
	}

	return nil
}
