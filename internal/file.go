package internal

import (
	"fmt"
	"os"

	"github.com/bmatcuk/doublestar/v4"
)

func FindPaths(globPatterns []string) ([]string, error) {
	paths := map[string]struct{}{}

	// Find all paths
	for _, pattern := range globPatterns {
		base, pattern := doublestar.SplitPattern(pattern)
		fsys := os.DirFS(base)
		matches, err := doublestar.Glob(fsys, pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during parsing %s, skip.\n", pattern)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return nil, err
		}
		for _, match := range matches {
			paths[base+"/"+match] = struct{}{}
		}
	}

	// Convert result into array
	result := make([]string, 0)
	for path := range paths {
		result = append(result, path)
	}

	return result, nil
}
