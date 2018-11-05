package grifts

import (
	"github.com/dosaki/owl_power_server/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
