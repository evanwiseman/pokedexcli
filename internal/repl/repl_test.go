package repl

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/evanwiseman/pokedexcli/internal/pokeapi"
)

type funcStep struct {
	fn             func(ctx *CommandContext, parameters []string) error
	expectContains string
	expectError    bool
}

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
		actual := CleanInput(c.input)

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
		if err := CommandExit(nil, nil); err != nil {
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
		if err := CommandHelp(nil, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
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

func TestCommandMapMapb(t *testing.T) {
	cases := []struct {
		steps []funcStep
	}{
		{
			steps: []funcStep{
				{fn: CommandMap, expectContains: "canalave-city-area", expectError: false},
				{fn: CommandMap, expectContains: "mt-coronet-1f-route-216", expectError: false},
			},
		},
		{
			steps: []funcStep{
				{fn: CommandMap, expectContains: "canalave-city-area", expectError: false},
				{fn: CommandMap, expectContains: "mt-coronet-1f-route-216", expectError: false},
				{fn: CommandMapb, expectContains: "canalave-city-area", expectError: false},
			},
		},
		{
			steps: []funcStep{
				{fn: CommandMapb, expectContains: "", expectError: true},
			},
		},
		{
			steps: []funcStep{
				{fn: CommandMap, expectContains: "canalave-city-area", expectError: false},
				{fn: CommandMapb, expectContains: "", expectError: true},
			},
		},
	}
	for _, c := range cases {
		ctx := CommandContext{
			Client: pokeapi.NewClient(),
			LocationConfig: &pokeapi.Config{
				Next:     strPtr(pokeapi.LocationAreaURL),
				Previous: nil,
			},
		}
		for _, step := range c.steps {
			r, w, _ := os.Pipe()
			old := os.Stdout
			os.Stdout = w

			err := step.fn(&ctx, nil)

			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)
			os.Stdout = old

			if step.expectError && err == nil {
				t.Errorf("expected error but got nil")
			} else if !step.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !strings.Contains(buf.String(), step.expectContains) {
				t.Errorf("expected output to contain %s, got %s", step.expectContains, buf.String())
			}
		}
	}
}

func TestCommandExplore(t *testing.T) {
	cases := []struct {
		steps []funcStep
	}{
		{
			steps: []funcStep{
				{
					fn: func(ctx *CommandContext, _ []string) error {
						return CommandExplore(ctx, []string{"canalave-city-area"})
					},
					expectContains: "staryu", // expect one Pokémon known in the area
					expectError:    false,
				},
			},
		},
		{
			steps: []funcStep{
				{
					fn: func(ctx *CommandContext, _ []string) error {
						return CommandExplore(ctx, []string{}) // no parameter
					},
					expectContains: "",
					expectError:    true,
				},
			},
		},
		{
			steps: []funcStep{
				{
					fn: func(ctx *CommandContext, _ []string) error {
						return CommandExplore(ctx, []string{"invalid-area"})
					},
					expectContains: "",
					expectError:    true,
				},
			},
		},
	}

	for _, c := range cases {
		ctx := CommandContext{
			Client: pokeapi.NewClient(),
		}

		for _, step := range c.steps {
			r, w, _ := os.Pipe()
			old := os.Stdout
			os.Stdout = w

			err := step.fn(&ctx, nil)

			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)
			os.Stdout = old

			if step.expectError && err == nil {
				t.Errorf("expected error but got nil")
			} else if !step.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if step.expectContains != "" && !strings.Contains(buf.String(), step.expectContains) {
				t.Errorf("expected output to contain %q, got %q", step.expectContains, buf.String())
			}
		}
	}
}

func TestCommandCatch(t *testing.T) {
	ctx := CommandContext{
		Client: pokeapi.NewClient(),
	}

	cases := []struct {
		parameters     []string
		expectContains string
		expectError    bool
	}{
		{
			parameters:     []string{"pikachu"},
			expectContains: "pikachu", // output should include the Pokémon name
			expectError:    false,
		},
		{
			parameters:     []string{},
			expectContains: "",
			expectError:    true, // no Pokémon name provided
		},
		{
			parameters:     []string{"invalidpokemon"},
			expectContains: "",
			expectError:    true, // Pokémon does not exist
		},
	}

	for _, c := range cases {
		r, w, _ := os.Pipe()
		old := os.Stdout
		os.Stdout = w

		err := CommandCatch(&ctx, c.parameters)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old

		if c.expectError && err == nil {
			t.Errorf("expected error but got nil")
		} else if !c.expectError && err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if c.expectContains != "" && !strings.Contains(buf.String(), c.expectContains) {
			t.Errorf("expected output to contain %q, got %q", c.expectContains, buf.String())
		}
	}
}
