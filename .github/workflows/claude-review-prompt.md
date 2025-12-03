# Pull Request Review

**REPO:** {{REPO}}
**PR:** {{PR_NUMBER}}

## Instructions

Review only the **changed files** in this pull request.

### Review Checklist

1. **Security**: Check for exposed secrets, credentials, API keys, or security vulnerabilities
2. **Go Best Practices**: Ensure code follows standard Go conventions and idioms
3. **Bugs**: Identify potential bugs, logic errors, or edge cases
4. **Code Quality**: Check for obvious issues in implementation

### Output Rules

- **If PR is good**: Post only `LGTM` - nothing else
- **If issues found**: Be concise. List specific issues with file:line references
- **Always use**: `gh pr comment {{PR_NUMBER}} --body "your review"` to post your review

### Steps

1. Run `gh pr diff {{PR_NUMBER}}` to see changes
2. Review the changed files only
3. Post your review using `gh pr comment`
