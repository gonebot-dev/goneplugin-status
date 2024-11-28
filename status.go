package status

import (
	"github.com/gonebot-dev/gonebot/adapter"
	"github.com/gonebot-dev/gonebot/message"
	"github.com/gonebot-dev/gonebot/plugin"
	"github.com/gonebot-dev/gonebot/plugin/rule"
	"github.com/gonebot-dev/goneplugin-status/renderer"
)

func handler(a *adapter.Adapter, msg message.Message) bool {
	result := message.NewReply(msg).Image(renderer.Render(a))
	a.SendMessage(result)
	return true
}

var Status plugin.GonePlugin

func init() {
	Status.Name = "Status"
	Status.Version = "v0.1.0"
	Status.Description = "Show bot status"

	Status.Handlers = append(Status.Handlers, plugin.GoneHandler{
		Rules:   rule.NewRules(rule.ToMe()).And(rule.Command("status")),
		Handler: handler,
	})
}
