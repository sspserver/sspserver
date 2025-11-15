package appcmd

import (
	"context"

	"github.com/demdxx/goconfig"
)

// ICommand is a command interface
type ICommand interface {
	String() string
	Cmd() string
	Help() string
	Run(ctx context.Context, args []string) error
}

// CommandFunc is a function that can be executed by the command line
type CommandFunc[T any] func(ctx context.Context, args []string, config *T) error

// ContextInitFunc is a function that initializes the context for a command.
type ContextInitFunc func(ctx context.Context) (context.Context, error)

// Command is a command that can be executed by the command line
type Command[T any] struct {
	Name        string
	HelpDesc    string
	Exec        CommandFunc[T]
	ContextInit ContextInitFunc
}

// Name of the command
func (c *Command[T]) String() string {
	return c.Name
}

// Cmd returns the command name
func (c *Command[T]) Cmd() string {
	return c.Name
}

// Name of the command
func (c *Command[T]) Help() string {
	return c.HelpDesc
}

// Run the command with the given context and arguments
func (c *Command[T]) Run(ctx context.Context, args []string) error {
	var config T
	// Parse config from args and environment
	err := goconfig.Load(
		&config,
		goconfig.WithDefaults(),
		goconfig.WithEnv(),
		goconfig.WithCustomArgs(args...),
	)
	if err != nil {
		return err
	}
	if c.ContextInit != nil {
		if ctx, err = c.ContextInit(ctx); err != nil {
			return err
		}
	}
	return c.Exec(ctx, args, &config)
}

func (c *Command[T]) WithInitContext(ctxWrapper ContextInitFunc) *Command[T] {
	newCmd := *c
	newCmd.ContextInit = ctxWrapper
	return &newCmd
}
