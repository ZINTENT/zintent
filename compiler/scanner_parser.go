package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// extractClassesHTMLParser walks the DOM and collects `class` attribute values (HTML only).
func extractClassesHTMLParser(content string, usedClasses map[string]bool) error {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return err
	}
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if strings.EqualFold(a.Key, "class") {
					addClassesFromValue(a.Val, usedClasses)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return nil
}

// extractVueTemplate attempts to find the content inside the <template> tag.
func extractVueTemplate(content string) string {
	startTag := "<template>"
	endTag := "</template>"
	start := strings.Index(content, startTag)
	if start == -1 {
		return ""
	}
	end := strings.LastIndex(content, endTag)
	if end == -1 || end <= start {
		return ""
	}
	return content[start+len(startTag) : end]
}

// extractClassesFromPath chooses regex vs HTML DOM parser based on file extension and --scanner mode.
func extractClassesFromPath(filePath string, content string, usedClasses map[string]bool, scanner string, verbose bool) {
	ext := strings.ToLower(filepath.Ext(filePath))
	isHTML := ext == ".html" || ext == ".htm" || ext == ".php"
	isVue := ext == ".vue"

	useHTMLParser := strings.EqualFold(scanner, "parser") && (isHTML || isVue)

	if useHTMLParser {
		parseContent := content
		if isVue {
			if vueTmpl := extractVueTemplate(content); vueTmpl != "" {
				parseContent = vueTmpl
			}
		}

		if err := extractClassesHTMLParser(parseContent, usedClasses); err != nil {
			if verbose {
				fmt.Printf("  [SCAN] HTML parser fallback for %s: %v\n", filePath, err)
			}
			extractClassesFromContent(content, usedClasses)
			return
		}

		// For PHP/Vue, we also run regex to catch dynamic classes in scripts/expressions
		if ext == ".php" || isVue {
			extractClassesFromContent(content, usedClasses)
		}
		return
	}
	extractClassesFromContent(content, usedClasses)
}
// AuditUX performs build-time UX and Accessibility analysis on HTML content.
func AuditUX(content string, filePath string) []string {
	findings := []string{}
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return findings
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// 1. Touch Target Compliance
			isSmall := false
			isBtn := n.Data == "button" || n.Data == "a" || strings.HasPrefix(n.Data, "z-btn")
			
			for _, a := range n.Attr {
				if a.Key == "size" && (a.Val == "sm" || a.Val == "xs") {
					isSmall = true
				}
				if a.Key == "class" && (strings.Contains(a.Val, "zi-text-xs") || strings.Contains(a.Val, "btn-sm")) {
					isSmall = true
				}
			}
			
			if isBtn && isSmall {
				findings = append(findings, fmt.Sprintf("Small interactive element <%s> in %s may violate touch-target minimums (44px).", n.Data, filepath.Base(filePath)))
			}

			// 2. Accessible Names for Icon Buttons
			if n.Data == "button" || strings.HasPrefix(n.Data, "z-btn") {
				hasAriaLabel := false
				for _, a := range n.Attr {
					if a.Key == "aria-label" && a.Val != "" {
						hasAriaLabel = true
					}
				}
				
				// Check if button is effectively empty (only icon or whitespace)
				isEmpty := true
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode && strings.TrimSpace(c.Data) != "" {
						isEmpty = false
						break
					}
					if c.Type == html.ElementNode && c.Data != "i" && c.Data != "svg" {
						isEmpty = false
						break
					}
				}
				
				if isEmpty && !hasAriaLabel {
					findings = append(findings, fmt.Sprintf("Icon-only button <%s> in %s needs an aria-label for screen readers.", n.Data, filepath.Base(filePath)))
				}
			}

			// 3. Nesting Violations
			if n.Data == "a" {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.ElementNode && (c.Data == "a" || c.Data == "button") {
						findings = append(findings, fmt.Sprintf("Invalid nesting: <%s> inside <a> in %s. This breaks interaction model.", c.Data, filepath.Base(filePath)))
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return findings
}
