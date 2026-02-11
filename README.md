# cih-mr

A command-line tool that generates merge/pull request descriptions using AI by analyzing git commits and diffs.

## Overview

`cih-mr` automates the creation of merge request descriptions by:
1. Fetching commit history between two branches
2. Analyzing the code diff
3. Using AI to generate a structured MR description based on a template

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

## Usage

```bash
cih-mr [branch-comparison] [flags]
```

### Flags

- `-b, --base` - Base branch to compare against (branch you're merging INTO) (default: `dev`)
- `-f, --head` - Head branch with changes (branch you're merging FROM) (default: current branch)
- `-m, --model` - AI model to use for generating descriptions (default: `gemini-3-pro`)
- `-p, --prompt` - Path to prompt file with AI instructions and project context (default: `prompt-context.md` if exists)

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

**Legacy usage (still supported):**
```bash
# Direct branch comparison string
cih-mr origin/dev..origin/feat

# Local branches
cih-mr main..feature-branch
```

## Output

The tool generates a structured merge request description following this template:

- **Summary** - Brief overview of what the MR accomplishes
- **Changes Made** - List of key changes
- **Testing** - Description of how changes were tested
- **Related Issues** - Links to related tickets or issues

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [go-openai](https://github.com/sashabaranov/go-openai) - OpenAI API client

## Requirements

- Go 1.25+
- Git installed and accessible in PATH
- Valid AI API key

## License

MIT
