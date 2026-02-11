package analysis

import (
	"testing"
)

func TestDuplicationDetector(t *testing.T) {
	detector := NewDuplicationDetector()
	detector.MinTokens = 10 // Lower threshold for testing

	files := map[string][]byte{
		"file1.go": []byte(`
			package main
			import "fmt"
			func duplicateLogic() {
				fmt.Println("This is a duplicate block of code")
				fmt.Println("It should be detected by the analyzer")
				a := 1 + 2
				b := a * 3
			}
		`),
		"file2.go": []byte(`
			package other
			import "fmt"
			func someOtherFunction() {
				// comments should be ignored
				fmt.Println("This is a duplicate block of code")
				fmt.Println("It should be detected by the analyzer")
				a := 1 + 2
				b := a * 3
			}
		`),
		"file3.go": []byte(`
			package unique
			func uniqueLogic() {
				// unique content
				x := 100
				y := x / 10
			}
		`),
	}

	issues, err := detector.Analyze(files)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if len(issues) == 0 {
		t.Errorf("Expected duplicates to be found, got 0 issues")
	}

	found := false
	for _, issue := range issues {
		if issue.Type == "duplication" {
			found = true
			t.Logf("Found issue: %s", issue.Description)
		}
	}

	if !found {
		t.Errorf("Did not find duplication issue")
	}
}
