# Milestone 2.3: Web-Ray Foundation

## Objective

Create a tool for agents to navigate websites and extract structured information using headless browser automation and accessibility tree extraction.

## Scope

### In-Scope

1. **Headless Browser**
   - Playwright integration
   - JavaScript execution
   - Form filling and interaction
   - Cookie/session management

2. **Accessibility Tree**
   - Extract DOM as accessible structure
   - Element descriptions and roles
   - Interactive element detection
   - State tracking

3. **Structured Output**
   - JSON representation of page state
   - Element coordinate detection (bounding boxes)
   - Visual hierarchy preservation
   - Screenshot support with annotations

4. **Web-Ray Tool** (`cmd/web-ray/`)
   - CLI interface for page navigation
   - Query language for element selection
   - Action execution (click, fill, wait)
   - Result formatting

### Out-of-Scope

- Full page rendering
- CSS/styling extraction
- Advanced JavaScript framework support
- Mobile testing

## Deliverables

1. **Browser Automation Wrapper** (`pkg/browser/`)
2. **Accessibility Tree Extractor** (`pkg/accessibility/`)
3. **web-ray Tool** (`cmd/web-ray/`)
4. **Query Engine**
5. **Documentation & Examples**

## Success Criteria

- ✅ Navigates complex websites
- ✅ Extracts accurate accessibility tree
- ✅ Detects interactive elements
- ✅ Executes actions (click, fill, navigate)
- ✅ >80% code coverage
- ✅ Performance: page load <5 seconds

## Timeline

- **Duration:** 2–3 weeks
- **Tasks:** ~20–25

---

## Example Usage

```bash
web-ray --url https://example.com --action "screenshot"
# Output: JSON with page structure, interactive elements, screenshot

web-ray --url https://example.com --query "find login button"
# Output: Element description, coordinates, state

web-ray --url https://example.com --action "click" --selector "login_button" --then "screenshot"
# Output: Result after interaction
```

## Dependencies

- Phase 1 completion
- Playwright browser automation
- Accessibility API libraries
