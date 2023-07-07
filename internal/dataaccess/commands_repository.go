package dataaccess

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type command = string
type answer = string

type CommandsRepository interface {
	GetCommandResponse(command string) answer
	GetCommands() []command
	UpdateCommand(command string, answer string)
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

func (r *commandsRepository) LoadCommands() error {
	commandsFile, err := os.OpenFile(r.file, os.O_RDONLY|os.O_CREATE, 0644)
	defer func() {
		err := commandsFile.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		return err
	}

	byteValue, err := ioutil.ReadAll(commandsFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &r.commands)
	if err != nil {
		return err
	}

	return nil
}

func (r *commandsRepository) SaveCommands() error {
	commandsFile, err := os.Create(r.file)
	defer func() {
		err := commandsFile.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		return err
	}

	byteValue, err := json.Marshal(r.commands)
	if err != nil {
		return err
	}

	commandsFile.Write(byteValue)
	return nil
}
