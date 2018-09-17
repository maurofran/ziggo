package zigbee

// Command is an alias for an empty interface
type Command interface{}

// CommandListener is the type of function receiving a command
type CommandListener interface {
	CommandReceived(Command)
}
