package ashell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ShellCommand represents a shell command with its parameters.
type ShellCommand struct {
	Type       ShellType `json:"type,omitempty"`
	WorkingDir string    `json:"workingDir,omitempty"`
	Command    string    `json:"command,omitempty"`
	Parameters []string  `json:"parameters,omitempty"`
	HideStdOut bool      `json:"hideStdOut,omitempty"`
	IsFakeRun  bool      `json:"isFakeRun,omitempty"`

	// Option to ignore running the command.
	Ignore bool `json:"ignore,omitempty"`

	// Optionally set by calling program.
	// While WorkingDir can be hardcoded into json,
	// cdRootDir can be set by a higher-level app manager.
	// If WorkingDir is not empty and absolute, then
	// cdRootDir is ignored, otherwise if cdRootDir is
	// not empty then cdRootDir is used in lieu of WorkingDir.
	cdRootDir string
}

// GetType returns the type of the ShellCommand.
func (sc *ShellCommand) GetType() ShellType {
	return sc.Type
}

// GetCommand returns the command string of the ShellCommand.
func (sc *ShellCommand) GetCommand() string {
	return sc.Command
}

// GetParameters returns the parameters of the ShellCommand.
func (sc *ShellCommand) GetParameters() []string {
	return sc.Parameters
}

// SetCDRootDir sets the cdRootDir for the ShellCommand.
func (sc *ShellCommand) SetCDRootDir(dir string) {
	sc.cdRootDir = dir
}

// GetCDRootDir returns the cdRootDir of the ShellCommand.
func (sc *ShellCommand) GetCDRootDir() string {
	return sc.cdRootDir
}

// RunCommand executes the shell command with its parameters.
func (sc *ShellCommand) RunCommand() (string, string, error) {
	return sc.RunCommandWithOptions(sc.HideStdOut, sc.IsFakeRun)
}

