package filter

import (
	"LingerAPI/pkg/spray"
	"log"
	"os"
	"strings"
)

// to read .txt filters
func ReadFilterFromFile(filename string) []string {
	data, err := os.ReadFile(filename)
	if err != nil {

		log.Println(spray.Rspray("Error reading file: "), err)
		return nil
	}

	lines := strings.Split(string(data), "\n")
	var filter []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			filter = append(filter, line)
		}
	}

	return filter
}
