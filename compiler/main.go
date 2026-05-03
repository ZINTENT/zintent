//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Intent struct {
	Styles string `json:"styles"`
	Hover  string `json:"hover,omitempty"`
}

type Registry struct {
	Intents      map[string]Intent            `json:"intents"`
	AriaMappings map[string]map[string]string `json:"aria_mappings"`
	Macros       map[string]string            `json:"macros"`
}

func loadRegistry(path string) (*Registry, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var reg Registry
	err = json.Unmarshal(data, &reg)
	return &reg, err
}

func injectAria(html string, ariaMappings map[string]map[string]string) string {
	modifiedHtml := html
	for cls, attrs := range ariaMappings {
		re := regexp.MustCompile(`class="([^"]*\b` + cls + `\b[^"]*)"`)
		modifiedHtml = re.ReplaceAllStringFunc(modifiedHtml, func(match string) string {
			classes := strings.TrimPrefix(strings.TrimSuffix(match, `"`), `class="`)
			attrStrings := []string{}
			for k, v := range attrs {
				attrStrings = append(attrStrings, fmt.Sprintf(`%s="%s"`, k, v))
			}
			return fmt.Sprintf(`class="%s" %s`, classes, strings.Join(attrStrings, " "))
		})
	}
	return modifiedHtml
}

func processMacros(html string, macros map[string]string) string {
	modifiedHtml := html
	for macro, replacement := range macros {
		// Simple tag replacement: <z-btn> -> <div class="...">
		openTag := fmt.Sprintf("<%s>", macro)
		closeTag := fmt.Sprintf("</%s>", macro)
		
		// Find replacement content placeholders
		// Assume replacement might look like "<div class='btn'>{{content}}</div>"
		if strings.Contains(replacement, "{{content}}") {
			re := regexp.MustCompile(fmt.Sprintf(`<%s>(.*?)</%s>`, macro, macro))
			modifiedHtml = re.ReplaceAllStringFunc(modifiedHtml, func(match string) string {
				content := re.FindStringSubmatch(match)[1]
				return strings.Replace(replacement, "{{content}}", content, 1)
			})
		} else {
			modifiedHtml = strings.ReplaceAll(modifiedHtml, openTag, strings.Split(replacement, ">")[0]+">")
			modifiedHtml = strings.ReplaceAll(modifiedHtml, closeTag, "</div>")
		}
	}
	return modifiedHtml
}

func compile(inputFile string, outputFile string, registry *Registry) {
	fmt.Printf("[%s] 🚀 ZINTENT: Compiling...\n", time.Now().Format("15:04:05"))

	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("❌ Error: Could not read input file: %v\n", err)
		return
	}

	html := string(content)

	// Phase 1: Macros
	html = processMacros(html, registry.Macros)

	// Phase 2: Accessibility
	html = injectAria(html, registry.AriaMappings)
	
	// Save the "processed" HTML for the browser to see
	// In a real dev server, we'd serve this in memory
	processedHtmlPath := filepath.Join(filepath.Dir(outputFile), "index.processed.html")
	_ = ioutil.WriteFile(processedHtmlPath, []byte(html), 0644)

	// Phase 3: Class Extraction
	usedClasses := make(map[string]bool)
	classRe := regexp.MustCompile(`class="([^"]+)"`)
	matches := classRe.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		classes := strings.Fields(match[1])
		for _, c := range classes {
			usedClasses[c] = true
		}
	}

	// Phase 4: CSS Generation
	var cssOutput strings.Builder
	cssOutput.WriteString("/* ZINTENT Generated Styles - Performance Build */\n\n")

	tokensPath := "core/tokens.css"
	if tokens, err := ioutil.ReadFile(tokensPath); err == nil {
		cssOutput.Write(tokens)
		cssOutput.WriteString("\n\n")
	}

	themesPath := "core/themes.css"
	if themes, err := ioutil.ReadFile(themesPath); err == nil {
		cssOutput.Write(themes)
		cssOutput.WriteString("\n\n")
	}

	for cls := range usedClasses {
		if intent, ok := registry.Intents[cls]; ok {
			cssOutput.WriteString(fmt.Sprintf(".%s {\n%s\n}\n\n", cls, intent.Styles))
			if intent.Hover != "" {
				cssOutput.WriteString(fmt.Sprintf(".%s:hover {\n%s\n}\n\n", cls, intent.Hover))
			}
		}
	}

	dir := filepath.Dir(outputFile)
	_ = os.MkdirAll(dir, os.ModePerm)
	_ = ioutil.WriteFile(outputFile, []byte(cssOutput.String()), 0644)

	fmt.Printf("✅ ZINTENT: Update complete.\n")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: zintent <input.html> -o <output.css> [--watch]")
		return
	}

	inputFile := args[0]
	outputFile := ""
	watch := false

	for i := 1; i < len(args); i++ {
		if args[i] == "-o" && i+1 < len(args) {
			outputFile = args[i+1]
			i++
		} else if args[i] == "--watch" {
			watch = true
		}
	}

	registry, err := loadRegistry("core/intent-registry.json")
	if err != nil {
		fmt.Printf("❌ Error loading registry: %v\n", err)
		return
	}

	if watch {
		fmt.Println("👀 ZINTENT: Watch mode active. Monitoring for changes...")
		lastSize := int64(0)
		for {
			info, err := os.Stat(inputFile)
			if err == nil {
				if info.Size() != lastSize {
					compile(inputFile, outputFile, registry)
					lastSize = info.Size()
				}
			}
			// Check registry changes too
			regInfo, err := os.Stat("core/intent-registry.json")
			if err == nil {
				// We don't track size here just yet for simplicity, but we could reload
				registry, _ = loadRegistry("core/intent-registry.json")
			}
			time.Sleep(500 * time.Millisecond)
		}
	} else {
		compile(inputFile, outputFile, registry)
	}
}
