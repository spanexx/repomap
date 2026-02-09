# Milestone 2.3: Quick Reference

## Summary

Web navigation and structure extraction tool. Agents navigate websites, extract accessibility trees, interact with elements, and extract structured data.

## Timeline

- **Duration:** 2–3 weeks
- **Tasks:** 25 sequential
- **Target Completion:** Week 13 of Phase 2

## Example Usage

```bash
# Navigate and screenshot
web-ray --url https://example.com --action "screenshot"

# Extract accessibility tree
web-ray --url https://github.com --action "extract_tree"

# Find and click element
web-ray --url https://example.com --action "click" --selector "login_button"

# Fill form and submit
web-ray --url https://example.com --action "fill" --selector "email" --value "test@example.com"

# Complex workflow
web-ray --url https://example.com \
  --actions "screenshot,click:login_button,fill:email:user@test.com,click:submit"
```

## Key Components

| Component | Path | Purpose |
|-----------|------|---------|
| Browser Wrapper | `pkg/browser/browser.go` | Playwright abstraction |
| Navigation | `pkg/browser/nav.go` | Page loading & waiting |
| Actions | `pkg/browser/actions.go` | Click, fill, select |
| DOM Extractor | `pkg/accessibility/dom.go` | DOM tree extraction |
| Tree Builder | `pkg/accessibility/tree.go` | Accessibility tree |
| Selector | `pkg/accessibility/selector.go` | Element selection |
| CLI Tool | `cmd/web-ray/main.go` | Command-line interface |

## Execution Strategy

1. **Browser Automation (Days 1–7)**
   - Evaluate and set up Playwright
   - Implement page navigation
   - Build interaction actions

2. **Structure Extraction (Days 8–14)**
   - Extract DOM and accessibility trees
   - Implement element descriptions
   - Build coordinate detection

3. **Interface & Testing (Days 15–21)**
   - CLI tool implementation
   - Query language
   - Integration testing

## Success Criteria

| Criterion | Target | Status |
|-----------|--------|--------|
| Navigation Success | >95% | — |
| Element Detection | >90% | — |
| Action Execution | >95% | — |
| Code Coverage | >80% | — |
| Page Load Time | <5s | — |

## Acceptance Criteria

- ✅ Navigates complex websites reliably
- ✅ Extracts accurate accessibility trees
- ✅ Detects and describes interactive elements
- ✅ Executes actions (click, fill, wait)
- ✅ Screenshots with element annotations
- ✅ >95% action success rate
- ✅ >80% code coverage
- ✅ Works on real websites

## Key Metrics

| Metric | Target |
|--------|--------|
| Page Load Time | <5s |
| Navigation Success | >95% |
| Element Detection | >90% |
| Action Success | >95% |
| Code Coverage | >80% |
