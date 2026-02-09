# Milestone 3.3: Quick Reference

## Summary

Intelligent merge conflict resolution. Automatically resolve Git merge conflicts using semantic code analysis and contextual understanding.

## Timeline

- **Duration:** 3–4 weeks
- **Tasks:** 25 sequential
- **Target Completion:** Week 26 of Phase 3

## Example Usage

```bash
# Check for conflicts and resolve
git-surgeon --file conflicted-file.py

# Auto-resolve all conflicts in working directory
git-surgeon --auto-resolve

# Show proposed resolutions
git-surgeon --show

# Accept resolution
git-surgeon --file file.py --accept

# Dry-run to preview changes
git-surgeon --file file.py --dry-run
```

## Success Criteria

| Criterion | Target | Status |
|-----------|--------|--------|
| Auto-Resolution Rate | >70% | — |
| Resolution Accuracy | >95% | — |
| Code Coverage | >80% | — |

## Acceptance Criteria

- ✅ Parses merge conflicts correctly
- ✅ >70% auto-resolution rate
- ✅ >95% resolution accuracy
- ✅ Validates resolutions
- ✅ Git integration working
- ✅ >80% code coverage
