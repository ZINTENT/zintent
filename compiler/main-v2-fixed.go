//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Intent struct {
	Styles string `json:"styles"`
	Hover  string `json:"hover,omitempty"`
	Focus  string `json:"focus,omitempty"`
}

type Registry struct {
	Meta         map[string]interface{}       `json:"_meta"`
	Intents      map[string]interface{}       `json:"intents"`
	AriaMappings map[string]map[string]string `json:"aria_mappings"`
	Macros       map[string]string            `json:"macros"`
	Dependencies struct {
		CSSFiles []string `json:"css_files"`
	} `json:"dependencies"`
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

func flattenIntents(nested map[string]interface{}) map[string]Intent {
	flat := make(map[string]Intent)
	
	for category, value := range nested {
		if category == "_meta" {
			continue
		}
		
		switch v := value.(type) {
		case map[string]interface{}:
			// Check if this is a nested category (contains sub-categories)
			hasSubCategories := false
			for _, subValue := range v {
				if _, isMap := subValue.(map[string]interface{}); isMap {
					hasSubCategories = true
					break
				}
			}
			
			if hasSubCategories {
				// Recursively flatten
				subFlat := flattenIntents(v)
				for k, intent := range subFlat {
					flat[k] = intent
				}
			} else {
				// This is an actual intent definition
				intent := Intent{}
				if styles, ok := v["styles"].(string); ok {
					intent.Styles = styles
				}
				if hover, ok := v["hover"].(string); ok {
					intent.Hover = hover
				}
				if focus, ok := v["focus"].(string); ok {
					intent.Focus = focus
				}
				flat[category] = intent
			}
		}
	}
	
	return flat
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
		openTag := fmt.Sprintf("<%s>", macro)
		closeTag := fmt.Sprintf("</%s>", macro)
		
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
	fmt.Printf("[%s] ZINTENT v2.0: Compiling with Phase 1 features...\n", time.Now().Format("15:04:05"))

	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error: Could not read input file: %v\n", err)
		return
	}

	html := string(content)

	// Phase 1: Macros
	html = processMacros(html, registry.Macros)

	// Phase 2: Accessibility
	html = injectAria(html, registry.AriaMappings)

	// Save processed HTML
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

	// Flatten intents from nested structure
	flatIntents := flattenIntents(registry.Intents)

	// Phase 4: CSS Generation
	var cssOutput strings.Builder
	cssOutput.WriteString("/* ZINTENT v2.0 Generated Styles - Phase 1 Build */\n")
	cssOutput.WriteString("/* Features: Container-First, AI Tokens, Animations, Antigravity */\n\n")

	// Include dependency CSS files
	for _, cssFile := range registry.Dependencies.CSSFiles {
		if css, err := ioutil.ReadFile(cssFile); err == nil {
			cssOutput.WriteString(fmt.Sprintf("/* === %s === */\n", cssFile))
			cssOutput.Write(css)
			cssOutput.WriteString("\n\n")
		}
	}

	// Generate used intent classes
	cssOutput.WriteString("/* === Generated Intent Classes === */\n")
	
	// Sort classes to ensure deterministic output and correct override order
	var sortedClasses []string
	for cls := range usedClasses {
		sortedClasses = append(sortedClasses, cls)
	}
	sort.Strings(sortedClasses)
	
	generatedCount := 0
	for _, cls := range sortedClasses {
		if intent, ok := flatIntents[cls]; ok {
			selector := regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(cls, "\\$0")
			cssOutput.WriteString(fmt.Sprintf(".%s {\n%s\n}\n\n", selector, intent.Styles))
			if intent.Hover != "" {
				cssOutput.WriteString(fmt.Sprintf(".%s:hover {\n%s\n}\n\n", selector, intent.Hover))
			}
			if intent.Focus != "" {
				cssOutput.WriteString(fmt.Sprintf(".%s:focus {\n%s\n}\n\n", selector, intent.Focus))
			}
			generatedCount++
		}
	}

	// Add container query support detection
	cssOutput.WriteString(`
/* === Container Query Support Detection === */
@supports not (container-type: inline-size) {
  .container-responsive,
  .container-layout,
  .container-auto,
  .intent-responsive {
    width: 100%;
  }
}
`)

	// Add reduced motion support
	cssOutput.WriteString(`
/* === Reduced Motion Support === */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
`)

	dir := filepath.Dir(outputFile)
	_ = os.MkdirAll(dir, os.ModePerm)
	_ = ioutil.WriteFile(outputFile, []byte(cssOutput.String()), 0644)

	fmt.Printf("Build complete.\n")
	fmt.Printf("   Generated %d intent classes\n", generatedCount)
	fmt.Printf("   Processed HTML: %s\n", processedHtmlPath)
	fmt.Printf("   Output CSS: %s\n", outputFile)
	fmt.Printf("   Features: Container-First, AI Tokens, Animations, Antigravity\n")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: zintent <input.html> -o <output.css> [--watch]")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  zintent src/index.html -o dist/styles.css")
		fmt.Println("  zintent src/index.html -o dist/styles.css --watch")
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

	if outputFile == "" {
		outputFile = "dist/styles.css"
	}

	registry, err := loadRegistry("core/intent-registry-v2.json")
	if err != nil {
		fmt.Printf("Warning: Could not load v2 registry: %v\n", err)
		fmt.Println("   Falling back to v1 registry...")
		registry, err = loadRegistry("core/intent-registry.json")
		if err != nil {
			fmt.Printf("Error loading registry: %v\n", err)
			return
		}
	}

	if watch {
		fmt.Println("Watch mode active. Monitoring for changes...")
		lastSize := int64(0)
		for {
			info, err := os.Stat(inputFile)
			if err == nil {
				if info.Size() != lastSize {
					compile(inputFile, outputFile, registry)
					lastSize = info.Size()
				}
			}
			
			// Check registry changes
			_, err = os.Stat("core/intent-registry-v2.json")
			if err == nil {
				registry, _ = loadRegistry("core/intent-registry-v2.json")
			}
			
			time.Sleep(500 * time.Millisecond)
		}
	} else {
		compile(inputFile, outputFile, registry)
	}
}
