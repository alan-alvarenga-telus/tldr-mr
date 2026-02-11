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

Set the `AI_KEY` environment variable with your API key:

```bash
export AI_KEY=your-api-key
```

## Usage

```bash
cih-mr <branch-comparison> [flags]
```

### Flags

- `-m, --model` - AI model to use for generating descriptions (default: `gemini-3-pro`)

### Examples

Compare a feature branch against the dev branch:
```bash
cih-mr origin/dev..origin/feat
```

Compare local branches:
```bash
cih-mr main..feature-branch
```

Use a specific AI model:
```bash
cih-mr main..dev/feat1 --model gemini-pro-3
```

Or use the short flag:
```bash
cih-mr main..dev/feat1 -m gpt-4
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
