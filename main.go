package status

import (
	"fmt"

	"github.com/gonebot-dev/gonebot/messages"
	"github.com/gonebot-dev/gonebot/plugins"
)

func handler(msg messages.IncomingStruct) (result messages.ResultStruct) {
	info := GetSysInfo()
	result.Text = fmt.Sprintf("CPU usage: %f\n", info.CpuUsedPercent)
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
