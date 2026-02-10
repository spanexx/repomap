package cli

import (
	"bytes"
	"flag"
	"fmt"
)

// GenerateHelp generates the help text for the application.
func (a *App) GenerateHelp() string {
	var b bytes.Buffer

	// 1. Name and Version
	fmt.Fprintf(&b, "%s v%s\n", a.Name, a.Version)
	if a.Description != "" {
		fmt.Fprintf(&b, "\n%s\n", a.Description)
	}

	// 2. Usage
	fmt.Fprintf(&b, "\nUsage:\n")
	if a.UsageText != "" {
		fmt.Fprintf(&b, "  %s\n", a.UsageText)
	} else {
		fmt.Fprintf(&b, "  %s [options]\n", a.Name)
	}

	// 3. Options (Flags)
	var flagsBuf bytes.Buffer
	a.flagSet.VisitAll(func(f *flag.Flag) {
		s := fmt.Sprintf("  --%s", f.Name)
		usage := f.Usage

		// Append default value if it's not empty
		if f.DefValue != "" {
			// For string flags, the default value is already unquoted text.
			// Ideally we would distinguish types to format nicely (e.g. quote strings)
			// but flag package stores everything as string.
			// We'll mimic standard format: (default: val)
			usage += fmt.Sprintf(" (default: %s)", f.DefValue)
		}

		// Simple alignment
		if len(s) < 25 {
			fmt.Fprintf(&flagsBuf, "%-25s %s\n", s, usage)
		} else {
			fmt.Fprintf(&flagsBuf, "%s\n%s%s\n", s, "                           ", usage)
		}
	})

	if flagsBuf.Len() > 0 {
		fmt.Fprintf(&b, "\nOptions:\n")
		b.Write(flagsBuf.Bytes())
	}

	// 4. Examples
	if len(a.Examples) > 0 {
		fmt.Fprintf(&b, "\nExamples:\n")
		for _, ex := range a.Examples {
			fmt.Fprintf(&b, "  %s\n", ex)
		}
	}

	return b.String()
}
