package dataaccess

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type command = string
type answer = string

type CommandsRepository interface {
	GetCommandResponse(string) answer
	GetCommands() []command
	UpdateCommand(string, string)
	LoadCommands() error
	SaveCommands() error
}

type commandsRepository struct {
	commands map[command]answer
	file     string
}

func NewCommandsRepository() *commandsRepository {
	return &commandsRepository{
		commands: make(map[command]answer),
		file:     "commands.json",
	}
}

func (r *commandsRepository) GetCommandResponse(command string) answer {
	return r.commands[command]
}

func (r *commandsRepository) GetCommands() []command {
	var commands []command
	for command := range r.commands {
		commands = append(commands, command)
	}
	return commands
}

func (r *commandsRepository) UpdateCommand(command string, answer string) {
	if answer == "" {
		delete(r.commands, command)
	} else {
		r.commands[command] = answer
	}
}

func (r *commandsRepository) LoadCommands() (err error) {
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	commandsFile, err := os.OpenFile(r.file, os.O_RDONLY|os.O_CREATE, 0644)
	byteValue, err := io.ReadAll(commandsFile)
	err = json.Unmarshal(byteValue, &r.commands)

	if err != nil {
		return err
	}
	defer commandsFile.Close()
	return nil
}

func (r *commandsRepository) SaveCommands() (err error) {
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	commandsFile, err := os.Create(r.file)
	byteValue, err := json.Marshal(r.commands)
	_, err = commandsFile.Write(byteValue)

	if err != nil {
		return err
	}
	defer commandsFile.Close()
	return nil
}
