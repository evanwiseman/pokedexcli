package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  HELLO  WORLD  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO WORLD",
			expected: []string{"hello", "world"},
		},
		{
			input:    "thisis/text and a field",
			expected: []string{"thisis/text", "and", "a", "field"},
		},
		{
			input:    "hElLo WoRlD =,./[]",
			expected: []string{"hello", "world", "=,./[]"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("actual length=%v, expected length=%v", len(actual), len(c.expected))
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("word '%v' != expected '%v'", word, expectedWord)
				t.Fail()
			}
		}
	}
}

func TestCommandExit(t *testing.T) {
	if os.Getenv("TEST_EXIT") == "1" {
		if err := commandExit(); err != nil {
			t.Fatalf("unexpected error: %v", err)
			t.Fail()
		}
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCommandExit")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1")

	output, err := cmd.CombinedOutput()

	// check exit code
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 0 {
			t.Fatalf("expected exit code 0, got %d", exitErr.ExitCode())
		}
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// check printed output
	if !strings.Contains(string(output), "Closing the Pokedex... Goodbye!") {
		t.Errorf("expected exit message, got: %q", string(output))
	}
}

func TestCommandHelp(t *testing.T) {
	if os.Getenv("TEST_HELP") == "1" {
		if err := commandHelp(); err != nil {
			t.Fatalf("unexpected error: %v", err)
			t.Fail()
		}
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCommandHelp")
	cmd.Env = append(os.Environ(), "TEST_HELP=1")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(output), "Usage") {
		t.Errorf("expected 'Usage' in output, got: %q", string(output))
	}
}
