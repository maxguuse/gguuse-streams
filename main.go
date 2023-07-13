package main

import (
	"flag"
	"log"
	"os"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/joho/godotenv"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
	"github.com/maxguuse/gguuse-streams/internal/handlers"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()

	username := os.Getenv("NICK")
	authToken := os.Getenv("AUTH_TOKEN")
	client := twitch.NewClient(username, authToken)

	channel := flag.String("channel", "gguuse", "Channel to join")
	flag.Parse()
	client.Join(*channel)

	commands := dataaccess.NewCommandsRepository()
	err := commands.LoadCommands()
	defer commands.SaveCommands()
	if err != nil {
		panic(err)
	}

	privateMessageHandler := handlers.NewPrivateMessageHandler(client, *channel, commands)

	client.OnPrivateMessage(privateMessageHandler.Handle)
	client.OnConnect(func() {
		log.Printf("Connected to #%s", *channel)
	})

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}
