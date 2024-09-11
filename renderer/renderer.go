package renderer

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"runtime"

	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/gonebot-dev/goneplugin-status/sysinfo"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
)

//go:embed assets
var assetsFS embed.FS

var bg image.Image
var contentFont font.Face
var titleFont font.Face

// Constants
const shadowOffsetX float64 = 10
const shadowOffsetY float64 = 10
const badgePaddingX float64 = 48
const badgePaddingY float64 = 24
const badgeMargin float64 = 32
const panelPadding float64 = 48
const panelMargin float64 = 48
const canvasWidth float64 = 1280

var golangBlue = "#007D9CC0"
var success = "#67C23AC0"
var warning = "#E6A23CC0"
var danger = "#F56C6CC0"
var shadow = "#00000070"
var panel = "#FFFFFF9C"

func init() {
	// Load background, asuming it to be 1280x...
	bgData, _ := assetsFS.Open("assets/background.png")
	defer bgData.Close()
	bg, _, _ = image.Decode(bgData)

	// Load font
	fontData, _ := assetsFS.ReadFile("assets/font.ttf")
	f, _ := freetype.ParseFont(fontData)
	titleFont = truetype.NewFace(f, &truetype.Options{
		Size:    64,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	contentFont = truetype.NewFace(f, &truetype.Options{
		Size:    36,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// Render renders the system info to an image and returns it as a base64 string
func Render() string {
	// Render process
	info := sysinfo.GetSysInfo()
	tmp := gg.NewContext(0, 0)
	tmp.SetFontFace(titleFont)
	_, titleLineHeight := tmp.MeasureString("T")
	tmp.SetFontFace(contentFont)
	_, contentLineHeight := tmp.MeasureString("T")
	canvasHeight := 0.0

	//! Define panels
	//* Panel for badges
	canvasHeight += (panelMargin + panelPadding) * 2
	//? Title badge
	canvasHeight += badgePaddingY*3 + titleLineHeight
	//? Adapter, Receive and Send badge
	canvasHeight += badgePaddingY*2 + contentLineHeight + badgeMargin
	//? Bot & System run time
	canvasHeight += badgePaddingY*2 + contentLineHeight + badgeMargin

	//* Panel for CPU usage
	canvasHeight += (panelPadding)*2 + panelMargin
	//? CPU title badge
	canvasHeight += badgePaddingY*2 + contentLineHeight
	//? CPU Info badge
	canvasHeight += badgePaddingY*2 + contentLineHeight + badgeMargin
	//? CPU progress bar
	canvasHeight += badgePaddingY + contentLineHeight + badgeMargin

	//* Panel for memory usage
	canvasHeight += (panelPadding)*2 + panelMargin
	//? Memory title badge
	canvasHeight += badgePaddingY*2 + contentLineHeight
	//? Memory progress bar
	canvasHeight += badgePaddingY + contentLineHeight + badgeMargin

	//* Panel for every disk
	for range info.Disks {
		//? Disk progress bar
		canvasHeight += (panelPadding)*2 + panelMargin + badgePaddingY*3 + contentLineHeight*2 + badgeMargin
	}

	//! Generate image
	img := gg.NewContext(int(canvasWidth), int(canvasHeight))
	img.SetFontFace(titleFont)

	//! Render Process
	//* Background
	tmpBg := bg
	if canvasHeight > canvasWidth {
		tmpBg = resize.Resize(0, uint(canvasHeight), bg, resize.Lanczos3)
	}
	img.DrawImageAnchored(tmpBg, int(canvasWidth/2), int(canvasHeight/2), 0.5, 0.5)

	//! Render Panels
	//* Panel for badges
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		panelMargin+shadowOffsetX,
		panelMargin+shadowOffsetY,
		canvasWidth-panelMargin*2,
		badgePaddingY*7+badgeMargin*2+contentLineHeight*2+titleLineHeight+panelPadding*2,
		32,
	)
	img.Fill()
	img.SetHexColor(panel)
	img.DrawRoundedRectangle(
		panelMargin,
		panelMargin,
		canvasWidth-panelMargin*2,
		badgePaddingY*7+badgeMargin*2+contentLineHeight*2+titleLineHeight+panelPadding*2,
		32,
	)
	img.Fill()
	//? Title badge
	logo := "\ue62a "
	if runtime.GOOS == "macos" {
		logo = "\uf179 "
	} else if runtime.GOOS == "linux" {
		logo = "\uf17c "
	}
	str := fmt.Sprintf("%sGonebot on %s %s", logo, info.OS, info.Arch)
	w, _ := img.MeasureString(str)
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+shadowOffsetX,
		panelMargin+panelPadding+shadowOffsetY,
		w+badgePaddingX*2,
		titleLineHeight+badgePaddingY*3,
		titleLineHeight/2.0+badgePaddingY*1.5,
	)
	img.Fill()
	img.SetHexColor(golangBlue)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding,
		panelMargin+panelPadding,
		w+badgePaddingX*2,
		titleLineHeight+badgePaddingY*3,
		titleLineHeight/2.0+badgePaddingY*1.5,
	)
	img.Fill()
	img.SetHexColor("#FFFFFF")
	img.DrawString(
		str,
		panelMargin+panelPadding+badgePaddingX,
		panelMargin+panelPadding+badgePaddingY*1.5+titleLineHeight,
	)
	//? Adapter, Receive and Send badge
	img.SetFontFace(contentFont)
	str = fmt.Sprintf("● %s", info.Backend)
	w, _ = img.MeasureString(str)
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+shadowOffsetX,
		panelMargin+panelPadding+badgePaddingY*3+titleLineHeight+badgeMargin+shadowOffsetY,
		w+badgePaddingX*2,
		contentLineHeight+badgePaddingY*2,
		contentLineHeight/2.0+badgePaddingY,
	)
	img.Fill()
	img.SetHexColor(success)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding,
		panelMargin+panelPadding+badgePaddingY*3+titleLineHeight+badgeMargin,
		w+badgePaddingX*2,
		contentLineHeight+badgePaddingY*2,
		contentLineHeight/2.0+badgePaddingY,
	)
	img.Fill()
	img.SetHexColor("#FFFFFF")
	img.DrawString(
		str,
		panelMargin+panelPadding+badgePaddingX,
		panelMargin+panelPadding+badgePaddingY*4+badgeMargin+titleLineHeight+contentLineHeight,
	)
	lastW := w + badgePaddingX*2 + badgeMargin
	str = fmt.Sprintf("● Recv: %d", info.ReceivedTotal)
	w, _ = img.MeasureString(str)
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+lastW+shadowOffsetX,
		panelMargin+panelPadding+badgePaddingY*3+titleLineHeight+badgeMargin+shadowOffsetY,
		w+badgePaddingX*2,
		contentLineHeight+badgePaddingY*2,
		contentLineHeight/2.0+badgePaddingY,
	)
	img.Fill()
	img.SetHexColor(warning)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+lastW,
		panelMargin+panelPadding+badgePaddingY*3+titleLineHeight+badgeMargin,
		w+badgePaddingX*2,
		contentLineHeight+badgePaddingY*2,
		contentLineHeight/2.0+badgePaddingY,
	)
	img.Fill()
	img.SetHexColor("#FFFFFF")
	img.DrawString(
		str,
		panelMargin+panelPadding+badgePaddingX+lastW,
		panelMargin+panelPadding+badgePaddingY*4+badgeMargin+titleLineHeight+contentLineHeight,
	)
	lastW += w + badgePaddingX*2 + badgeMargin
	str = fmt.Sprintf("● Sent: %d", info.SentTotal)
	w, _ = img.MeasureString(str)
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+lastW+shadowOffsetX,
		panelMargin+panelPadding+badgePaddingY*3+titleLineHeight+badgeMargin+shadowOffsetY,
		w+badgePaddingX*2,
		contentLineHeight+badgePaddingY*2,
		contentLineHeight/2.0+badgePaddingY,
	)
	img.Fill()
	img.SetHexColor(danger)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+lastW,
		panelMargin+panelPadding+badgePaddingY*3+titleLineHeight+badgeMargin,
		w+badgePaddingX*2,
		contentLineHeight+badgePaddingY*2,
		contentLineHeight/2.0+badgePaddingY,
	)
	img.Fill()
	img.SetHexColor("#FFFFFF")
	img.DrawString(
		str,
		panelMargin+panelPadding+badgePaddingX+lastW,
		panelMargin+panelPadding+badgePaddingY*4+badgeMargin+titleLineHeight+contentLineHeight,
	)
	//? Bot & System run time
	dayOrDays := "Day"
	if info.Days > 1 {
		dayOrDays = "Days"
	}
	str = fmt.Sprintf("Sys: %d %s %02d:%02d:%02d", info.Days, dayOrDays, info.Hours, info.Minutes, info.Seconds)
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+shadowOffsetX,
		panelMargin+panelPadding+badgePaddingY*5+titleLineHeight+contentLineHeight+badgeMargin*2+shadowOffsetY,
		(canvasWidth-badgeMargin-(panelPadding+panelMargin)*2)/2.0,
		contentLineHeight+badgePaddingY*2,
		contentLineHeight/2.0+badgePaddingY,
	)
	img.Fill()
	img.SetHexColor("#2222229C")
	img.DrawRoundedRectangle(
		panelMargin+panelPadding,
		panelMargin+panelPadding+badgePaddingY*5+titleLineHeight+contentLineHeight+badgeMargin*2,
		(canvasWidth-badgeMargin-(panelPadding+panelMargin)*2)/2.0,
		contentLineHeight+badgePaddingY*2,
		contentLineHeight/2.0+badgePaddingY,
	)
	img.Fill()
	img.SetHexColor("#FFFFFF")
	img.DrawStringAnchored(
		str,
		panelMargin+panelPadding+(canvasWidth-badgeMargin-(panelPadding+panelMargin)*2)/4.0,
		panelMargin+panelPadding+badgePaddingY*5+titleLineHeight+contentLineHeight*2+badgeMargin*2,
		0.5, 0.5,
	)
	if info.BotDays > 1 {
		dayOrDays = "Days"
	}
	str = fmt.Sprintf("Bot: %d %s %02d:%02d:%02d", info.BotDays, dayOrDays, info.BotHours, info.BotMinutes, info.BotSeconds)
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+(canvasWidth-badgeMargin-(panelPadding+panelMargin)*2)/2.0+badgeMargin+shadowOffsetX,
		panelMargin+panelPadding+badgePaddingY*5+titleLineHeight+contentLineHeight+badgeMargin*2+shadowOffsetY,
		(canvasWidth-badgeMargin-(panelPadding+panelMargin)*2)/2.0,
		contentLineHeight+badgePaddingY*2,
		contentLineHeight/2.0+badgePaddingY,
	)
	img.Fill()
	img.SetHexColor("#2222229C")
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+(canvasWidth-badgeMargin-(panelPadding+panelMargin)*2)/2.0+badgeMargin,
		panelMargin+panelPadding+badgePaddingY*5+titleLineHeight+contentLineHeight+badgeMargin*2,
		(canvasWidth-badgeMargin-(panelPadding+panelMargin)*2)/2.0,
		contentLineHeight+badgePaddingY*2,
		contentLineHeight/2.0+badgePaddingY,
	)
	img.Fill()
	img.SetHexColor("#FFFFFF")
	img.DrawStringAnchored(
		str,
		panelMargin+panelPadding+(canvasWidth-badgeMargin-(panelPadding+panelMargin)*2)*0.75+badgeMargin,
		panelMargin+panelPadding+badgePaddingY*5+titleLineHeight+contentLineHeight*2+badgeMargin*2,
		0.5, 0.5,
	)

	//* Panel for CPU usage
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		panelMargin+shadowOffsetX,
		panelMargin*2+panelPadding*2+badgePaddingY*7+titleLineHeight+contentLineHeight*2+badgeMargin*2+shadowOffsetY,
		canvasWidth-panelMargin*2,
		badgePaddingY*5+badgeMargin*2+contentLineHeight*3+panelPadding*2,
		32,
	)
	img.Fill()
	img.SetHexColor(panel)
	img.DrawRoundedRectangle(
		panelMargin,
		panelMargin*2+panelPadding*2+badgePaddingY*7+titleLineHeight+contentLineHeight*2+badgeMargin*2,
		canvasWidth-panelMargin*2,
		badgePaddingY*5+badgeMargin*2+contentLineHeight*3+panelPadding*2,
		32,
	)
	img.Fill()
	//? CPU title badge
	str = fmt.Sprintf("● CPU | Cores: %d", info.CpuCores)
	w, _ = img.MeasureString(str)
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		(canvasWidth-(w+badgePaddingX*2))/2.0+shadowOffsetX,
		panelMargin*2+panelPadding*3+badgePaddingY*7+titleLineHeight+contentLineHeight*2+badgeMargin*2+shadowOffsetY,
		w+badgePaddingX*2,
		badgePaddingY*2+contentLineHeight,
		badgePaddingY+contentLineHeight/2.0,
	)
	img.Fill()
	img.SetHexColor(danger)
	img.DrawRoundedRectangle(
		(canvasWidth-(w+badgePaddingX*2))/2.0,
		panelMargin*2+panelPadding*3+badgePaddingY*7+titleLineHeight+contentLineHeight*2+badgeMargin*2,
		w+badgePaddingX*2,
		badgePaddingY*2+contentLineHeight,
		badgePaddingY+contentLineHeight/2.0,
	)
	img.Fill()
	img.SetHexColor("#FFFFFF")
	img.DrawString(
		str,
		(canvasWidth-(w+badgePaddingX*2))/2.0+badgePaddingX,
		panelMargin*2+panelPadding*3+badgePaddingY*8+titleLineHeight+contentLineHeight*3+badgeMargin*2,
	)
	//? CPU Info badge
	str = fmt.Sprint(info.CpuInfo)
	w, _ = img.MeasureString(str)
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		(canvasWidth-(w+badgePaddingX*2))/2.0+shadowOffsetX,
		panelMargin*2+panelPadding*3+badgePaddingY*9+titleLineHeight+contentLineHeight*3+badgeMargin*3+shadowOffsetY,
		w+badgePaddingX*2,
		badgePaddingY*2+contentLineHeight,
		badgePaddingY+contentLineHeight/2.0,
	)
	img.Fill()
	img.SetHexColor(golangBlue)
	img.DrawRoundedRectangle(
		(canvasWidth-(w+badgePaddingX*2))/2.0,
		panelMargin*2+panelPadding*3+badgePaddingY*9+titleLineHeight+contentLineHeight*3+badgeMargin*3,
		w+badgePaddingX*2,
		badgePaddingY*2+contentLineHeight,
		badgePaddingY+contentLineHeight/2.0,
	)
	img.Fill()
	img.SetHexColor("#FFFFFF")
	img.DrawString(
		str,
		(canvasWidth-(w+badgePaddingX*2))/2.0+badgePaddingX,
		panelMargin*2+panelPadding*3+badgePaddingY*10+titleLineHeight+contentLineHeight*4+badgeMargin*3,
	)
	//? CPU progress bar
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+contentLineHeight*5+shadowOffsetX,
		panelMargin*2+panelPadding*3+badgePaddingY*11+titleLineHeight+contentLineHeight*4+badgeMargin*4+shadowOffsetY,
		canvasWidth-panelMargin*2-panelPadding*2-contentLineHeight*6,
		badgePaddingY+contentLineHeight,
		(badgePaddingY+contentLineHeight)/2.0,
	)
	img.Fill()
	img.SetHexColor(panel)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+contentLineHeight*5,
		panelMargin*2+panelPadding*3+badgePaddingY*11+titleLineHeight+contentLineHeight*4+badgeMargin*4,
		canvasWidth-panelMargin*2-panelPadding*2-contentLineHeight*6,
		badgePaddingY+contentLineHeight,
		(badgePaddingY+contentLineHeight)/2.0,
	)
	img.Fill()
	if info.CpuUsedPercent < 40 {
		img.SetHexColor(success)
	} else if info.CpuUsedPercent < 80 {
		img.SetHexColor(warning)
	} else {
		img.SetHexColor(danger)
	}
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+contentLineHeight*5,
		panelMargin*2+panelPadding*3+badgePaddingY*11+titleLineHeight+contentLineHeight*4+badgeMargin*4,
		(canvasWidth-panelMargin*2-panelPadding*2-contentLineHeight*6)*info.CpuUsedPercent/100.0,
		badgePaddingY+contentLineHeight,
		(badgePaddingY+contentLineHeight)/2.0,
	)
	img.Fill()
	img.SetHexColor("#000000")
	img.DrawString(
		fmt.Sprintf("%5.2f%%", info.CpuUsedPercent),
		panelMargin+panelPadding,
		panelMargin*2+panelPadding*3+badgePaddingY*11.5+titleLineHeight+contentLineHeight*5+badgeMargin*4,
	)
	img.SetHexColor("#000000")
	img.DrawStringAnchored(
		fmt.Sprintf("Load: %.2f / %.2f / %.2f", info.CpuLoad1, info.CpuLoad5, info.CpuLoad15),
		(canvasWidth+contentLineHeight*5)/2.0,
		panelMargin*2+panelPadding*3+badgePaddingY*11.5+titleLineHeight+contentLineHeight*4.5+badgeMargin*4,
		0.5, 0.5,
	)

	//* Panel for memory usage
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		panelMargin+shadowOffsetX,
		panelMargin*3+panelPadding*4+badgePaddingY*12+titleLineHeight+contentLineHeight*5+badgeMargin*4+shadowOffsetY,
		canvasWidth-panelMargin*2,
		badgePaddingY*3+badgeMargin*1+contentLineHeight*2+panelPadding*2,
		32,
	)
	img.Fill()
	img.SetHexColor(panel)
	img.DrawRoundedRectangle(
		panelMargin,
		panelMargin*3+panelPadding*4+badgePaddingY*12+titleLineHeight+contentLineHeight*5+badgeMargin*4,
		canvasWidth-panelMargin*2,
		badgePaddingY*3+badgeMargin*1+contentLineHeight*2+panelPadding*2,
		32,
	)
	img.Fill()
	//? Memory title badge
	str = fmt.Sprintf("● Memory | Total: %.2f GB", float64(info.MemAll)/1024.0)
	w, _ = img.MeasureString(str)
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		(canvasWidth-(w+badgePaddingX*2))/2.0+shadowOffsetX,
		panelMargin*3+panelPadding*5+badgePaddingY*12+titleLineHeight+contentLineHeight*5+badgeMargin*4+shadowOffsetY,
		w+badgePaddingX*2,
		badgePaddingY*2+contentLineHeight,
		badgePaddingY+contentLineHeight/2.0,
	)
	img.Fill()
	img.SetHexColor(warning)
	img.DrawRoundedRectangle(
		(canvasWidth-(w+badgePaddingX*2))/2.0,
		panelMargin*3+panelPadding*5+badgePaddingY*12+titleLineHeight+contentLineHeight*5+badgeMargin*4,
		w+badgePaddingX*2,
		badgePaddingY*2+contentLineHeight,
		badgePaddingY+contentLineHeight/2.0,
	)
	img.Fill()
	img.SetHexColor("#FFFFFF")
	img.DrawString(
		str,
		(canvasWidth-(w+badgePaddingX*2))/2.0+badgePaddingX,
		panelMargin*3+panelPadding*5+badgePaddingY*13+titleLineHeight+contentLineHeight*6+badgeMargin*4,
	)
	//? Memory progress bar
	img.SetHexColor(shadow)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+contentLineHeight*5+shadowOffsetX,
		panelMargin*3+panelPadding*5+badgePaddingY*14+titleLineHeight+contentLineHeight*6+badgeMargin*5+shadowOffsetY,
		canvasWidth-panelMargin*2-panelPadding*2-contentLineHeight*6,
		badgePaddingY+contentLineHeight,
		(badgePaddingY+contentLineHeight)/2.0,
	)
	img.Fill()
	img.SetHexColor(panel)
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+contentLineHeight*5,
		panelMargin*3+panelPadding*5+badgePaddingY*14+titleLineHeight+contentLineHeight*6+badgeMargin*5,
		canvasWidth-panelMargin*2-panelPadding*2-contentLineHeight*6,
		badgePaddingY+contentLineHeight,
		(badgePaddingY+contentLineHeight)/2.0,
	)
	img.Fill()
	if info.MemUsedPercent < 40 {
		img.SetHexColor(success)
	} else if info.MemUsedPercent < 80 {
		img.SetHexColor(warning)
	} else {
		img.SetHexColor(danger)
	}
	img.DrawRoundedRectangle(
		panelMargin+panelPadding+contentLineHeight*5,
		panelMargin*3+panelPadding*5+badgePaddingY*14+titleLineHeight+contentLineHeight*6+badgeMargin*5,
		(canvasWidth-panelMargin*2-panelPadding*2-contentLineHeight*6)*info.MemUsedPercent/100.0,
		badgePaddingY+contentLineHeight,
		(badgePaddingY+contentLineHeight)/2.0,
	)
	img.Fill()
	img.SetHexColor("#000000")
	img.DrawString(
		fmt.Sprintf("%5.2f%%", info.MemUsedPercent),
		panelMargin+panelPadding,
		panelMargin*3+panelPadding*5+badgePaddingY*14.5+titleLineHeight+contentLineHeight*7+badgeMargin*5,
	)
	img.SetHexColor("#000000")
	img.DrawStringAnchored(
		fmt.Sprintf("%.2f GB / %.2f GB", float64(info.MemUsed)/1024.0, float64(info.MemAll)/1024.0),
		(canvasWidth+contentLineHeight*5)/2.0,
		panelMargin*3+panelPadding*5+badgePaddingY*14.5+titleLineHeight+contentLineHeight*6.5+badgeMargin*5,
		0.5, 0.5,
	)

	//* Panel for every disk
	for i := range info.Disks {
		index := float64(i)
		img.SetHexColor(shadow)
		img.DrawRoundedRectangle(
			panelMargin+shadowOffsetX,
			panelMargin*(4+index)+panelPadding*(6+2*index)+badgePaddingY*(15+3*index)+titleLineHeight+contentLineHeight*(7+2*index)+badgeMargin*(5+index)+shadowOffsetY,
			canvasWidth-panelMargin*2,
			badgePaddingY*3+badgeMargin*1+contentLineHeight*2+panelPadding*2,
			32,
		)
		img.Fill()
		img.SetHexColor(panel)
		img.DrawRoundedRectangle(
			panelMargin,
			panelMargin*(4+index)+panelPadding*(6+2*index)+badgePaddingY*(15+3*index)+titleLineHeight+contentLineHeight*(7+2*index)+badgeMargin*(5+index),
			canvasWidth-panelMargin*2,
			badgePaddingY*3+badgeMargin*1+contentLineHeight*2+panelPadding*2,
			32,
		)
		img.Fill()
		//? Disk badge
		str = fmt.Sprintf("● Disk: \"%s\" | Total: %.2f GB", info.Disks[i].Name, float64(info.Disks[i].Total)/1024.0)
		w, _ = img.MeasureString(str)
		img.SetHexColor(shadow)
		img.DrawRoundedRectangle(
			(canvasWidth-(w+badgePaddingX*2))/2.0+shadowOffsetX,
			panelMargin*(4+index)+panelPadding*(7+2*index)+badgePaddingY*(15+3*index)+titleLineHeight+contentLineHeight*(7+2*index)+badgeMargin*(5+index)+shadowOffsetY,
			w+badgePaddingX*2,
			badgePaddingY*2+contentLineHeight,
			badgePaddingY+contentLineHeight/2.0,
		)
		img.Fill()
		img.SetHexColor(success)
		img.DrawRoundedRectangle(
			(canvasWidth-(w+badgePaddingX*2))/2.0,
			panelMargin*(4+index)+panelPadding*(7+2*index)+badgePaddingY*(15+3*index)+titleLineHeight+contentLineHeight*(7+2*index)+badgeMargin*(5+index),
			w+badgePaddingX*2,
			badgePaddingY*2+contentLineHeight,
			badgePaddingY+contentLineHeight/2.0,
		)
		img.Fill()
		img.SetHexColor("#FFFFFF")
		img.DrawString(
			str,
			(canvasWidth-(w+badgePaddingX*2))/2.0+badgePaddingX,
			panelMargin*(4+index)+panelPadding*(7+2*index)+badgePaddingY*(16+3*index)+titleLineHeight+contentLineHeight*(8+2*index)+badgeMargin*(5+index),
		)
		//? Disk progress bar
		img.SetHexColor(shadow)
		img.DrawRoundedRectangle(
			panelMargin+panelPadding+contentLineHeight*5+shadowOffsetX,
			panelMargin*(4+index)+panelPadding*(7+2*index)+badgePaddingY*(17+3*index)+titleLineHeight+contentLineHeight*(8+2*index)+badgeMargin*(6+index)+shadowOffsetY,
			canvasWidth-panelMargin*2-panelPadding*2-contentLineHeight*6,
			badgePaddingY+contentLineHeight,
			(badgePaddingY+contentLineHeight)/2.0,
		)
		img.Fill()
		img.SetHexColor(panel)
		img.DrawRoundedRectangle(
			panelMargin+panelPadding+contentLineHeight*5,
			panelMargin*(4+index)+panelPadding*(7+2*index)+badgePaddingY*(17+3*index)+titleLineHeight+contentLineHeight*(8+2*index)+badgeMargin*(6+index),
			canvasWidth-panelMargin*2-panelPadding*2-contentLineHeight*6,
			badgePaddingY+contentLineHeight,
			(badgePaddingY+contentLineHeight)/2.0,
		)
		img.Fill()
		if info.Disks[i].UsedPercent < 40 {
			img.SetHexColor(success)
		} else if info.Disks[i].UsedPercent < 80 {
			img.SetHexColor(warning)
		} else {
			img.SetHexColor(danger)
		}
		img.DrawRoundedRectangle(
			panelMargin+panelPadding+contentLineHeight*5,
			panelMargin*(4+index)+panelPadding*(7+2*index)+badgePaddingY*(17+3*index)+titleLineHeight+contentLineHeight*(8+2*index)+badgeMargin*(6+index),
			(canvasWidth-panelMargin*2-panelPadding*2-contentLineHeight*6)*info.Disks[i].UsedPercent/100.0,
			badgePaddingY+contentLineHeight,
			(badgePaddingY+contentLineHeight)/2.0,
		)
		img.Fill()
		img.SetHexColor("#000000")
		img.DrawString(
			fmt.Sprintf("%5.2f%%", info.Disks[i].UsedPercent),
			panelMargin+panelPadding,
			panelMargin*(4+index)+panelPadding*(7+2*index)+badgePaddingY*(17.5+3*index)+titleLineHeight+contentLineHeight*(9+2*index)+badgeMargin*(6+index),
		)
		img.SetHexColor("#000000")
		img.DrawStringAnchored(
			fmt.Sprintf("%.2f GB / %.2f GB", float64(info.Disks[i].Used)/1024.0, float64(info.Disks[i].Total)/1024.0),
			(canvasWidth+contentLineHeight*5)/2.0,
			panelMargin*(4+index)+panelPadding*(7+2*index)+badgePaddingY*(17.5+3*index)+titleLineHeight+contentLineHeight*(8.5+2*index)+badgeMargin*(6+index),
			0.5, 0.5,
		)
	}

	//! Calculate result
	var result bytes.Buffer
	writer := io.Writer(&result)
	img.EncodePNG(writer)

	return "base64://" + base64.StdEncoding.EncodeToString(result.Bytes())
}
