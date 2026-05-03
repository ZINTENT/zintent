package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Pre-compiled global regex for extreme performance (Phase 10 Hardening)
var (
	classRegex       = regexp.MustCompile(`(?:class|className)\s*=\s*["']([^"']+)["']`)
	classExprRegex   = regexp.MustCompile("(?:class|className)\\s*=\\s*\\{[\"'`].+?[\"'`]\\}")
	vueClassRegex    = regexp.MustCompile(`:class\s*=\s*["']([^"']+)["']`)
	// Stronger TSX/JS parsing: Match object literals { 'cls-name': true } and template literals
	jsObjectRegex    = regexp.MustCompile("[\"'`]+([a-zA-Z0-9@:/_-]+)[\"'`]+\\s*:\\s*(?:true|false|!!|[^,}]+)")
	templateLitRegex = regexp.MustCompile("`([^`]+)`")
	placeholderRegex = regexp.MustCompile(`{{([a-zA-Z0-9_-]+)}}`)
	attrRegex        = regexp.MustCompile(`([a-zA-Z0-9-]+)=["']([^"']*)["']`)
	slotRegex        = regexp.MustCompile(`(?s)<z-slot\s+name=["']([^"']*)["']>(.*?)</z-slot>`)
	variableRegex    = regexp.MustCompile(`(--zi-[a-zA-Z0-9-]+):\s*([^;]+);`)
)

type Intent struct {
	Styles string            `json:"styles"`
	Hover  string            `json:"hover,omitempty"`
	Focus  string            `json:"focus,omitempty"`
	Nested map[string]Intent `json:"nested,omitempty"`
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

func mergeRegistry(main *Registry, plugin *Registry) {
	// Merge Meta
	if plugin.Meta != nil {
		if main.Meta == nil {
			main.Meta = make(map[string]interface{})
		}
		for k, v := range plugin.Meta {
			main.Meta[k] = v
		}
	}
	// Merge Intents
	if plugin.Intents != nil {
		if main.Intents == nil {
			main.Intents = make(map[string]interface{})
		}
		for k, v := range plugin.Intents {
			main.Intents[k] = v
		}
	}
	// Merge AriaMappings
	if plugin.AriaMappings != nil {
		if main.AriaMappings == nil {
			main.AriaMappings = make(map[string]map[string]string)
		}
		for k, v := range plugin.AriaMappings {
			main.AriaMappings[k] = v
		}
	}
	// Merge Macros
	if plugin.Macros != nil {
		if main.Macros == nil {
			main.Macros = make(map[string]string)
		}
		for k, v := range plugin.Macros {
			main.Macros[k] = v
		}
	}
	// Merge Dependencies (append unique)
	for _, dep := range plugin.Dependencies.CSSFiles {
		found := false
		for _, mainDep := range main.Dependencies.CSSFiles {
			if mainDep == dep {
				found = true
				break
			}
		}
		if !found {
			main.Dependencies.CSSFiles = append(main.Dependencies.CSSFiles, dep)
		}
	}
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
				} else if css, ok := v["css"].(string); ok {
					intent.Styles = css
				}
				if hover, ok := v["hover"].(string); ok {
					intent.Hover = hover
				}
				if focus, ok := v["focus"].(string); ok {
					intent.Focus = focus
				}
				if nest, ok := v["nested"].(map[string]interface{}); ok {
					intent.Nested = make(map[string]Intent)
					for query, subInt := range nest {
						if subVal, ok := subInt.(map[string]interface{}); ok {
							si := Intent{}
							if s, ok := subVal["styles"].(string); ok {
								si.Styles = s
							} else if c, ok := subVal["css"].(string); ok {
								si.Styles = c
							}
							intent.Nested[query] = si
						}
					}
				}
				flat[category] = intent
			}
		}
	}

	return flat
}

// Phase 6 Experiment: CSS-in-Go Internal Intents
// These are defined in Go source to avoid disk I/O and enable ultra-fast "reset" styles.
func getInternalIntents() map[string]Intent {
	return map[string]Intent{
		"zi-reset": {
			Styles: "margin: 0; padding: 0; box-sizing: border-box; -webkit-font-smoothing: antialiased;",
		},
		"zi-base-body": {
			Styles: "min-height: 100vh; overflow-x: hidden; scroll-behavior: smooth;",
		},
		"zi-antigravity-root": {
			Styles: "display: grid; min-height: 100vh; grid-template-rows: auto 1fr auto;",
		},
	}
}

func injectAria(html string, ariaMappings map[string]map[string]string) string {
	modifiedHtml := html
	// Map to track tags we've already processed for a specific attribute to avoid duplicates
	// This is a simple implementation; a full HTML parser would be more robust but slower.

	for cls, attrs := range ariaMappings {
		// Matches tags that have the specific class, handling different quote types
		re := regexp.MustCompile(`(<[a-zA-Z0-9]+\s+[^>]*class=["']([^"']*\b` + cls + `\b[^"']*)["'][^>]*>)`)

		modifiedHtml = re.ReplaceAllStringFunc(modifiedHtml, func(tag string) string {
			for attr, val := range attrs {
				// Don't inject if the attribute already exists in the tag
				attrPattern := regexp.MustCompile(`\s` + attr + `=["']`)
				if !attrPattern.MatchString(tag) {
					// Inject before the closing bracket
					tag = strings.TrimSuffix(tag, ">") + fmt.Sprintf(` %s="%s">`, attr, val)
				}
			}
			return tag
		})
	}
	return modifiedHtml
}

