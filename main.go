package status

import (
	renderer "github.com/Kingcxp/go-sysinfo-renderer"
	"github.com/gonebot-dev/gonebot/messages"
	"github.com/gonebot-dev/gonebot/plugins"
)

func handler(msg messages.IncomingStruct) (result messages.ResultStruct) {
	result.Imgs = append(result.Imgs, renderer.Render())
	return result
}

func init() {
	status := plugins.GonePlugin{}
	status.Name = "status"
	status.Description = "Picture status"
	status.Handlers = append(status.Handlers,
		plugins.GoneHandler{
			Command: []string{"status", "状态", "stat"},
			Handler: handler,
		})
	plugins.LoadPlugin(status)
}
