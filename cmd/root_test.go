package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestGlobalFlags_JSON(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantJSON bool
	}{
		{
			name:     "no flags",
			args:     []string{},
			wantJSON: false,
		},
		{
			name:     "short json flag",
			args:     []string{"-j"},
			wantJSON: true,
		},
		{
			name:     "long json flag",
			args:     []string{"--json"},
			wantJSON: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetFlags()
			cmd := GetRootCmd()
			cmd.SetArgs(tt.args)
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_ = cmd.Execute()

			flags := GetGlobalFlags()
			if flags.JSON != tt.wantJSON {
				t.Errorf("JSON flag = %v, want %v", flags.JSON, tt.wantJSON)
			}
		})
	}
}

func TestGlobalFlags_Pretty(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantPretty bool
	}{
		{
			name:       "no flags",
			args:       []string{},
			wantPretty: false,
		},
		{
			name:       "short pretty flag",
			args:       []string{"-p"},
			wantPretty: true,
		},
		{
			name:       "long pretty flag",
			args:       []string{"--pretty"},
			wantPretty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetFlags()
			cmd := GetRootCmd()
			cmd.SetArgs(tt.args)
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_ = cmd.Execute()

			flags := GetGlobalFlags()
			if flags.Pretty != tt.wantPretty {
				t.Errorf("Pretty flag = %v, want %v", flags.Pretty, tt.wantPretty)
			}
		})
	}
}

func TestGlobalFlags_Combined(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantJSON   bool
		wantPretty bool
	}{
		{
			name:       "both short flags",
			args:       []string{"-j", "-p"},
			wantJSON:   true,
			wantPretty: true,
		},
		{
			name:       "both long flags",
			args:       []string{"--json", "--pretty"},
			wantJSON:   true,
			wantPretty: true,
		},
		{
			name:       "mixed flags",
			args:       []string{"-j", "--pretty"},
			wantJSON:   true,
			wantPretty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetFlags()
			cmd := GetRootCmd()
			cmd.SetArgs(tt.args)
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_ = cmd.Execute()

			flags := GetGlobalFlags()
			if flags.JSON != tt.wantJSON {
				t.Errorf("JSON flag = %v, want %v", flags.JSON, tt.wantJSON)
			}
			if flags.Pretty != tt.wantPretty {
				t.Errorf("Pretty flag = %v, want %v", flags.Pretty, tt.wantPretty)
			}
		})
	}
}

func TestVersionFlag(t *testing.T) {
	ResetFlags()
	SetVersionInfo("1.0.0", "abc123", "2025-01-01")

	cmd := GetRootCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--version"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	// Note: printVersion writes to os.Stdout, not cmd's output
	// The test verifies the flag is parsed correctly
}

func TestHelpOutput(t *testing.T) {
	ResetFlags()
	cmd := GetRootCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--help"})

	_ = cmd.Execute()

	output := buf.String()

	// Check that help contains expected content
	expectedStrings := []string{
		"gobpftool",
		"prog",
		"map",
		"--json",
		"--pretty",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Help output missing expected string: %s", expected)
		}
	}
}

func TestInvalidSubcommand(t *testing.T) {
	ResetFlags()
	cmd := GetRootCmd()
	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(errBuf)
	cmd.SetArgs([]string{"invalidcmd"})

	err := cmd.Execute()

	// Cobra should return an error for unknown commands
	if err == nil {
		t.Error("Expected error for invalid subcommand, got nil")
	}
}

func TestRootCommandNoArgs(t *testing.T) {
	ResetFlags()
	cmd := GetRootCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	// Should show help when no args provided
	output := buf.String()
	if !strings.Contains(output, "gobpftool") {
		t.Error("Expected help output when no args provided")
	}
}