func lintIntents(usedClasses map[string]bool, flatIntents map[string]Intent) {
	fmt.Printf("[%s] ZINTENT: Running intent-specific linting...\n", time.Now().Format("15:04:05"))

	// 1. Conflicting Layout Intents (heuristic based)
	layoutConflicts := []struct {
		a, b string
		msg  string
	}{
		{"intent-center", "intent-cluster", "Conflicting display modes: display:grid (center) vs display:flex (cluster)"},
		{"intent-stack", "intent-cluster", "Conflicting flex directions: column (stack) vs row (cluster)"},
		{"intent-full-height", "intent-auto-grid", "Unusual combination: Fixed full-height might clash with auto-grid content overflow"},
	}

	// This is a global check for the entire file.
	// In the future, we could check per-element, but this provides good broad warnings.
	for _, rule := range layoutConflicts {
		if usedClasses[rule.a] && usedClasses[rule.b] {
			fmt.Printf("  [LINT WARNING] %s\n", rule.msg)
		}
	}

	// 2. Dead Intent Detector
	deadCount := 0
	for name := range flatIntents {
		if !usedClasses[name] && !strings.HasPrefix(name, "zi-") {
			deadCount++
		}
	}
	if deadCount > 0 {
		fmt.Printf("  [LINT INFO] %d unused intents found in registry (tree-shaken from output).\n", deadCount)
	}
}

func extractPlaceholders(macroBody string) []string {
	re := regexp.MustCompile(`{{([a-zA-Z0-9_-]+)}}`)
	matches := re.FindAllStringSubmatch(macroBody, -1)
	seen := make(map[string]bool)
	var placeholders []string
	for _, m := range matches {
		if m[1] != "content" && m[1] != "attributes" && !seen[m[1]] {
			placeholders = append(placeholders, m[1])
			seen[m[1]] = true
		}
	}
	return placeholders
}

func generateCustomData(registry *Registry, outputPath string) {
	fmt.Printf("[%s] ZINTENT: Generating VS Code CustomData...\n", time.Now().Format("15:04:05"))

	tags := []map[string]interface{}{}

	// Sort macros for deterministic output
	var macroKeys []string
	for k := range registry.Macros {
		macroKeys = append(macroKeys, k)
	}
	sort.Strings(macroKeys)

	for _, name := range macroKeys {
		placeholders := extractPlaceholders(registry.Macros[name])
		attributes := []map[string]interface{}{}

		for _, p := range placeholders {
			attributes = append(attributes, map[string]interface{}{
				"name":        p,
				"description": fmt.Sprintf("ZINTENT Macro Property: %s", p),
			})
		}

		tags = append(tags, map[string]interface{}{
			"name":        name,
			"description": fmt.Sprintf("ZINTENT High-Fidelity Macro: %s", name),
			"attributes":  attributes,
		})
	}

	customData := map[string]interface{}{
		"version": 1.1,
		"tags":    tags,
	}

	data, _ := json.MarshalIndent(customData, "", "  ")
	_ = ioutil.WriteFile(outputPath, data, 0644)
	fmt.Printf("  [DONE] CustomData saved to %s (Intellisense Active for %d macros)\n", outputPath, len(tags))
}

func generateSchema(registry *Registry, outputPath string) {
	fmt.Printf("[%s] ZINTENT: Generating IDE JSON Schema...\n", time.Now().Format("15:04:05"))

	schema := map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title":   "ZINTENT Intent Registry Schema",
		"type":    "object",
		"properties": map[string]interface{}{
			"intents": map[string]interface{}{
				"type": "object",
				"patternProperties": map[string]interface{}{
					"^[a-z0-9_-]+$": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"styles": map[string]interface{}{"type": "string"},
							"hover":  map[string]interface{}{"type": "string"},
							"focus":  map[string]interface{}{"type": "string"},
							"nested": map[string]interface{}{"type": "object"},
						},
						"required":             []string{"styles"},
						"additionalProperties": false,
					},
				},
			},
			"macros": map[string]interface{}{
				"type": "object",
				"patternProperties": map[string]interface{}{
					"^[z]-[a-z0-9_-]+$": map[string]interface{}{
						"type":        "string",
						"description": "ZINTENT Macro Definition (HTML Template)",
					},
				},
			},
		},
	}

	data, _ := json.MarshalIndent(schema, "", "  ")
	_ = ioutil.WriteFile(outputPath, data, 0644)
	fmt.Printf("  [DONE] Schema saved to %s\n", outputPath)
}

func generateVariableMap(registry *Registry, outputPath string) {
	fmt.Printf("[%s] ZINTENT: Generating CSS Variable Map...\n", time.Now().Format("15:04:05"))

	varMap := make(map[string]string)

	for _, cssFile := range registry.Dependencies.CSSFiles {
		content, err := ioutil.ReadFile(cssFile)
		if err != nil {
			continue
		}

		re := regexp.MustCompile(`(--zi-[a-zA-Z0-9-]+):\s*([^;]+);`)
		matches := re.FindAllStringSubmatch(string(content), -1)
		for _, m := range matches {
			varMap[m[1]] = strings.TrimSpace(m[2])
		}
	}

	data, _ := json.MarshalIndent(varMap, "", "  ")
	_ = ioutil.WriteFile(outputPath, data, 0644)
	fmt.Printf("  [DONE] Variable map saved to %s (%d tokens)\n", outputPath, len(varMap))
}

