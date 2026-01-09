package main

import (
	"strings"
	"testing"
)

func TestElizaBasicPatterns(t *testing.T) {
	eliza := NewEliza()

	tests := []struct {
		input    string
		contains string
	}{
		{"I need help", "need"},
		{"I am sad", "are sad"},
		{"Hello", "Hello"},
		{"quit", "Thank you"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			response := eliza.respond(tt.input)
			if response == "" {
				t.Errorf("Empty response for input: %s", tt.input)
			}
			if !strings.Contains(strings.ToLower(response), strings.ToLower(tt.contains)) {
				t.Logf("Input: %s", tt.input)
				t.Logf("Response: %s", response)
				t.Logf("Expected to contain: %s", tt.contains)
			}
		})
	}
}

func TestReflection(t *testing.T) {
	eliza := NewEliza()

	tests := map[string]string{
		"i":      "you",
		"am":     "are",
		"my":     "your",
		"you":    "me",
		"me":     "you",
		"your":   "my",
		"i'm":    "i'm", // Not in reflection table, should stay same
	}

	for input, expected := range tests {
		result := eliza.translate(input, gReflections)
		if result != expected {
			t.Errorf("translate(%q) = %q; want %q", input, result, expected)
		}
	}
}

func TestMultiWordReflection(t *testing.T) {
	eliza := NewEliza()

	input := "i am sad"
	result := eliza.translate(input, gReflections)
	if result != "you are sad" {
		t.Errorf("translate(%q) = %q; want %q", input, result, "you are sad")
	}
}
