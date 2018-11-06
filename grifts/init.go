package grifts

import (
	"github.com/dosaki/emote_combat_server/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