func initProject() {
	fmt.Printf("[%s] ZINTENT: Initializing new project...\n", time.Now().Format("15:04:05"))

	_ = os.MkdirAll("src", 0755)
	_ = os.MkdirAll("dist", 0755)

	htmlContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My ZINTENT Project</title>
    <link rel="stylesheet" href="../dist/styles.css">
</head>
<body class="zi-reset zi-base-body">
    <z-nav>
        <div class="zi-text-xl zi-font-bold">My Project</div>
    </z-nav>
    <main class="zi-container zi-py-12">
        <z-card>
            <h1>Welcome to ZINTENT</h1>
            <p>Ready to build something amazing?</p>
            <z-btn variant="primary">Get Started</z-btn>
        </z-card>
    </main>
</body>
</html>`

	_ = ioutil.WriteFile("src/index.html", []byte(htmlContent), 0644)
	fmt.Println("  [DONE] Created src/index.html")
	fmt.Println("  [TIP] Run 'go run ./compiler --input src/index.html' to build.")
}

func processMacros(html string, macros map[string]string) string {
	modifiedHtml := html
	for macro, replacement := range macros {
		// Matches <macro-name attr="val" ...>content</macro-name>
		// or <macro-name attr="val" ... />
		re := regexp.MustCompile(`(?s)<` + macro + `(\s+[^>]*)?>(.*?)</` + macro + `>|<` + macro + `(\s+[^>]*)?/>`)

		modifiedHtml = re.ReplaceAllStringFunc(modifiedHtml, func(match string) string {
			sub := re.FindStringSubmatch(match)
			attrsRaw := sub[1] + sub[3] // One of them will be empty
			content := sub[2]

			result := replacement

			// Handle attributes (props)
			attrMap := make(map[string]string)
			attrRe := regexp.MustCompile(`([a-zA-Z0-9-]+)=["']([^"']*)["']`)
			attrMatches := attrRe.FindAllStringSubmatch(attrsRaw, -1)
			for _, m := range attrMatches {
				attrMap[m[1]] = m[2]
				result = strings.ReplaceAll(result, fmt.Sprintf("{{%s}}", m[1]), m[2])
			}

			// Handle slots
			if strings.Contains(result, "{{content}}") {
				result = strings.ReplaceAll(result, "{{content}}", content)
			}

			// Handle named slots: <z-slot name="title">...</z-slot>
			slotRe := regexp.MustCompile(`(?s)<z-slot\s+name=["']([^"']*)["']>(.*?)</z-slot>`)
			slotMatches := slotRe.FindAllStringSubmatch(content, -1)
			for _, m := range slotMatches {
				result = strings.ReplaceAll(result, fmt.Sprintf("{{slot_%s}}", m[1]), m[2])
			}

			// Clean up unused placeholders
			placeholderRe := regexp.MustCompile(`{{[a-zA-Z0-9_-]+}}`)
			result = placeholderRe.ReplaceAllString(result, "")

			return result
		})
	}
	return modifiedHtml
}

func addClassesFromValue(raw string, usedClasses map[string]bool) {
	for _, token := range strings.Fields(raw) {
		clean := strings.TrimSpace(token)
		clean = strings.Trim(clean, "\"'`{}(),")
		if clean == "" {
			continue
		}
		if strings.Contains(clean, "${") || strings.Contains(clean, "{{") {
			continue
		}
		if strings.HasPrefix(clean, "text-") && strings.Contains(clean, "$") {
			continue
		}
		usedClasses[clean] = true
	}
}

func extractClassesFromContent(content string, usedClasses map[string]bool) {
	for _, match := range classRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			addClassesFromValue(match[1], usedClasses)
		}
	}
	for _, match := range classExprRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			addClassesFromValue(match[1], usedClasses)
		}
	}
	for _, match := range vueClassRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			addClassesFromValue(match[1], usedClasses)
		}
	}
	// Stronger TSX/JS matches
	for _, match := range jsObjectRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			addClassesFromValue(match[1], usedClasses)
		}
	}
	for _, match := range templateLitRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			addClassesFromValue(match[1], usedClasses)
		}
	}
}

