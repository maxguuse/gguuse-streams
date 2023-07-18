package dataaccess

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	twitch_config "github.com/maxguuse/gguuse-streams/configs/twitch"

	"golang.org/x/exp/maps"
)

type jsonCommandsRepository struct {
	commands map[string]string
	file     string
}

func NewJsonCommandsRepository() *jsonCommandsRepository {
	return &jsonCommandsRepository{
		commands: make(map[string]string),
		file:     fmt.Sprintf("json_commands/%s_commands.json", twitch_config.Channel),
	}
}

func (r *jsonCommandsRepository) GetResponse(command string) string {
	return r.commands[command]
}

func (r *jsonCommandsRepository) GetCommands() (cmds []string) {
	return maps.Keys(r.commands)
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
