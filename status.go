package status

import (
	"reflect"

	"github.com/gonebot-dev/gonebot/adapter"
	"github.com/gonebot-dev/gonebot/message"
	"github.com/gonebot-dev/gonebot/plugin"
	"github.com/gonebot-dev/gonebot/rule"
	"github.com/gonebot-dev/goneplugin-status/renderer"
)

func handler(a *adapter.Adapter, msg message.Message) bool {
	imageSerializer := message.GetSerializer("image", "")
	imageSegment := message.MessageSegment{
		Type: "image",
		Data: imageSerializer.Serialize(message.ImageType{
			File: renderer.Render(a),
		}, reflect.TypeOf(imageSerializer)),
	}
	var result message.Message
	result.Sender = msg.Self
	result.Receiver = msg.Sender
	result.Group = msg.Group
	result.Self = msg.Self
	result.Attach(imageSegment, imageSerializer)
	a.SendChannel.Push(result, false)
	return true
}

var Status plugin.GonePlugin

func init() {
	Status.Name = "Status"
	Status.Version = "v0.1.0"
	Status.Description = "Show bot status"

	Status.Handlers = append(Status.Handlers, plugin.GoneHandler{
		Rules:   []rule.FilterBundle{{Filters: []rule.FilterRule{rule.ToMe(), rule.Command([]string{"status"})}}},
		Handler: handler,
	})
}