// RunCommandWithOptions executes the shell command with its parameters, considering the isFakeRun option.
func (sc *ShellCommand) RunCommandWithOptions(hideStdOut bool, isFakeRun bool) (string, string, error) {
	// Ignore command skips running any commands.
	if sc.Ignore {
		return "", "", nil
	}
	// Validate that Command and Parameters have values.
	if strings.TrimSpace(sc.Command) == "" {
		if sc.Parameters == nil || len(sc.Parameters) == 0 {
			return "", "", fmt.Errorf("command cannot be empty")
		}
	}
	for _, param := range sc.Parameters {
		if strings.TrimSpace(param) == "" {
			return "", "", fmt.Errorf("parameters cannot contain empty values")
		}
	}
	// Determine the working directory
	workingDir := sc.WorkingDir
	if workingDir == "" && sc.cdRootDir != "" {
		workingDir = sc.cdRootDir
	} else if workingDir != "" && !filepath.IsAbs(workingDir) && sc.cdRootDir != "" {
		workingDir = filepath.Join(sc.cdRootDir, workingDir)
	}

	// Clean the working directory path
	workingDir = filepath.Clean(workingDir)

	// Validate that WorkingDir exists, if applicable
	if workingDir != "" {
		if _, err := os.Stat(workingDir); os.IsNotExist(err) {
			return "", "", fmt.Errorf("working directory does not exist")
		}
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	// Validate the ShellType and check for the correct operating system.
	switch sc.Type {
	case SHELLTYPE_BASH, SHELLTYPE_SH:
		if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
			return "", "", fmt.Errorf("unsupported operating system for shell type: %s", sc.Type)
		}
	case SHELLTYPE_WIN:
		if runtime.GOOS != "windows" {
			return "", "", fmt.Errorf("unsupported operating system for shell type: %s", sc.Type)
		}
	case SHELLTYPE_PWS:
		if runtime.GOOS != "windows" {
			return "", "", fmt.Errorf("unsupported operating system for shell type: %s", sc.Type)
		}
	case SHELLTYPE_OSX:
		if runtime.GOOS != "darwin" {
			return "", "", fmt.Errorf("unsupported operating system for shell type: %s", sc.Type)
		}
	case SHELLTYPE_PYTHON, SHELLTYPE_PYTHON3, SHELLTYPE_DOCKER, SHELLTYPE_DCOMPOSE:
		// No specific OS restrictions for Python, Docker, and Docker Compose
	default:
		return "", "", fmt.Errorf("unsupported shell type: %s", sc.Type)
	}

	// Debug
	// Print the absolute path of the working directory for debugging
	//absWorkingDir, err := filepath.Abs(workingDir)
	//if err != nil {
	//	return "", "", fmt.Errorf("failed to get absolute path of working directory: %v", err)
	//}
	//fmt.Printf("Absolute working directory: %s\n", absWorkingDir)

	// If isFakeRun is true, simulate the command execution.
	if isFakeRun {
		simulatedOutput := fmt.Sprintf("Simulated command: %s %s", sc.Command, strings.Join(sc.Parameters, " "))
		stdoutBuf.WriteString(simulatedOutput)
		return stdoutBuf.String(), "", nil
	}

	fnCmdArgs := func(name string, arg ...string) *exec.Cmd {
		arr := []string{}
		if arg != nil && len(arg) > 0 {
			for _, ar := range arg {
				if strings.TrimSpace(ar) != "" {
					arr = append(arr, ar)
				}
			}
		}
		return exec.Command(name, arr...)
	}
	// Set up the command based on the ShellType.
	var cmd *exec.Cmd
	switch sc.Type {
	case SHELLTYPE_BASH:
		cmd = fnCmdArgs("bash", "-c", sc.Command)
	case SHELLTYPE_SH:
		cmd = fnCmdArgs("sh", "-c", sc.Command)
	case SHELLTYPE_WIN:
		cmd = fnCmdArgs("cmd", "/C", sc.Command)
	case SHELLTYPE_PWS:
		cmd = fnCmdArgs("powershell", "-Command", sc.Command)
	case SHELLTYPE_OSX:
		cmd = fnCmdArgs("sh", "-c", sc.Command)
	case SHELLTYPE_PYTHON:
		cmd = fnCmdArgs("python", sc.Command)
	case SHELLTYPE_PYTHON3:
		cmd = fnCmdArgs("python3", sc.Command)
	case SHELLTYPE_DOCKER:
		cmd = fnCmdArgs("docker", sc.Command)
	case SHELLTYPE_DCOMPOSE:
		cmd = fnCmdArgs("docker", "compose", sc.Command)
	}

	// Add parameters to the command
	if len(sc.Parameters) > 0 {
		cmd.Args = append(cmd.Args, sc.Parameters...)
	}

	// Change the working directory if WorkingDir is not empty.
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	stdout := io.Writer(&stdoutBuf)
	stderr := io.Writer(&stderrBuf)
	if !hideStdOut {
		stdout = io.MultiWriter(os.Stdout, stdout)
		stderr = io.MultiWriter(os.Stderr, stderr)
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("cmd.Run failed with %v", err)
		return "", "", err
	}

	outStr, errStr := stdoutBuf.String(), stderrBuf.String()
	return outStr, errStr, nil
}

// ShellCommands is a slice of pointers to ShellCommand structs.
type ShellCommands []*ShellCommand

// Validate runs the RunCommandWithOptions(true, true) for each ShellCommand and returns the index of the command that failed.
func (scs *ShellCommands) Validate() (int, error) {
	for i, sc := range *scs {
		_, _, err := sc.RunCommandWithOptions(true, true)
		if err != nil {
			return i, err
		}
	}
	return -1, nil
}

// RemoveAtIndex removes the ShellCommand at the specified index.
func (scs *ShellCommands) RemoveAtIndex(index int) error {
	if index < 0 || index >= len(*scs) {
		return fmt.Errorf("index out of range")
	}
	*scs = append((*scs)[:index], (*scs)[index+1:]...)
	return nil
}

// ReplaceAtIndex replaces the ShellCommand at the specified index with a new ShellCommand.
func (scs *ShellCommands) ReplaceAtIndex(index int, newCommand *ShellCommand) error {
	if index < 0 || index >= len(*scs) {
		return fmt.Errorf("index out of range")
	}
	(*scs)[index] = newCommand
	return nil
}
