package status

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gonebot-dev/goneplugin-status/sysinfo"
)

func TestSysinfo(t *testing.T) {
	info := sysinfo.GetSysInfo()
	result, _ := json.MarshalIndent(info, "", "  ")
	fmt.Println(string(result))
}
