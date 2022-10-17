package jsonnettest

import (
	"strings"

	"github.com/fatih/color"
)

// ColorizeDiffOutput colorizes the diff output from myers.ComputeEdits
func ColorizeDiffOutput(diff string) string {
	newString := make([]string, 0)

	for _, l := range strings.Split(diff, "\n") {
		if strings.HasPrefix(l, "-") {

			newString = append(newString, color.New(color.FgRed).Sprint(l))
			continue
		}

		if strings.HasPrefix(l, "+") {
			newString = append(newString, color.New(color.FgGreen).Sprint(l))
			continue
		}

		if strings.HasPrefix(l, "@") {
			newString = append(newString, color.New(color.FgCyan).Sprint(l))
			continue
		}

		newString = append(newString, l)
	}

	return strings.Join(newString, "\n")
}
