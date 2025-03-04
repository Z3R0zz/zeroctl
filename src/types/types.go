package types

type Command struct {
	Name        string
	Description string
	Handler     func() string
}

var CommandRegistry = map[string]Command{}

func RegisterCommand(cmd Command) {
	CommandRegistry[cmd.Name] = cmd
}

func GetCommand(name string) (Command, bool) {
	cmd, exists := CommandRegistry[name]
	return cmd, exists
}
