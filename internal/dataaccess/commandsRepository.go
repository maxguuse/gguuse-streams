package dataaccess

type CommandsRepository interface {
	GetResponse(string) string
	GetCommands() []string
	UpdateCommand(string, string)

	LoadCommands() error
	SaveCommands() error
}
