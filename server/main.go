package main

import (
	"github.com/mattermost/mattermost-server/v6/plugin"

	"github.com/brightscout/mattermost-plugin-autolink/server/autolinkplugin"
)

func main() {
	plugin.ClientMain(autolinkplugin.New())
}
