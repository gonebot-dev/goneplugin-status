package renderer

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/gonebot-dev/goneplugin-status/sysinfo"
	"golang.org/x/image/font"
)

//go:embed assets
var assetsFS embed.FS

var bg image.Image
var myfont font.Face

func init() {
	// Load background, asuming it to be 1280x1280
	bgData, err := assetsFS.Open("assets/background.png")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer bgData.Close()
	bg, _, err = image.Decode(bgData)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Load font
	fontData, err := assetsFS.ReadFile("assets/font.ttf")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	f, err := freetype.ParseFont(fontData)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	myfont = truetype.NewFace(f, &truetype.Options{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// Render renders the system info to an image and returns it as a base64 string
func Render() string {
	// Render process
	info := sysinfo.GetSysInfo()
	img := gg.NewContextForImage(bg)
	img.SetFontFace(myfont)

	// Panel1 shadow
	img.SetRGBA255(0, 0, 0, 112)
	img.DrawRoundedRectangle(58, 58, 1184, 440, 64)
	img.Fill()

	// Panel1
	img.SetRGBA255(255, 255, 255, 156)
	img.DrawRoundedRectangle(48, 48, 1184, 440, 64)
	img.Fill()

	// Measure badge text width
	str := fmt.Sprintf("● Gonebot on %s %s", info.OS, info.Arch)
	w, h := img.MeasureString(str)

	// Bot badge shadow
	img.SetRGBA255(0, 0, 0, 112)
	img.DrawRoundedRectangle(110, 110, 116+w, h+64, 50)
	img.Fill()

	// Bot badge
	img.SetRGBA255(0, 125, 156, 192)
	img.DrawRoundedRectangle(100, 100, 116+w, h+64, 50)
	img.Fill()

	// Bot badge text
	img.SetRGBA(1, 1, 1, 1)
	img.DrawStringAnchored(str, 158, 132+h/2, 0, 0.5)
	img.Fill()

	// Measure badge text width
	if info.Days > 1 {
		str = fmt.Sprintf("● Time elapsed since system boot:\n%d Days %02d:%02d:%02d", info.Days, info.Hours, info.Minutes, info.Seconds)
	} else {
		str = fmt.Sprintf("● Time elapsed since system boot:\n%d Day %02d:%02d:%02d", info.Days, info.Hours, info.Minutes, info.Seconds)
	}
	w, h = img.MeasureMultilineString(str, 2)

	// Running time shadow
	img.SetRGBA255(0, 0, 0, 112)
	img.DrawRoundedRectangle(110, 258, 116+w, h+64, 50)
	img.Fill()

	// Running time
	img.SetRGBA255(103, 194, 58, 192)
	img.DrawRoundedRectangle(100, 248, 116+w, h+64, 50)
	img.Fill()

	// Running time text
	img.SetRGBA(0, 0, 0, 1)
	img.DrawStringWrapped(str, 158, 280+h/2, 0, 0.5, w, 2, gg.AlignCenter)
	img.Fill()

	// Panel2 shadow
	img.SetRGBA255(0, 0, 0, 112)
	img.DrawRoundedRectangle(58, 546, 1184, 690, 64)
	img.Fill()

	// Panel2
	img.SetRGBA255(255, 255, 255, 156)
	img.DrawRoundedRectangle(48, 536, 1184, 690, 64)
	img.Fill()

	// Dashboard witdh
	width := 28.0

	// CPU dashboard
	img.SetRGBA(0, 0, 0, 0.3)
	for i := 0.0; i < width; i++ {
		img.DrawCircle(224, 728, 128-i)
		img.Stroke()
	}
	if info.CpuUsedPercent < 50 {
		img.SetRGBA255(103, 194, 58, 192)
	} else if info.CpuUsedPercent < 80 {
		img.SetRGBA255(230, 162, 60, 192)
	} else {
		img.SetRGBA255(245, 108, 108, 192)
	}
	for i := 0.0; i < width; i++ {
		img.DrawArc(224, 728, 128-i, gg.Radians(-90.0), gg.Radians(-90.0+3.6*info.CpuUsedPercent))
		img.Stroke()
	}

	// CPU dashboard text
	img.SetRGBA(0, 0, 0, 1)
	str = fmt.Sprintf("CPU: %d%%\nCores: %d", int8(math.Round(info.CpuUsedPercent)), info.CpuCores)
	w, _ = img.MeasureMultilineString(str, 2)
	img.DrawStringWrapped(str, 388, 750, 0, 0.5, w, 2, gg.AlignLeft)

	// Memory dashboard
	img.SetRGBA(0, 0, 0, 0.3)
	for i := 0.0; i < width; i++ {
		img.DrawCircle(224, 1034, 128-i)
		img.Stroke()
	}
	if info.MemUsedPercent < 30 {
		img.SetRGBA255(103, 194, 58, 192)
	} else if info.MemUsedPercent < 80 {
		img.SetRGBA255(230, 162, 60, 192)
	} else {
		img.SetRGBA255(245, 108, 108, 192)
	}
	for i := 0.0; i < width; i++ {
		img.DrawArc(224, 1034, 128-i, gg.Radians(-90.0), gg.Radians(-90.0+3.6*info.MemUsedPercent))
		img.Stroke()
	}

	// Memory dashboard text
	img.SetRGBA(0, 0, 0, 1)
	str = fmt.Sprintf("Memory: %d%%\n%.1f/%.1f GB", int8(math.Round(info.MemUsedPercent)), float32(info.MemUsed)/1024.0, float32(info.MemAll)/1024.0)
	w, _ = img.MeasureMultilineString(str, 2)
	img.DrawStringWrapped(str, 388, 1056, 0, 0.5, w, 2, gg.AlignLeft)

	var result bytes.Buffer
	writer := io.Writer(&result)
	img.EncodePNG(writer)

	return "base64://" + base64.StdEncoding.EncodeToString(result.Bytes())
}
