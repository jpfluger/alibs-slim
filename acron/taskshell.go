package acron

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/aerr"
	"github.com/jpfluger/alibs-slim/ashell"
)

type ITaskShell interface {
	ITask
}

const TASKTYPE_SHELL TaskType = "shell"

// TaskShell is a concrete type that implements ITask.
type TaskShell struct {
	Type          TaskType             `json:"type"`
	ShellCommands ashell.ShellCommands `json:"shellCommands,omitempty"`
	cronCC        ICronControlCenter
}

// GetType returns the type of the TaskShell.
func (ts *TaskShell) GetType() TaskType {
	return ts.Type
}

// Validate ensures quality control on this struct.
func (ts *TaskShell) Validate() error {
	if ts.Type.IsEmpty() {
		ts.Type = TASKTYPE_SHELL
	}
	for _, cmd := range ts.ShellCommands {
		if cmd.Ignore {
			continue
		}
		// Let's save the workingDir temporarily. There is a potential error that the
		// SetCDRootDir is required for a relative path for WorkingDir to be validated
		// without error. The check would happen naturally in Run.
		saveWorkingDir := cmd.WorkingDir
		cmd.WorkingDir = ""
		// Run simulated command execution
		_, _, err := cmd.RunCommandWithOptions(true, true)
		if err != nil {
			return fmt.Errorf("failed to validate command: %s, error: %v", cmd.Command, err)
		}
		cmd.WorkingDir = saveWorkingDir
	}
	return nil
}

// Run executes the shell commands.
func (ts *TaskShell) Run(ccc ICronControlCenter) error {
	if ccc == nil {
		return fmt.Errorf("nil cronControlCenter")
	}
	cccOS, ok := ccc.(ICronControlCenterShell)
	if !ok {
		return fmt.Errorf("ccc is not a ICronControlCenterShell")
	}
	for ii, cmd := range ts.ShellCommands {
		if cmd.Ignore {
			cccOS.GetJRun().Logger().Info().Msgf("shell command %d of %d: skipping", ii+1, len(ts.ShellCommands))
			continue
		}
		cmd.SetCDRootDir(cccOS.GetCDRootDir())

		cccOS.GetJRun().Logger().Info().Msgf("shell command %d of %d: begin", ii+1, len(ts.ShellCommands))

		stdOut, stdErr, err := cmd.RunCommandWithOptions(cccOS.GetHideStdOut(), false)
		if err != nil {
			cccOS.GetJRun().Logger().Info().Msgf("err: %v", err)
		}
		cccOS.SetLogStd(&LogStd{
			Index:  ii,
			StdOut: stdOut,
			StdErr: stdErr,
			Error:  aerr.NewError(err),
		})
		cccOS.GetJRun().Logger().Info().Msgf("shell command %d of %d: finish", ii+1, len(ts.ShellCommands))
		if err != nil {
			return err
		}
	}
	return nil
}
