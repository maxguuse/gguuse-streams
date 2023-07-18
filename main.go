package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/maxguuse/gguuse-streams/configs/repositories"
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
	userToken := os.Getenv("USER_TOKEN")
	appToken := os.Getenv("APP_TOKEN")

	oauthAppToken := fmt.Sprintf("oauth:%s", appToken)

	twitchIrcClient := twitch.NewClient(username, oauthAppToken)
	log.Println("Twitch IRC client initialized")

	twitchApiClient, err := helix.NewClient(&helix.Options{
		UserAccessToken: userToken,
		AppAccessToken:  appToken,

		ClientID: os.Getenv("CLIENT_ID"),
	})
	if err != nil {
		panic(err)
	}
	twitch_config.ApiClient = twitchApiClient
	log.Println("Twitch API client initialized")

	twitch_config.IrcClient = twitchIrcClient
	twitch_config.ApiClient = twitchApiClient
	twitch_config.Channel = *channel

	log.Println("Twitch config sat up")

	commands := dataaccess.NewJsonCommandsRepository()
	err = commands.LoadCommands()
	if err != nil {
		panic(err)
	}
	log.Println("Commands source initialized")

	anns := dataaccess.NewJsonAnnouncementsRepository()
	err = anns.LoadAnnouncements()
	if err != nil {
		panic(err)
	}
	log.Println("Announcements source initialized")

	repositories.Commands = commands
	repositories.Announcements = anns

	announcements_helper.InitAnnouncements()

	twitch_config.IrcClient.Join(twitch_config.Channel)

	privateMessageHandler := handlers.NewPrivateMessageHandler()

	twitch_config.IrcClient.OnConnect(func() {
		log.Printf("Connected to #%s", twitch_config.Channel)
	})
	twitch_config.IrcClient.OnPrivateMessage(privateMessageHandler.Handle)

	err = twitch_config.IrcClient.Connect()
	if err != nil {
		panic(err)
	}
}
