package commands

type ICommands []ICommand

// Get command by name
func (c ICommands) Get(name string) ICommand {
	for _, cmd := range c {
		if cmd.Cmd() == name {
			return cmd
		}
	}
	return nil
}

// Commands list
var Commands = ICommands{}
