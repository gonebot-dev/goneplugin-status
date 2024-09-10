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
	canvasHeight := panelMargin * 2

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
	canvasHeight += badgePaddingY*3 + titleLineHeight
	//? CPU Info badge
	canvasHeight += badgePaddingY*2 + contentLineHeight + badgeMargin
	//? CPU progress bar
	canvasHeight += badgePaddingY*2 + contentLineHeight + badgeMargin

	//* Panel for memory usage
	canvasHeight += (panelPadding)*2 + panelMargin
	//? Memory title badge
	canvasHeight += badgePaddingY*3 + titleLineHeight
	//? Memory Info badge
	canvasHeight += badgePaddingY*2 + contentLineHeight + badgeMargin
	//? Memory progress bar
	canvasHeight += badgePaddingY*2 + contentLineHeight + badgeMargin

	//* Panel for every disk
	for range info.Disks {
		//? Disk progress bar
		canvasHeight += (panelPadding)*2 + panelMargin + badgePaddingY*2 + contentLineHeight
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
	str = fmt.Sprintf("Sys: %d %s %02d:%02d:%02d", info.SentTotal, dayOrDays, info.Hours, info.Minutes, info.Seconds)
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
	str = fmt.Sprintf("Bot: %d %s %02d:%02d:%02d", info.SentTotal, dayOrDays, info.Hours, info.Minutes, info.Seconds)
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
	//? CPU title badge
	//? CPU Info badge
	//? CPU progress bar
	if info.CpuUsedPercent < 40 {
		img.SetHexColor(success)
	} else if info.CpuUsedPercent < 80 {
		img.SetHexColor(warning)
	} else {
		img.SetHexColor(danger)
	}

	//* Panel for memory usage
	//? Memory title badge
	//? Memory Info badge
	//? Memory progress bar
	if info.CpuUsedPercent < 40 {
		img.SetHexColor(success)
	} else if info.CpuUsedPercent < 80 {
		img.SetHexColor(warning)
	} else {
		img.SetHexColor(danger)
	}

	//* Panel for every disk
	for range info.Disks {
		//? Disk progress bar
		canvasHeight += badgePaddingY*2 + contentLineHeight
	}

	//! Calculate result
	var result bytes.Buffer
	writer := io.Writer(&result)
	img.EncodePNG(writer)

	return "base64://" + base64.StdEncoding.EncodeToString(result.Bytes())
}
