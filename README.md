# cih-mr

A flexible command-line tool that uses AI to generate documentation from git commits and diffs. Perfect for creating merge request descriptions, changelogs, release notes, and more.

## Overview

`cih-mr` automates the creation of git-based documentation by:
1. Fetching commit history between two branches
2. Analyzing the code diff
3. Using AI to generate structured output based on customizable templates and prompts

**Use Cases:**
- 📝 Merge/Pull Request descriptions
- 📋 Changelogs
- 🚀 Release notes
- 📊 Sprint summaries
- 📖 Any git history-based documentation

## Installation

```bash
go install cih-mr@latest
```

Or build from source:

```bash
git clone <repository-url>
cd cih-mr
go build -o cih-mr
```

## Configuration

### API Key

Set the `AI_KEY` environment variable with your API key:

```bash
export AI_KEY=your-api-key
```

### Custom AI Prompt (Optional)

Create a `prompt-context.md` file in your project root to customize how the AI generates descriptions. This file should include:
- Instructions for the AI's writing style and tone
- Project-specific context (what the project does, tech stack, team conventions)
- Any specific guidance for MR descriptions

**Example `prompt-context.md`:**

```markdown
You are a developer writing a merge request description for your team.

Write like a human. Be direct. Use simple words. Avoid buzzwords and jargon.

Your goal is clarity, not perfection. Reviewers should understand:
- What changed
- Why it changed
- What to watch out for

Keep it professional but real. Be concise. Skip the fluff.

## About This Project

We're building an e-commerce platform for small businesses. The main users are 
shop owners who need reliability over fancy features. We prioritize uptime and
data accuracy.

Tech stack: Go backend, React frontend, PostgreSQL database.
```

The tool will automatically use `prompt-context.md` if it exists. You can also specify a different file:

```bash
cih-mr --prompt ./custom-prompt.md
cih-mr -p ./prompts/detailed.md
```

If no prompt file is provided, the tool uses a sensible built-in default.

### Custom Output Templates (Optional)

Create a `template.md` file to define the structure and format of the generated output. This makes the tool flexible for different use cases.

**Example `template.md` (for MR descriptions):**

```markdown
# Merge Request Description

## Summary
[Brief overview]

## Changes Made
[List key changes]

## Testing Done
[Describe testing]

## Breaking Changes
[List or note "None"]
```

The tool will automatically use `template.md` if it exists. You can also specify a different template:

```bash
cih-mr --template ./examples/changelog-template.md
cih-mr -t ./examples/release-notes-template.md
```

If no template file is provided, the tool uses a built-in MR template.

## Usage

```bash
cih-mr [branch-comparison] [flags]
```

### Flags

- `-b, --base` - Base branch to compare against (branch you're merging INTO) (default: `dev`)
- `-f, --head` - Head branch with changes (branch you're merging FROM) (default: current branch)
- `-m, --model` - AI model to use for generating descriptions (default: `gemini-3-pro`)
- `-p, --prompt` - Path to prompt file with AI instructions and project context (default: `prompt-context.md` if exists)
- `-t, --template` - Path to template file for output structure (default: `template.md` if exists)

### Examples

**Simple usage (recommended):**
```bash
# Compare current branch against dev (uses defaults)
cih-mr

# Compare current branch against main
cih-mr -b main

# Compare specific branches explicitly
cih-mr -b dev -f feature-123

# With custom AI model
cih-mr -b main -m gpt-4

# With custom prompt file
cih-mr -b main --prompt ./prompts/detailed.md
```

**Different output types:**
```bash
# Generate a changelog
cih-mr -t examples/changelog-template.md -p examples/changelog-prompt.md

# Generate release notes
cih-mr -t examples/release-notes-template.md -p examples/release-notes-prompt.md

# Custom workflow for your team
cih-mr -t ./team-templates/sprint-summary.md -p ./team-prompts/agile.md
```

**Legacy usage (still supported):**
```bash
# Direct branch comparison string
cih-mr origin/dev..origin/feat

# Local branches
cih-mr main..feature-branch
```

## Output

The tool generates structured output based on the template you provide (or the built-in default). The default template produces a merge request description with:

- **Summary** - Brief overview of what this accomplishes
- **Changes Made** - List of key changes
- **Type of Change** - Categorization (feature, bug fix, etc.)
- **Testing Done** - Description of testing approach
- **Breaking Changes** - Any breaking changes
- **Related Issues** - Links to related tickets

## Example Templates

The repository includes example templates and prompts in the `examples/` directory:

- **Changelog**: `changelog-template.md` + `changelog-prompt.md`
- **Release Notes**: `release-notes-template.md` + `release-notes-prompt.md`

These demonstrate how to customize the tool for different documentation needs.

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [go-openai](https://github.com/sashabaranov/go-openai) - OpenAI API client

## Requirements

- Go 1.25+
- Git installed and accessible in PATH
- Valid AI API key

## License

MIT
