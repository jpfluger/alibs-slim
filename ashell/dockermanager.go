package ashell

import (
	"fmt"
	"strings"
)

type DockerManager struct {
	WorkingDir    string
	DCFile        string // Path to the Docker Compose file
	Project       string // Project used for Docker services
	EnvFile       string // Path to environment file (optional)
	DisplayStdErr bool   // Option to display standard error
	HideStdOut    bool   // Option to hide standard output
	DoPrint       bool
}

// StartDocker starts the Docker Compose services.
func (dm *DockerManager) StartDocker() error {
	fmt.Println("Starting Docker services...")

	params := []string{"-f", dm.DCFile, "-p", dm.Project}
	dm.EnvFile = strings.TrimSpace(dm.EnvFile)
	if dm.EnvFile != "" {
		params = append(params, "--env-file", dm.EnvFile)
	}
	params = append(params, "up", "-d")

	cmd := ShellCommand{
		Type:       SHELLTYPE_DCOMPOSE,
		HideStdOut: dm.HideStdOut,
		WorkingDir: dm.WorkingDir,
		Parameters: params,
		//Parameters: []string{"-f", dm.DCFile, "-p", dm.Project, "up", "-d"},
		//Command:    fmt.Sprintf("-f %s -p %s up -d", dm.DCFile, dm.Project),
	}
	if _, stdErr, err := cmd.RunCommand(); err != nil {
		return fmt.Errorf("failed to start Docker Compose services: %v", err)
	} else if strings.TrimSpace(stdErr) != "" {
		if dm.DoPrint {
			dm.printStdErr(stdErr)
		}
	}
	return nil
}

// StopDocker stops and removes the Docker Compose services.
func (dm *DockerManager) StopDocker() error {
	fmt.Println("Stopping Docker services...")
	cmd := ShellCommand{
		Type:       SHELLTYPE_DCOMPOSE,
		HideStdOut: dm.HideStdOut,
		WorkingDir: dm.WorkingDir,
		Parameters: []string{"-f", dm.DCFile, "-p", dm.Project, "down", "-v"},
		//Command:    fmt.Sprintf("-f %s -p %s down -v", dm.DCFile, dm.Project),
	}
	if _, stdErr, err := cmd.RunCommand(); err != nil {
		return fmt.Errorf("failed to stop Docker Compose services: %v", err)
	} else if strings.TrimSpace(stdErr) != "" {
		if dm.DoPrint {
			dm.printStdErr(stdErr)
		}
	}
	return nil
}

// IsDockerRunning checks if Docker services matching the prefix are running.
func (dm *DockerManager) IsDockerRunning() (bool, error) {
	cmd := ShellCommand{
		Type:       SHELLTYPE_DOCKER,
		HideStdOut: dm.HideStdOut,
		WorkingDir: dm.WorkingDir,
		//Command:    fmt.Sprintf("ps -a -q --filter name=%s*", dm.Project),
		Command: "ps",
		//Parameters: []string{"ps", "-a", "-q", "--filter", fmt.Sprintf("name=%s*", dm.Project)},
		Parameters: []string{"-a", "-q", "--filter", fmt.Sprintf("name=%s*", dm.Project)},
	}
	stdOut, stdErr, err := cmd.RunCommand()
	if err != nil {
		return true, fmt.Errorf("failed to detect Docker services: %v", err)
	} else if strings.TrimSpace(stdErr) != "" {
		if dm.DoPrint {
			dm.printStdErr(stdErr)
		}
		return true, fmt.Errorf("docker detection returned an error: %v", err)
	} else if strings.TrimSpace(stdOut) != "" {
		return true, nil
	}
	return false, nil
}

// printStdErr optionally prints standard error output if enabled.
func (dm *DockerManager) printStdErr(stdErr string) {
	if dm.DisplayStdErr {
		fmt.Println("....... standard-error (begin) .......")
		fmt.Println(strings.TrimSpace(stdErr))
		fmt.Println("....... standard-error (end) .......")
	}
}
