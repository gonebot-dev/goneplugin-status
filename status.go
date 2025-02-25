package status

import (
	"github.com/gonebot-dev/gonebot/message"
	"github.com/gonebot-dev/gonebot/plugin"
	"github.com/gonebot-dev/gonebot/plugin/handler"
	"github.com/gonebot-dev/goneplugin-status/renderer"
)

var TriggerCommand = "status"

var Status plugin.GonePlugin

func statusHandler(incomingMsg message.Message, resultMsg *message.Message) bool {
	resultMsg.AddImageSegment(renderer.Render())
	return true
}

func statusMatcher(incomingMsg message.Message) bool {
	return incomingMsg.IsToMe && incomingMsg.HasPrefix(TriggerCommand)
}

func init() {
	Status.Name = "Status"
	Status.Version = "v0.2.0"
	Status.Description = "Show bot status"

	Status.Handlers = append(Status.Handlers, handler.GoneHandler{
		Matcher: statusMatcher,
		Handler: statusHandler,
	})
}
