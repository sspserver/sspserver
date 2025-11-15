package appcmd

import (
	"context"
	"fmt"
)

// BeforeCommandRunFunc is a function type that is called before executing a command.
type BeforeCommandRunFunc func(ctx context.Context, cmd ICommand) (context.Context, error)

// App represents the application with its commands and metadata.
type App struct {
	Name        string
	Description string
	Version     string
	BuildDate   string
	BuildCommit string
	CmdList     ICommands

	BeforeCommandRun BeforeCommandRunFunc
}

// Run executes the application with the provided arguments.
func (app *App) Run(ctx context.Context, args []string) error {
	if len(args) < 2 {
		app.printCommandsUsage()
		return nil
	}

	// Get command name
	cmdName := args[1]

	// Run command by name
	icmd := app.CmdList.Get(cmdName)

	// Print help if command not found
	if cmdName == "help" || icmd == nil {
		app.printCommandsUsage()
		return nil
	}

	if app.BeforeCommandRun != nil {
		var err error
		ctx, err = app.BeforeCommandRun(ctx, icmd)
		if err != nil {
			return fmt.Errorf("before command run: %w", err)
		}
	}

	// Execute command
	return icmd.Run(ctx, args[2:])
}

func (app *App) printCommandsUsage() {
	fmt.Printf("Usage: %s <command> [options]\n", app.Name)
	fmt.Printf("Version: %s\n", app.Version)
	fmt.Printf("Build Date: %s\n", app.BuildDate)
	fmt.Printf("Build Commit: %s\n", app.BuildCommit)
	fmt.Println()
	fmt.Printf("Description: –\n%s\n", app.Description)

	fmt.Println("Commands:")
	for _, cmd := range app.CmdList {
		fmt.Printf("  % 10s - %s\n", cmd.Cmd(), cmd.Help())
	}
	fmt.Println("  help       - print this help")
}
