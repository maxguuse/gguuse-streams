package twitch

import (
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nicklaw5/helix/v2"
)

var (
	IrcClient *twitch.Client
	ApiClient *helix.Client

	Channel string
)
