# Milestone 3.5: Quick Reference

## Summary

Mathematical computation co-processor. Solve equations, perform symbolic math, matrix operations, visualize results.

## Timeline

- **Duration:** 2–3 weeks
- **Tasks:** 25 sequential

## Success Criteria

| Criterion | Target |
|-----------|--------|
| Expression Parsing | <10ms |
| Solving Accuracy | 100% |
| Code Coverage | >80% |

## Acceptance Criteria

- ✅ Parses and evaluates expressions
- ✅ Solves equations symbolically
- ✅ Performs matrix operations
- ✅ Generates visualizations
- ✅ >80% code coverage
- ✅ Complete documentation

## Example Usage

```bash
# Solve equation
calc-bridge "x^2 + 2*x + 1 = 0"

# Matrix operations
calc-bridge "[[1, 2], [3, 4]] @ [[5, 6], [7, 8]]"

# Symbolic differentiation
calc-bridge "diff(x^3 + 2*x^2, x)"

# Numerical computation
calc-bridge "integral(sin(x), 0, pi)"
```
