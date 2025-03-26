package ashell

// IShellCommand defines the interface for shell commands.
type IShellCommand interface {
	GetType() ShellType
	RunCommand() (string, string, error)
	RunCommandWithOptions(hideStdOut bool, isFakeRun bool) (string, string, error)
	GetCommand() string
	GetParameters() []string
}

type IShellCommands []IShellCommand
