package commands

import (
	"os"
	"path/filepath"
)

// GetOpenrtbTestingData loads all *.json files from baseDir/hostname and returns
// their contents as a slice of raw JSON strings for use with SimulationRTBRequester.
func GetOpenrtbTestingData(baseDir, hostname string) []string {
	pattern := filepath.Join(baseDir, hostname, "*.json")
	paths, err := filepath.Glob(pattern)
	if err != nil || len(paths) == 0 {
		return nil
	}
	results := make([]string, 0, len(paths))
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		results = append(results, string(data))
	}
	return results
}
