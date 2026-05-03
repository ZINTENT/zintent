package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractClassesHTMLParser_NestedClasses(t *testing.T) {
	root := findProjectRoot(t)
	fixture := filepath.Join(root, "tests", "fixtures", "scanner-nested.html")
	b, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	used := make(map[string]bool)
	if err := extractClassesHTMLParser(string(b), used); err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := []string{"outer", "inner", "deep", "zi-fixture-a", "zi-fixture-b", "@sm:compact"}
	for _, w := range want {
		if !used[w] {
			t.Errorf("missing class token %q", w)
		}
	}
}

func findProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("go.mod not found from test cwd")
	return ""
}

func TestExtractClassesFromPath_RegexJSX(t *testing.T) {
	used := make(map[string]bool)
	src := `export function X() { return <div className="zi-jsx-one zi-jsx-two">hi</div>; }`
	extractClassesFromPath("app.tsx", src, used, "regex", false)
	if !used["zi-jsx-one"] || !used["zi-jsx-two"] {
		t.Errorf("expected JSX className tokens, got keys: %v", keysSample(used, 10))
	}
}

func keysSample(m map[string]bool, n int) []string {
	var out []string
	for k := range m {
		out = append(out, k)
		if len(out) >= n {
			break
		}
	}
	return out
}

func TestExtractClassesFromPath_ParserSkipsVueForDOM(t *testing.T) {
	// .vue uses regex path; parser mode does not treat whole SFC as HTML document.
	used := make(map[string]bool)
	src := `<template><div class="zi-vue-a"></div></template>`
	extractClassesFromPath("c.vue", src, used, "parser", false)
	if !used["zi-vue-a"] {
		t.Error("expected regex extraction for .vue in parser mode")
	}
}

func TestShouldIncludeDependency_CoreOmitsCrossBrowser(t *testing.T) {
	if !shouldIncludeDependency("core/tokens.css", nil, "core") {
		t.Error("core should include tokens")
	}
	if shouldIncludeDependency("core/cross-browser.css", nil, "core") {
		t.Error("core should omit cross-browser.css")
	}
	if !shouldIncludeDependency("core/cross-browser.css", nil, "full") {
		t.Error("full should include cross-browser.css")
	}
}

func TestShouldIncludeDependency_MinimalAntigravityIntentPrefix(t *testing.T) {
	used := map[string]bool{"intent-stack-md": true}
	if !shouldIncludeDependency("core/antigravity-layouts.css", used, "minimal") {
		t.Error("minimal should include antigravity when any intent-* layout class is used")
	}
}
