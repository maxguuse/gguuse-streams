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

	commands := dataaccess.NewJsonCommandsRepository(*channel)
	err := commands.LoadCommands()
	if err != nil {
		panic(err)
	}

	anns := dataaccess.NewJsonAnnouncmentsRepository(*channel)
	err = anns.LoadAnnouncments()
	if err != nil {
		panic(err)
	}

	privateMessageHandler := handlers.NewPrivateMessageHandler(client, *channel, commands, anns)

	client.OnConnect(func() {
		log.Printf("Connected to #%s", *channel)
	})
	client.OnPrivateMessage(privateMessageHandler.Handle)

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}