func collectContentFiles(path string, exts map[string]bool) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	files := []string{}
	if !info.IsDir() {
		ext := strings.ToLower(filepath.Ext(path))
		if exts[ext] {
			return []string{path}, nil
		}
		return files, nil
	}

	err = filepath.Walk(path, func(filePath string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if fi.IsDir() {
			base := strings.ToLower(fi.Name())
			if base == "node_modules" || base == "dist" || strings.HasPrefix(base, ".git") {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(filePath))
		if exts[ext] {
			files = append(files, filePath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

func classSetHasPrefix(usedClasses map[string]bool, prefix string) bool {
	for cls := range usedClasses {
		if strings.HasPrefix(cls, prefix) {
			return true
		}
	}
	return false
}

// minimalNeedsAntigravityLayout returns true if any used class likely requires antigravity-layouts.css.
func minimalNeedsAntigravityLayout(usedClasses map[string]bool) bool {
	for cls := range usedClasses {
		if strings.HasPrefix(cls, "intent-") || strings.HasPrefix(cls, "antigravity-") {
			return true
		}
	}
	return false
}

// minimalNeedsAnimations returns true if animations.css is needed under minimal preset.
func minimalNeedsAnimations(usedClasses map[string]bool) bool {
	if classSetHasPrefix(usedClasses, "transition-") ||
		classSetHasPrefix(usedClasses, "animate-") ||
		classSetHasPrefix(usedClasses, "stagger-") ||
		classSetHasPrefix(usedClasses, "delay-") ||
		classSetHasPrefix(usedClasses, "duration-") ||
		classSetHasPrefix(usedClasses, "fill-") ||
		classSetHasPrefix(usedClasses, "iteration-") ||
		usedClasses["stagger-children"] ||
		usedClasses["intent-card-hover"] ||
		usedClasses["intent-press"] ||
		usedClasses["intent-link"] ||
		usedClasses["animate-paused"] ||
		usedClasses["animate-running"] {
		return true
	}
	return false
}

func shouldIncludeDependency(cssFile string, usedClasses map[string]bool, preset string) bool {
	file := strings.ReplaceAll(strings.ToLower(cssFile), "\\", "/")

	// Legacy compatibility shims: ship only in `full` preset.
	if strings.HasSuffix(file, "core/cross-browser.css") {
		return preset == "full"
	}

	// `full` and `core` include all standard dependency CSS (core omits cross-browser above).
	if preset == "full" || preset == "core" {
		return true
	}

	// `minimal`: tree-shake heavy modules when unused.
	if strings.HasSuffix(file, "core/tokens.css") || strings.HasSuffix(file, "core/themes.css") || strings.HasSuffix(file, "core/a11y.css") {
		return true
	}

	if strings.HasSuffix(file, "core/container-queries.css") {
		return usedClasses["container-responsive"] ||
			usedClasses["container-layout"] ||
			usedClasses["container-auto"] ||
			usedClasses["container-fixed"] ||
			usedClasses["intent-responsive"] ||
			classSetHasPrefix(usedClasses, "@")
	}

	if strings.HasSuffix(file, "core/animations.css") {
		return minimalNeedsAnimations(usedClasses)
	}

	if strings.HasSuffix(file, "core/antigravity-layouts.css") {
		return minimalNeedsAntigravityLayout(usedClasses)
	}

	return true
}

func compile(inputFile string, outputFile string, registry *Registry, analyze bool, budgetKB int, scope string, contentPath string, preset string, scanner string) {
	fmt.Printf("[%s] ZINTENT v2.1.0: Compiling...\n", time.Now().Format("15:04:05"))

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
	processedBase := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	processedHtmlPath := filepath.Join(filepath.Dir(outputFile), processedBase+".processed.html")
	_ = ioutil.WriteFile(processedHtmlPath, []byte(html), 0644)

	// Phase 3: Class Extraction
	usedClasses := make(map[string]bool)
	extractClassesFromPath(inputFile, html, usedClasses, scanner, true)

	externalFilesScanned := 0
	if contentPath != "" {
		extensions := map[string]bool{
			".html": true, ".htm": true, ".php": true,
			".js": true, ".jsx": true, ".ts": true, ".tsx": true,
			".vue": true, ".svelte": true,
		}
		files, scanErr := collectContentFiles(contentPath, extensions)
		if scanErr != nil {
			fmt.Printf("Warning: content scan skipped (%v)\n", scanErr)
		} else {
			for _, file := range files {
				bytes, fileErr := ioutil.ReadFile(file)
				if fileErr != nil {
					continue
				}
				extractClassesFromPath(file, string(bytes), usedClasses, scanner, false)
			}
			externalFilesScanned = len(files)
		}
	}

	// Flatten intents from nested structure
	flatIntents := flattenIntents(registry.Intents)

	// Phase 6 Experiment: Merge Internal CSS-in-Go Intents
	internal := getInternalIntents()
	for k, v := range internal {
		if _, exists := flatIntents[k]; !exists {
			flatIntents[k] = v
		}
	}

	// Phase 6: Intent-specific Linting
	lintIntents(usedClasses, flatIntents)

	// Phase 9: Interaction Audit Engine
	uxFindings := AuditUX(html, inputFile)
	if len(uxFindings) > 0 {
		fmt.Printf("[%s] ZINTENT: Interaction Audit Results:\n", time.Now().Format("15:04:05"))
		for _, f := range uxFindings {
			fmt.Printf("  [UX WARNING] %s\n", f)
		}
	}

	// Phase 4: CSS Generation
	var cssOutput strings.Builder
	cssOutput.WriteString("/* ZINTENT v2.1.0 Generated Styles */\n")
	cssOutput.WriteString("/* Features: Container-First, Tokens, Animations, Antigravity */\n\n")

	// Predictable cascade + optional isolation (Modern CSS)
	// Order matters: tokens -> themes -> core -> generated intents
	cssOutput.WriteString("@layer zi-tokens, zi-themes, zi-core, zi-intents;\n\n")

	// Include dependency CSS files
	for _, cssFile := range registry.Dependencies.CSSFiles {
		if !shouldIncludeDependency(cssFile, usedClasses, preset) {
			continue
		}
		if css, err := ioutil.ReadFile(cssFile); err == nil {
			layer := "zi-core"
			if strings.HasSuffix(cssFile, "core/tokens.css") {
				layer = "zi-tokens"
			} else if strings.HasSuffix(cssFile, "core/themes.css") {
				layer = "zi-themes"
			}

			cssOutput.WriteString(fmt.Sprintf("/* === %s === */\n", cssFile))
			cssOutput.WriteString(fmt.Sprintf("@layer %s {\n", layer))
			cssOutput.Write(css)
			cssOutput.WriteString("\n}\n\n")
		}
	}

	// Generate used intent classes
	cssOutput.WriteString("/* === Generated Intent Classes === */\n")
	cssOutput.WriteString("@layer zi-intents {\n")

	// Phase 8: Deduplication Engine (Optimization)
	// We group identical style blocks to shrink output size
	styleGroups := make(map[string][]string)

	// Sort classes to ensure deterministic output
	var sortedClasses []string
	for cls := range usedClasses {
		sortedClasses = append(sortedClasses, cls)
	}
	sort.Strings(sortedClasses)

	generatedCount := 0
	for _, cls := range sortedClasses {
		if intent, ok := flatIntents[cls]; ok {
			// Generate full style block including hover/focus/nested
			fullStyle := intent.Styles
			if intent.Hover != "" {
				fullStyle += "\n&:hover {\n" + intent.Hover + "\n}"
			}
			if intent.Focus != "" {
				fullStyle += "\n&:focus {\n" + intent.Focus + "\n}"
			}
			if len(intent.Nested) > 0 {
				var nestedQueries []string
				for q := range intent.Nested {
					nestedQueries = append(nestedQueries, q)
				}
				sort.Strings(nestedQueries)
				for _, q := range nestedQueries {
					fullStyle += "\n" + q + " {\n" + intent.Nested[q].Styles + "\n}"
				}
			}

			styleGroups[fullStyle] = append(styleGroups[fullStyle], cls)
			generatedCount++
		}
	}

	// Output grouped CSS
	var styleKeys []string
	for k := range styleGroups {
		styleKeys = append(styleKeys, k)
	}
	sort.Strings(styleKeys)

	for _, style := range styleKeys {
		selectors := styleGroups[style]
		var selectorList []string
		for _, s := range selectors {
			escaped := regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(s, "\\$0")
			if scope != "" {
				escaped = scope + "-" + escaped
			}
			selectorList = append(selectorList, "."+escaped)
		}

		cssOutput.WriteString(fmt.Sprintf("  %s {\n%s\n  }\n\n", strings.Join(selectorList, ", "), indentCSS(style, "    ")))
	}

	cssOutput.WriteString("}\n\n")

	// Add container query support detection
	cssOutput.WriteString(`
/* === Container Query Support Detection === */
@layer zi-core {
  @supports not (container-type: inline-size) {
    .container-responsive,
    .container-layout,
    .container-auto,
    .intent-responsive {
      width: 100%;
    }
  }
}
`)

	// Add reduced motion support
	cssOutput.WriteString(`
/* === Reduced Motion Support === */
@layer zi-core {
  @media (prefers-reduced-motion: reduce) {
    *, *::before, *::after {
      animation-duration: 0.01ms !important;
      animation-iteration-count: 1 !important;
      transition-duration: 0.01ms !important;
    }
  }
}
`)

	dir := filepath.Dir(outputFile)
	_ = os.MkdirAll(dir, os.ModePerm)
	_ = ioutil.WriteFile(outputFile, []byte(cssOutput.String()), 0644)

	var outputSizeBytes int64 = 0
	if info, err := os.Stat(outputFile); err == nil {
		outputSizeBytes = info.Size()
	}

	fmt.Printf("Build complete.\n")
	fmt.Printf("   Generated %d intent classes\n", generatedCount)
	fmt.Printf("   Processed HTML: %s\n", processedHtmlPath)
	fmt.Printf("   Output CSS: %s\n", outputFile)
	if contentPath != "" {
		fmt.Printf("   Scanned content files: %d\n", externalFilesScanned)
	}
	fmt.Printf("   Preset: %s\n", preset)
	fmt.Printf("   Scanner: %s\n", scanner)
	fmt.Printf("   Features: Container-First, AI Tokens, Animations, Antigravity\n")

	if analyze {
		fmt.Printf("Analysis:\n")
		fmt.Printf("   Used classes in HTML: %d\n", len(usedClasses))
		fmt.Printf("   Output size: %.2f KB\n", float64(outputSizeBytes)/1024.0)
	}

	if budgetKB > 0 {
		if float64(outputSizeBytes)/1024.0 > float64(budgetKB) {
			fmt.Printf("Budget exceeded: output is %.2f KB (budget %d KB)\n", float64(outputSizeBytes)/1024.0, budgetKB)
			os.Exit(2)
		}
	}
}

func indentCSS(css string, indent string) string {
	css = strings.TrimSpace(css)
	if css == "" {
		return ""
	}

	// Handle variable replacement if any
	// (Future Phase: --zi-var injection)

	parts := strings.Split(css, "\n")
	for i, p := range parts {
		parts[i] = indent + strings.TrimRight(p, " \t")
	}
	return strings.Join(parts, "\n")
}

func main() {
	inputFile := ""
	outputFile := ""
	registryFile := "core/intent-registry-v2.json"
	watch := false
	analyze := false
	budgetKB := 0
	scope := ""
	contentPath := ""
	preset := "full"
	scanner := "regex"
	schemaPath := ""
	docsPath := ""
	varsPath := ""
	pluginFiles := []string{}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "--input":
			if i+1 < len(os.Args) {
				inputFile = os.Args[i+1]
				i++
			}
		case "--output":
			if i+1 < len(os.Args) {
				outputFile = os.Args[i+1]
				i++
			}
		case "-o":
			if i+1 < len(os.Args) {
				outputFile = os.Args[i+1]
				i++
			}
		case "--registry":
			if i+1 < len(os.Args) {
				registryFile = os.Args[i+1]
				i++
			}
		case "--analyze":
			analyze = true
		case "--watch":
			watch = true
		case "--budget-kb":
			if i+1 < len(os.Args) {
				budgetKB, _ = strconv.Atoi(os.Args[i+1])
				i++
			}
		case "--scope":
			if i+1 < len(os.Args) {
				scope = os.Args[i+1]
				i++
			}
		case "--content":
			if i+1 < len(os.Args) {
				contentPath = os.Args[i+1]
				i++
			}
		case "--preset":
			if i+1 < len(os.Args) {
				preset = strings.ToLower(os.Args[i+1])
				i++
			}
		case "--scanner":
			if i+1 < len(os.Args) {
				scanner = strings.ToLower(os.Args[i+1])
				i++
			}
		case "--schema":
			if i+1 < len(os.Args) {
				schemaPath = os.Args[i+1]
				i++
			}
		case "--vars":
			if i+1 < len(os.Args) {
				varsPath = os.Args[i+1]
				i++
			}
		case "--docs":
			if i+1 < len(os.Args) {
				docsPath = os.Args[i+1]
				i++
			}
		case "--plugin":
			if i+1 < len(os.Args) {
				pluginFiles = append(pluginFiles, os.Args[i+1])
				i++
			}
		case "--init":
			initProject()
			return
		default:
			// Backward compatibility: support positional input/output args.
			if strings.HasPrefix(arg, "-") {
				continue
			}
			if inputFile == "" {
				inputFile = arg
			} else if outputFile == "" {
				outputFile = arg
			}
		}
	}

	if inputFile == "" && schemaPath == "" && docsPath == "" && varsPath == "" {
		fmt.Println("Usage: go run ./compiler --input <file> --output <file> [options]")
		fmt.Println("Options:")
		fmt.Println("  --registry <path>   Path to intent registry JSON")
		fmt.Println("  --output, -o <path> Output CSS file path")
		fmt.Println("  --analyze          Output bundle size analysis")
		fmt.Println("  --budget-kb <n>    Set a bundle size budget in KB")
		fmt.Println("  --watch            Watch for changes and recompile")
		fmt.Println("  --scope <prefix>   Scoping prefix for generated selectors")
		fmt.Println("  --content <path>   Scan classes from extra files/folders")
		fmt.Println("  --preset <mode>    full (default), core, or minimal")
		fmt.Println("  --scanner <mode>   regex (default) or parser (HTML DOM for .html/.htm/.php)")
		fmt.Println("  --schema <path>    Generate JSON schema for IDE autocomplete")
		fmt.Println("  --vars <path>      Generate CSS variable map JSON")
		fmt.Println("  --docs <path>      Generate interactive ZINTENT Explorer documentation")
		fmt.Println("  --plugin <path>    Load extra intents/macros from a plugin JSON file")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  Production (lean):  go run ./compiler -i src/index.html -o dist/zintent.css --preset minimal --content src")
		fmt.Println("  SPA (React/Vite):   go run ./compiler -i index.html -o dist/zintent.css --preset core --content src")
		fmt.Println("  Laravel/Blade:      go run ./compiler -i resources/views/app.blade.php -o public/css/zintent.css --scanner parser --content resources/views")
		fmt.Println("  Watch dev:          go run ./compiler --watch -i src/page.html -o dist/out.css --preset minimal")
		return
	}

	reg, err := loadRegistry(registryFile)
	if err != nil {
		fmt.Printf("Error: Could not load registry: %v\n", err)
		return
	}

	// Load and merge plugins
	for _, pluginFile := range pluginFiles {
		plugin, err := loadRegistry(pluginFile)
		if err != nil {
			fmt.Printf("Warning: Could not load plugin %s: %v\n", pluginFile, err)
			continue
		}
		mergeRegistry(reg, plugin)
		fmt.Printf("[%s] ZINTENT: Merged plugin %s\n", time.Now().Format("15:04:05"), pluginFile)
	}

	if schemaPath != "" {
		generateSchema(reg, schemaPath)
		if inputFile == "" && docsPath == "" && varsPath == "" {
			return
		}
	}

	if varsPath != "" {
		generateVariableMap(reg, varsPath)
		if inputFile == "" && docsPath == "" {
			return
		}
	}

	if docsPath != "" {
		generateExplorer(reg, docsPath)
		if inputFile == "" {
			return
		}
	}

	if outputFile == "" {
		outputFile = "dist/styles.css"
	}
	if preset != "full" && preset != "minimal" && preset != "core" {
		fmt.Printf("Invalid preset '%s'. Use 'full', 'core', or 'minimal'.\n", preset)
		return
	}
	if scanner != "regex" && scanner != "parser" {
		fmt.Printf("Invalid scanner '%s'. Use 'regex' or 'parser'.\n", scanner)
		return
	}

	if watch {
		fmt.Println("Watch mode active. Monitoring for changes...")
		lastInputMod := time.Time{}
		lastRegistryMod := time.Time{}
		for {
			info, err := os.Stat(inputFile)
			if err == nil {
				if info.ModTime().After(lastInputMod) {
					compile(inputFile, outputFile, reg, analyze, budgetKB, scope, contentPath, preset, scanner)
					lastInputMod = info.ModTime()
				}
			}

			// Check registry changes
			if regInfo, regErr := os.Stat(registryFile); regErr == nil {
				if regInfo.ModTime().After(lastRegistryMod) {
					if nextReg, loadErr := loadRegistry(registryFile); loadErr == nil {
						reg = nextReg
						compile(inputFile, outputFile, reg, analyze, budgetKB, scope, contentPath, preset, scanner)
						lastRegistryMod = regInfo.ModTime()
					}
				}
			}

			time.Sleep(500 * time.Millisecond)
		}
	} else {
		compile(inputFile, outputFile, reg, analyze, budgetKB, scope, contentPath, preset, scanner)
	}
}

func generateExplorer(reg *Registry, outputPath string) {
	fmt.Printf("[%s] ZINTENT: Generating Interactive Explorer...\n", time.Now().Format("15:04:05"))

	flatIntents := flattenIntents(reg.Intents)
	intentData, _ := json.Marshal(flatIntents)
	macroData, _ := json.Marshal(reg.Macros)

	html := `<!DOCTYPE html>
<html lang="en" data-theme="nordic">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ZINTENT Explorer</title>
    <link rel="stylesheet" href="styles.css">
    <style>
        :root { --sidebar-w: 300px; }
        body { margin: 0; display: flex; height: 100vh; background: var(--zi-bg-base); color: var(--zi-text-base); font-family: system-ui, sans-serif; overflow: hidden; }
        .sidebar { width: var(--sidebar-w); background: var(--zi-bg-glass-heavy); border-right: 1px solid var(--zi-border-glass); backdrop-filter: blur(20px); display: flex; flex-direction: column; }
        .main { flex: 1; overflow-y: auto; padding: var(--zi-space-8); scroll-behavior: smooth; }
        
        .nav-header { padding: var(--zi-space-6); border-bottom: 1px solid var(--zi-border-subtle); }
        .theme-switcher { padding: var(--zi-space-4); border-bottom: 1px solid var(--zi-border-subtle); display: flex; gap: var(--zi-space-2); flex-wrap: wrap; }
        .theme-pill { padding: 4px 10px; border-radius: 99px; font-size: 0.65rem; font-weight: 700; cursor: pointer; border: 1px solid var(--zi-border-subtle); text-transform: uppercase; transition: 0.2s; }
        .theme-pill:hover { background: var(--zi-accent-primary); color: white; }
        .theme-pill.active { background: var(--zi-accent-primary); color: white; border-color: var(--zi-accent-primary); }
        
        .search { width: 100%; padding: var(--zi-space-3); background: var(--zi-bg-surface); border: 1px solid var(--zi-neutral-700); border-radius: var(--zi-radius-md); color: white; margin-top: var(--zi-space-4); }
        
        .nav-list { flex: 1; overflow-y: auto; padding: var(--zi-space-2); }
        .nav-group { margin-bottom: var(--zi-space-6); }
        .nav-group-title { padding: var(--zi-space-2) var(--zi-space-4); font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--zi-text-muted); letter-spacing: 0.1em; }
        .nav-item { padding: var(--zi-space-2) var(--zi-space-4); cursor: pointer; border-radius: var(--zi-radius-sm); font-size: 0.9rem; transition: 0.2s; color: var(--zi-text-muted); }
        .nav-item:hover { background: var(--zi-bg-elevated); color: var(--zi-accent-primary); }
        .nav-item.active { background: var(--zi-accent-primary); color: white; }
        
        .card { background: var(--zi-bg-surface); border: 1px solid var(--zi-border-subtle); border-radius: var(--zi-radius-lg); padding: var(--zi-space-8); margin-bottom: var(--zi-space-8); box-shadow: var(--zi-shadow-md); animation: slideUp 0.4s ease-out; }
        @keyframes slideUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
        
        .code { background: #000; padding: var(--zi-space-4); border-radius: var(--zi-radius-md); font-family: ui-monospace, monospace; color: #a5b4fc; overflow-x: auto; font-size: 0.9rem; margin: var(--zi-space-4) 0; border: 1px solid var(--zi-neutral-800); }
        .preview { background: var(--zi-bg-base); border: 1px dashed var(--zi-neutral-600); border-radius: var(--zi-radius-md); padding: var(--zi-space-8); margin-top: var(--zi-space-4); min-height: 120px; display: grid; place-items: center; position: relative; overflow: hidden; }
        .preview::before { content: 'PREVIEW'; position: absolute; top: 8px; right: 8px; font-size: 0.6rem; font-weight: 900; opacity: 0.2; letter-spacing: 0.2em; }
    </style>
</head>
<body class="zi-antigravity-root">
    <aside class="sidebar">
        <div class="nav-header">
            <div class="zi-text-2xl zi-font-black text-accent">ZINTENT <span class="zi-text-xs zi-text-muted" style="font-weight: 400;">EXPLORER</span></div>
            <input type="text" class="search" placeholder="Filter intents..." id="filter">
        </div>
        <div class="theme-switcher" id="themes">
            <div class="theme-pill active" onclick="setTheme('nordic', this)">Nordic</div>
            <div class="theme-pill" onclick="setTheme('midnight', this)">Midnight</div>
            <div class="theme-pill" onclick="setTheme('vibrant', this)">Vibrant</div>
            <div class="theme-pill" onclick="setTheme('high-contrast', this)">Contrast</div>
            <div class="theme-pill" onclick="setTheme('forest', this)">Forest</div>
        </div>
        <div class="nav-list" id="navList"></div>
    </aside>
    
    <main class="main" id="main">
        <div id="welcome">
            <h1 class="zi-text-fluid-5xl zi-font-black">Design <span class="text-accent">System</span> Explorer</h1>
            <p class="zi-text-muted zi-text-fluid-lg zi-box-prose">Select any intent or macro from the sidebar to view its specification, generated CSS, and a live rendering.</p>
            <div class="zi-grid-3 zi-gap-6 zi-mt-12">
                <div class="zi-card-flat zi-p-6">
                    <div class="zi-text-xl zi-font-bold text-accent">Intents</div>
                    <p class="zi-text-sm">Semantic styling utilities powered by the core engine.</p>
                </div>
                <div class="zi-card-flat zi-p-6">
                    <div class="zi-text-xl zi-font-bold text-accent">Macros</div>
                    <p class="zi-text-sm">Zero-runtime components with props and named slots.</p>
                </div>
                <div class="zi-card-flat zi-p-6">
                    <div class="zi-text-xl zi-font-bold text-accent">Themes</div>
                    <p class="zi-text-sm">High-fidelity tokens including glassmorphism and depth.</p>
                </div>
            </div>
        </div>
        <div id="inspector"></div>
    </main>

    <script>
        const intentDataJSON = ` + "`" + string(intentData) + "`" + `;
        const macroDataJSON = ` + "`" + string(macroData) + "`" + `;
        const registry = JSON.parse(intentDataJSON);
        const components = JSON.parse(macroDataJSON);
        
        const navList = document.getElementById('navList');
        const inspector = document.getElementById('inspector');
        const welcome = document.getElementById('welcome');
        const filterInput = document.getElementById('filter');

        function renderNav(query = '') {
            navList.innerHTML = '';
            
            // Render Macros
            const mGroup = createGroup('Components (Macros)');
            Object.keys(components).sort().forEach(name => {
                if (name.toLowerCase().includes(query.toLowerCase())) {
                    mGroup.appendChild(createItem(name, () => inspectMacro(name)));
                }
            });
            if (mGroup.children.length > 1) navList.appendChild(mGroup);

            // Render Intents
            const iGroup = createGroup('Styling (Intents)');
            Object.keys(registry).sort().forEach(name => {
                if (name.toLowerCase().includes(query.toLowerCase())) {
                    iGroup.appendChild(createItem(name, () => inspectIntent(name)));
                }
            });
            if (iGroup.children.length > 1) navList.appendChild(iGroup);
        }

        function createGroup(title) {
            const div = document.createElement('div');
            div.className = 'nav-group';
            div.innerHTML = ` + "`" + `<div class="nav-group-title">${title}</div>` + "`" + `;
            return div;
        }

        function createItem(text, onClick) {
            const div = document.createElement('div');
            div.className = 'nav-item';
            div.innerText = text;
            div.onclick = (e) => {
                document.querySelectorAll('.nav-item').forEach(i => i.classList.remove('active'));
                div.classList.add('active');
                onClick();
            };
            return div;
        }

        function inspectMacro(name) {
            welcome.style.display = 'none';
            const template = components[name];
            
            // Expansion demo
            let demo = template.replace(/{{content}}/g, 'Macro Content Block')
                               .replace(/{{slot_header}}/g, 'Header Slot Content')
                               .replace(/{{slot_footer}}/g, 'Footer Slot Content')
                               .replace(/{{attributes}}/g, '')
                               .replace(/{{variant}}/g, 'primary')
                               .replace(/{{size}}/g, 'md');

            inspector.innerHTML = ` + "`" + `
                <div class="card">
                    <h2 class="zi-text-3xl zi-font-black text-accent">${name}</h2>
                    <p class="zi-text-muted">Zero-Runtime Component Definition</p>
                    <div class="code">${template.replace(/</g, '&lt;').replace(/>/g, '&gt;')}</div>
                    
                    <h3 class="zi-text-xl zi-font-bold zi-mt-8">Live Rendering</h3>
                    <div class="preview">
                        ${demo}
                    </div>
                </div>
            ` + "`" + `;
        }

        function inspectIntent(name) {
            welcome.style.display = 'none';
            const intent = registry[name];
            inspector.innerHTML = ` + "`" + `
                <div class="card">
                    <h2 class="zi-text-3xl zi-font-black text-accent">${name}</h2>
                    <p class="zi-text-muted">Semantic Intent Utility</p>
                    <div class="code">
.intent-${name} {
${intent.styles}
${intent.hover ? ` + "`" + `  &:hover { ${intent.hover} }\n` + "`" + ` : ''}
${intent.focus ? ` + "`" + `  &:focus { ${intent.focus} }\n` + "`" + ` : ''}
}</div>
                    
                    <h3 class="zi-text-xl zi-font-bold zi-mt-8">Visual Preview</h3>
                    <div class="preview">
                        <div class="intent-${name} zi-p-8 zi-bg-glass zi-radius-lg zi-border-glass" style="min-width: 200px; text-align: center;">
                            The ${name} Intent
                        </div>
                    </div>
                </div>
            ` + "`" + `;
        }

        function setTheme(theme, el) {
            document.documentElement.setAttribute('data-theme', theme);
            document.querySelectorAll('.theme-pill').forEach(p => p.classList.remove('active'));
            el.classList.add('active');
        }

        filterInput.oninput = (e) => renderNav(e.target.value);
        renderNav();
    </script>
</body>
</html>`

	_ = ioutil.WriteFile(outputPath, []byte(html), 0644)
	fmt.Printf("  [DONE] Explorer generated at %s\n", outputPath)
}
