package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	sourceDir = "i18n"
	output    = "cmd/wolmanager/i18n.generated.js"
)

var messagesPattern = regexp.MustCompile(`(?s)const\s+messages\s*=\s*(\{.*\})\s+as\s+const\s*;`)

func main() {
	files, err := filepath.Glob(filepath.Join(sourceDir, "*.ts"))
	if err != nil {
		fatal(err)
	}
	if len(files) == 0 {
		fatal(fmt.Errorf("no i18n files found in %s", sourceDir))
	}

	sort.Strings(files)
	all := map[string]map[string]string{}
	for _, file := range files {
		lang := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		messages, err := readMessages(file)
		if err != nil {
			fatal(err)
		}
		all[lang] = messages
	}
	if _, ok := all["en"]; !ok {
		fatal(fmt.Errorf("missing required fallback language: en"))
	}
	if err := validateKeys(all); err != nil {
		fatal(err)
	}

	payload, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		fatal(err)
	}
	content := "window.WOL_I18N = " + string(payload) + ";\n"
	if err := os.WriteFile(output, []byte(content), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("Generated %s from %d language files\n", output, len(files))
}

func readMessages(file string) (map[string]string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	match := messagesPattern.FindSubmatch(content)
	if match == nil {
		return nil, fmt.Errorf("%s: expected `const messages = {...} as const;`", file)
	}
	var messages map[string]string
	if err := json.Unmarshal(match[1], &messages); err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	return messages, nil
}

func validateKeys(all map[string]map[string]string) error {
	base := all["en"]
	for lang, messages := range all {
		for key := range base {
			if _, ok := messages[key]; !ok {
				return fmt.Errorf("%s is missing key %q", lang, key)
			}
		}
		for key := range messages {
			if _, ok := base[key]; !ok {
				return fmt.Errorf("%s has unknown key %q", lang, key)
			}
		}
	}
	return nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
