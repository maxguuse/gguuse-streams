package dataaccess

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type jsonCommandsRepository struct {
	commands map[string]string
	file     string
}

func NewJsonCommandsRepository(channel string) *jsonCommandsRepository {
	return &jsonCommandsRepository{
		commands: make(map[string]string),
		file:     fmt.Sprintf("json_commands/%s_commands.json", channel),
	}
}

func (r *jsonCommandsRepository) GetResponse(command string) string {
	return r.commands[command]
}

func (r *jsonCommandsRepository) GetCommands() (cmds []string) {
	for command := range r.commands {
		cmds = append(cmds, command)
	}
	return
}

func (r *jsonCommandsRepository) UpdateCommand(command string, answer string) {
	if answer == "" {
		delete(r.commands, command)
	} else {
		r.commands[command] = answer
	}
}

func (r *jsonCommandsRepository) LoadCommands() (err error) {
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	commandsFile, err := os.OpenFile(r.file, os.O_RDONLY|os.O_CREATE, 0644)
	byteValue, err := io.ReadAll(commandsFile)

	if !json.Valid(byteValue) {
		r.SaveCommands()
		byteValue, _ = io.ReadAll(commandsFile)
	}

	err = json.Unmarshal(byteValue, &r.commands)

	if err != nil {
		return err
	}
	defer commandsFile.Close()
	return nil
}

func (r *jsonCommandsRepository) SaveCommands() (err error) {
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
