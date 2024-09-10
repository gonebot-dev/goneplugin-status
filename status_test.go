package status

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"testing"

	"github.com/fogleman/gg"
	"github.com/gonebot-dev/goneplugin-status/renderer"
	"github.com/gonebot-dev/goneplugin-status/sysinfo"
)

func TestSysinfo(t *testing.T) {
	info := sysinfo.GetSysInfo()
	result, _ := json.MarshalIndent(info, "", "  ")
	fmt.Println(string(result))
}

func TestRenderer(t *testing.T) {
	b64Str := renderer.Render()[9:]
	imgData, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		t.Fatal(err)
	}
	gg.SavePNG("test.png", img)
	//! Inspect it yourself
}
