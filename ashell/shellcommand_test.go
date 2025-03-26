package ashell

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellCommand_GetType(t *testing.T) {
	sc := &ShellCommand{Type: SHELLTYPE_BASH}
	assert.Equal(t, SHELLTYPE_BASH, sc.GetType())
}

func TestShellCommand_GetCommand(t *testing.T) {
	sc := &ShellCommand{Command: "echo"}
	assert.Equal(t, "echo", sc.GetCommand())
}

func TestShellCommand_GetParameters(t *testing.T) {
	sc := &ShellCommand{Parameters: []string{"Hello, World!"}}
	assert.Equal(t, []string{"Hello, World!"}, sc.GetParameters())
}

func TestShellCommand_RunCommand_EmptyCommand(t *testing.T) {
	sc := &ShellCommand{Command: ""}
	_, _, err := sc.RunCommand()
	assert.Error(t, err)
	assert.Equal(t, "command cannot be empty", err.Error())
}

func TestShellCommand_RunCommand_EmptyParameter(t *testing.T) {
	sc := &ShellCommand{Command: "echo", Parameters: []string{"Hello", ""}}
	_, _, err := sc.RunCommand()
	assert.Error(t, err)
	assert.Equal(t, "parameters cannot contain empty values", err.Error())
}

func TestShellCommand_RunCommand_FakeRun_Bash(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-Unix-like OS")
	}
	sc := &ShellCommand{
		Type:       SHELLTYPE_BASH,
		Command:    "echo",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  true,
	}
	outStr, errStr, err := sc.RunCommand()
	assert.NoError(t, err)
	assert.Equal(t, true, strings.HasPrefix(outStr, "Simulated command:"))
	assert.Equal(t, "", errStr)
}

func TestShellCommand_RunCommand_FakeRun_Sh(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-Unix-like OS")
	}
	sc := &ShellCommand{
		Type:       SHELLTYPE_SH,
		Command:    "echo",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  true,
	}
	outStr, errStr, err := sc.RunCommand()
	assert.NoError(t, err)
	assert.Equal(t, true, strings.HasPrefix(outStr, "Simulated command:"))
	assert.Equal(t, "", errStr)
}

func TestShellCommand_RunCommand_FakeRun_Win(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping test on non-Windows OS")
	}
	sc := &ShellCommand{
		Type:       SHELLTYPE_WIN,
		Command:    "echo",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  true,
	}
	outStr, errStr, err := sc.RunCommand()
	assert.NoError(t, err)
	assert.Equal(t, true, strings.HasPrefix(outStr, "Simulated command:"))
	assert.Equal(t, "", errStr)
}

func TestShellCommand_RunCommand_FakeRun_Pws(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping test on non-Windows OS")
	}
	sc := &ShellCommand{
		Type:       SHELLTYPE_PWS,
		Command:    "Write-Output",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  true,
	}
	outStr, errStr, err := sc.RunCommand()
	assert.NoError(t, err)
	assert.Equal(t, true, strings.HasPrefix(outStr, "Simulated command:"))
	assert.Equal(t, "", errStr)
}

func TestShellCommand_RunCommand_FakeRun_Osx(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-MacOSX OS")
	}
	sc := &ShellCommand{
		Type:       SHELLTYPE_OSX,
		Command:    "echo",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  true,
	}
	outStr, errStr, err := sc.RunCommand()
	assert.NoError(t, err)
	assert.Equal(t, true, strings.HasPrefix(outStr, "Simulated command:"))
	assert.Equal(t, "", errStr)
}

func TestShellCommand_RunCommand_UnsupportedShellType(t *testing.T) {
	sc := &ShellCommand{
		Type:       ShellType("unsupported"),
		Command:    "echo",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  true,
	}
	outStr, errStr, err := sc.RunCommand()
	assert.Error(t, err)
	assert.Equal(t, "", outStr)
	assert.Equal(t, "", errStr)
	assert.Equal(t, "unsupported shell type: unsupported", err.Error())
}

func TestShellCommand_RunCommand_RealRun_Bash(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-Unix-like OS")
	}
	sc := &ShellCommand{
		Type:    SHELLTYPE_BASH,
		Command: `echo "Hello, World!"`,
		//Parameters: []string{`"Hello, World!"`}, // DOES NOT WORK!
		IsFakeRun:  false,
		HideStdOut: true,
	}
	outStr, errStr, err := sc.RunCommand()
	assert.NoError(t, err)
	assert.Contains(t, strings.TrimSpace(outStr), "Hello, World!")
	assert.Equal(t, "", strings.TrimSpace(errStr))
}

func TestShellCommand_RunCommand_RealRun_Sh(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-Unix-like OS")
	}
	sc := &ShellCommand{
		Type:    SHELLTYPE_SH,
		Command: `echo "Hello, World!"`,
		//Parameters: []string{"Hello, World!"}, // DOES NOT WORK
		IsFakeRun:  false,
		HideStdOut: true,
	}
	outStr, errStr, err := sc.RunCommand()
	assert.NoError(t, err)
	assert.Contains(t, strings.TrimSpace(outStr), "Hello, World!")
	assert.Equal(t, "", strings.TrimSpace(errStr))
}

