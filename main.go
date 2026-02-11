package main

import (
	"fmt"
	"os"
	"strings"

	"cih-mr/ai"
	"cih-mr/git"

	"github.com/spf13/cobra"
)

var modelFlag string

var rootCmd = &cobra.Command{
	Use:   "cih-mr [branch-comparison]",
	Short: "Generate merge request description using AI",
	Long: `Analyzes git changes between banches and generates a merge request
	description using AI based on commits and diffs.
	Example: cih-mr origin/dev..origin/feat`,
	Args: cobra.ExactArgs(1),
	Run:  runMRGenerator,
}

func init() {
	rootCmd.Flags().StringVarP(&modelFlag, "model", "m", "gemini-3-pro", "AI model to use for generating descriptions")
}

func runMRGenerator(cmd *cobra.Command, args []string) {
	branchComparison := args[0]
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

	// Generate MR description
	description, err := aiClient.GenerateMRDescription(commits, diff, modelFlag)
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
