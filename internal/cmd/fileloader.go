package cmd

import (
	"bufio"
	"os"
	"strings"

	"github.com/AdguardTeam/golibs/log"
)

// loadRulesFromFile loads domain patterns from a CSV file.
// Lines starting with # are treated as comments and ignored.
// Empty lines are ignored.
// Returns the loaded rules or nil if file doesn't exist or is empty.
func loadRulesFromFile(filePath string) (rules []string, err error) {
	if filePath == "" {
		return nil, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug("cmd: rules file %s does not exist, skipping", filePath)
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		rules = append(rules, line)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	log.Info("cmd: loaded %d rules from %s", len(rules), filePath)
	return rules, nil
}