func TestShellCommand_RunCommand_RealRun_Win(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping test on non-Windows OS")
	}
	sc := &ShellCommand{
		Type:       SHELLTYPE_WIN,
		Command:    "echo",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  false,
		HideStdOut: true,
	}
	outStr, errStr, err := sc.RunCommand()
	assert.NoError(t, err)
	assert.Contains(t, outStr, "Hello, World!")
	assert.Equal(t, "", errStr)
}

func TestShellCommand_RunCommand_RealRun_Pws(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping test on non-Windows OS")
	}
	sc := &ShellCommand{
		Type:       SHELLTYPE_PWS,
		Command:    "Write-Output",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  false,
		HideStdOut: true,
	}
	outStr, errStr, err := sc.RunCommand()
	assert.NoError(t, err)
	assert.Contains(t, outStr, "Hello, World!")
	assert.Equal(t, "", errStr)
}

func TestShellCommand_RunCommand_RealRun_Osx(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-MacOSX OS")
	}
	sc := &ShellCommand{
		Type:       SHELLTYPE_OSX,
		Command:    "echo",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  false,
		HideStdOut: true,
	}
	outStr, errStr, err := sc.RunCommand()
	assert.NoError(t, err)
	assert.Contains(t, outStr, "Hello, World!")
	assert.Equal(t, "", errStr)
}

func TestRunCommandWithOptions_Bash(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-Linux/non-macOS system")
	}

	sc := &ShellCommand{
		Type:       SHELLTYPE_BASH,
		Command:    "echo",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  true,
	}

	stdout, stderr, err := sc.RunCommandWithOptions(false, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stdout == "" {
		t.Errorf("Expected stdout, got empty string")
	}
	if stderr != "" {
		t.Errorf("Expected no stderr, got %v", stderr)
	}
}

func TestRunCommandWithOptions_Win(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping test on non-Windows system")
	}

	sc := &ShellCommand{
		Type:       SHELLTYPE_WIN,
		Command:    "echo",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  true,
	}

	stdout, stderr, err := sc.RunCommandWithOptions(false, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stdout == "" {
		t.Errorf("Expected stdout, got empty string")
	}
	if stderr != "" {
		t.Errorf("Expected no stderr, got %v", stderr)
	}
}

func TestRunCommandWithOptions_PWS(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping test on non-Windows system")
	}

	sc := &ShellCommand{
		Type:       SHELLTYPE_PWS,
		Command:    "Write-Output",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  true,
	}

	stdout, stderr, err := sc.RunCommandWithOptions(false, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stdout == "" {
		t.Errorf("Expected stdout, got empty string")
	}
	if stderr != "" {
		t.Errorf("Expected no stderr, got %v", stderr)
	}
}

func TestRunCommandWithOptions_OSX(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-macOS system")
	}

	sc := &ShellCommand{
		Type:       SHELLTYPE_OSX,
		Command:    "echo",
		Parameters: []string{"Hello, World!"},
		IsFakeRun:  true,
	}

	stdout, stderr, err := sc.RunCommandWithOptions(false, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stdout == "" {
		t.Errorf("Expected stdout, got empty string")
	}
	if stderr != "" {
		t.Errorf("Expected no stderr, got %v", stderr)
	}
}

func TestShellCommands_Validate(t *testing.T) {
	tests := []struct {
		name     string
		commands ShellCommands
		wantErr  bool
	}{
		{"Valid commands", ShellCommands{
			&ShellCommand{Type: SHELLTYPE_BASH, Command: "echo", Parameters: []string{"Hello, World!"}},
			&ShellCommand{Type: SHELLTYPE_PYTHON3, Command: "print('Hello, World!')"},
		}, false},
		{"Invalid command", ShellCommands{
			&ShellCommand{Type: SHELLTYPE_BASH, Command: "echo", Parameters: []string{"Hello, World!"}},
			&ShellCommand{Type: ShellType("unsupported"), Command: "echo", Parameters: []string{"Hello, World!"}},
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.commands.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShellCommands_RemoveAtIndex(t *testing.T) {
	commands := ShellCommands{
		&ShellCommand{Type: SHELLTYPE_BASH, Command: "echo", Parameters: []string{"Hello, World!"}},
		&ShellCommand{Type: SHELLTYPE_PYTHON3, Command: "print('Hello, World!')"},
	}

	tests := []struct {
		name    string
		index   int
		wantErr bool
	}{
		{"Valid index", 1, false},
		{"Invalid index", 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := commands.RemoveAtIndex(tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveAtIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShellCommands_ReplaceAtIndex(t *testing.T) {
	commands := ShellCommands{
		&ShellCommand{Type: SHELLTYPE_BASH, Command: "echo", Parameters: []string{"Hello, World!"}},
		&ShellCommand{Type: SHELLTYPE_PYTHON3, Command: "print('Hello, World!')"},
	}

	newCommand := &ShellCommand{Type: SHELLTYPE_DOCKER, Command: "hello-world"}

	tests := []struct {
		name    string
		index   int
		wantErr bool
	}{
		{"Valid index", 1, false},
		{"Invalid index", 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := commands.ReplaceAtIndex(tt.index, newCommand)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceAtIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
