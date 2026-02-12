package main

import (
	"fmt"
	"os"
	"strings"

	"cih-mr/ai"
	"cih-mr/git"

	"github.com/spf13/cobra"
)

const (
	defaultPromptFile   = "prompt-context.md"
	defaultTemplateFile = "template.md"
)

var (
	modelFlag    string
	baseFlag     string
	headFlag     string
	promptFlag   string
	templateFlag string
)

var rootCmd = &cobra.Command{
	Use:   "cih-mr [branch-comparison]",
	Short: "Generate merge request description using AI",
	Long: `Analyzes git changes between branches and generates a merge request
	description using AI based on commits and diffs.
	
	Examples:
	  cih-mr                              # Compare current branch against dev
	  cih-mr -b main                      # Compare current branch against main
	  cih-mr -b dev -h feature-123        # Compare specific branches
	  cih-mr origin/dev..origin/feat      # Legacy: direct comparison string`,
	Args: cobra.MaximumNArgs(1),
	Run:  runMRGenerator,
}

// init registers command-line flags for the root command
func init() {
	rootCmd.Flags().StringVarP(&modelFlag, "model", "m", "gemini-3-pro", "AI model to use for generating descriptions")
	rootCmd.Flags().StringVarP(&baseFlag, "base", "b", "dev", "Base branch to compare against (branch you're merging INTO)")
	rootCmd.Flags().StringVarP(&headFlag, "head", "f", "", "Head branch with changes (branch you're merging FROM, defaults to current branch)")
	rootCmd.Flags().StringVarP(&promptFlag, "prompt", "p", "", "Path to prompt file with AI instructions and project context (default: prompt-context.md if exists)")
	rootCmd.Flags().StringVarP(&templateFlag, "template", "t", "", "Path to template file for output structure (default: template.md if exists)")
}

// fatalError prints an error message and exits with status code 1
func fatalError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// loadFileWithDefault loads a file with fallback to default filename
// Returns empty string if neither explicit nor default file exists (unless explicit was provided)
func loadFileWithDefault(explicitPath, defaultPath, fileType string) (string, error) {
	if explicitPath != "" {
		// Explicit flag provided - must exist
		content, err := os.ReadFile(explicitPath)
		if err != nil {
			return "", fmt.Errorf("failed to read %s file %s: %v", fileType, explicitPath, err)
		}
		return string(content), nil
	}

	// Check for default file
	if content, err := os.ReadFile(defaultPath); err == nil {
		return string(content), nil
	}

	// No file found - return empty to use built-in default
	return "", nil
}

// determineBranchComparison builds the branch comparison string from args or flags
func determineBranchComparison(args []string) (string, error) {
	if len(args) > 0 {
		// Legacy mode: use positional argument
		return args[0], nil
	}

	// New mode: build from flags
	head := headFlag
	if head == "" {
		// Auto-detect current branch
		currentBranch, err := git.GetCurrentBranch()
		if err != nil {
			return "", err
		}
		head = currentBranch
	}
	return fmt.Sprintf("%s...%s", baseFlag, head), nil
}

// fetchGitData retrieves commit history and diff for the given branch comparison
func fetchGitData(branchComparison string) (commits, diff string, err error) {
	fmt.Println("📝 Fetching commit history...")
	commits, err = git.GetLog(branchComparison)
	if err != nil {
		return "", "", err
	}

	fmt.Println("📊 Fetching changes...")
	diff, err = git.GetDiff(branchComparison)
	if err != nil {
		return "", "", err
	}

	return commits, diff, nil
}

// loadConfigFiles loads prompt and template files with fallback to defaults
func loadConfigFiles() (systemPrompt, template string, err error) {
	systemPrompt, err = loadFileWithDefault(promptFlag, defaultPromptFile, "prompt")
	if err != nil {
		return "", "", err
	}

	template, err = loadFileWithDefault(templateFlag, defaultTemplateFile, "template")
	if err != nil {
		return "", "", err
	}

	return systemPrompt, template, nil
}

// printResult outputs the generated MR description in a formatted box
func printResult(description string) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("MERGE REQUEST DESCRIPTION")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(description)
	fmt.Println(strings.Repeat("=", 60))
}

func runMRGenerator(cmd *cobra.Command, args []string) {
	// Determine branch comparison
	branchComparison, err := determineBranchComparison(args)
	if err != nil {
		fatalError(err)
	}

	fmt.Printf("Analyzing: %s\n\n", branchComparison)

	// Fetch git data
	commits, diff, err := fetchGitData(branchComparison)
	if err != nil {
		fatalError(err)
	}

	// Load configuration files
	systemPrompt, template, err := loadConfigFiles()
	if err != nil {
		fatalError(err)
	}

	// Create AI client and generate description
	fmt.Println("🤖 Generating MR description...")
	aiClient, err := ai.NewClient()
	if err != nil {
		fatalError(err)
	}

	description, err := aiClient.GenerateMRDescription(commits, diff, modelFlag, systemPrompt, template)
	if err != nil {
		fatalError(err)
	}

	// Output result
	printResult(description)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
