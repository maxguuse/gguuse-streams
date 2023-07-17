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
	log.Println("Starting bot...")

	loadEnv()

	username := os.Getenv("NICK")
	authToken := os.Getenv("AUTH_TOKEN")
	client := twitch.NewClient(username, authToken)

	log.Println("Environment variables loaded")

	channel := flag.String("channel", "gguuse", "Channel to join")
	flag.Parse()
	client.Join(*channel)

	commands := dataaccess.NewJsonCommandsRepository(*channel)
	err := commands.LoadCommands()
	if err != nil {
		panic(err)
	}

	log.Println("Commands source initialized")

	anns := dataaccess.NewJsonAnnouncementsRepository(*channel)
	err = anns.LoadAnnouncements()
	if err != nil {
		panic(err)
	}

	log.Println("Announcements source initialized")

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
