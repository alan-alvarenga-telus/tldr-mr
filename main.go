package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"cih-mr/ai"
	"cih-mr/git"

	"github.com/spf13/cobra"
)

var (
	modelFlag  string
	baseFlag   string
	headFlag   string
	promptFlag string
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

func init() {
	rootCmd.Flags().StringVarP(&modelFlag, "model", "m", "gemini-3-pro", "AI model to use for generating descriptions")
	rootCmd.Flags().StringVarP(&baseFlag, "base", "b", "dev", "Base branch to compare against (branch you're merging INTO)")
	rootCmd.Flags().StringVarP(&headFlag, "head", "f", "", "Head branch with changes (branch you're merging FROM, defaults to current branch)")
	rootCmd.Flags().StringVarP(&promptFlag, "prompt", "p", "", "Path to prompt file with AI instructions and project context (default: prompt-context.md if exists)")
}

// getCurrentBranch returns the name of the current git branch
func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// loadPromptFile loads the system prompt from file or returns empty string for built-in default
func loadPromptFile() (string, error) {
	var promptPath string

	if promptFlag != "" {
		// Explicit flag provided - must exist
		promptPath = promptFlag
		content, err := os.ReadFile(promptPath)
		if err != nil {
			return "", fmt.Errorf("failed to read prompt file %s: %v", promptPath, err)
		}
		return string(content), nil
	}

	// Check for default prompt-context.md
	promptPath = "prompt-context.md"
	if content, err := os.ReadFile(promptPath); err == nil {
		return string(content), nil
	}

	// No prompt file found - return empty to use built-in default
	return "", nil
}

func runMRGenerator(cmd *cobra.Command, args []string) {
	var branchComparison string

	// Build branch comparison string
	if len(args) > 0 {
		// Legacy mode: use positional argument
		branchComparison = args[0]
	} else {
		// New mode: build from flags
		head := headFlag
		if head == "" {
			// Auto-detect current branch
			currentBranch, err := getCurrentBranch()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			head = currentBranch
		}
		branchComparison = fmt.Sprintf("%s..%s", baseFlag, head)
	}

	fmt.Printf("Analyzing: %s\n\n", branchComparison)

	// Get git log
	fmt.Println("📝 Fetching commit history...")
	commits, err := git.GetLog(branchComparison)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get git diff
	fmt.Println("📊 Fetching changes...")
	diff, err := git.GetDiff(branchComparison)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create AI client
	fmt.Println("🤖 Generating MR description...")
	aiClient, err := ai.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Load prompt file
	systemPrompt, err := loadPromptFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Generate MR description
	description, err := aiClient.GenerateMRDescription(commits, diff, modelFlag, systemPrompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output result
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("MERGE REQUEST DESCRIPTION")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(description)
	fmt.Println(strings.Repeat("=", 60))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
