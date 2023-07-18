package main

import (
	"flag"
	"log"
	"os"

	twitch_config "github.com/maxguuse/gguuse-streams/configs/twitch"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/joho/godotenv"
	"github.com/maxguuse/gguuse-streams/internal/announcements/announcements_helper"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
	"github.com/maxguuse/gguuse-streams/internal/handlers"
	"github.com/nicklaw5/helix/v2"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	channel := flag.String("channel", "gguuse", "Channel to join")
	flag.Parse()

	log.Println("Starting bot...")

	loadEnv()
	log.Println("Environment variables loaded")

	username := os.Getenv("NICK")
	authToken := os.Getenv("AUTH_TOKEN")
	twitch_config.IrcClient = twitch.NewClient(username, authToken)
	twitch_config.IrcClient.Join(*channel)
	log.Println("Twitch IRC client initialized")

	twitchApiClient, err := helix.NewClient(&helix.Options{
		UserAccessToken: authToken[6:],
		AppAccessToken:  authToken[6:],
		ClientID:        os.Getenv("CLIENT_ID"),
	})
	if err != nil {
		panic(err)
	}
	twitch_config.ApiClient = twitchApiClient
	log.Println("Twitch API client initialized")

	commands := dataaccess.NewJsonCommandsRepository(*channel)
	err = commands.LoadCommands()
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

	announcements_helper.InitAnnouncements(anns, *channel)

	privateMessageHandler := handlers.NewPrivateMessageHandler(*channel, commands, anns)

	twitch_config.IrcClient.OnConnect(func() {
		log.Printf("Connected to #%s", *channel)
	})
	twitch_config.IrcClient.OnPrivateMessage(privateMessageHandler.Handle)

	err = twitch_config.IrcClient.Connect()
	if err != nil {
		panic(err)
	}
}
