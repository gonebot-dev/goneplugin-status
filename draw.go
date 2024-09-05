package status

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"log"
)

func Draw() string {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	var imgBuffer bytes.Buffer
	err := png.Encode(&imgBuffer, img)
	if err != nil {
		log.Printf("[status] Draw status fail.\n")
		return ""
	}
	imgBase64 := base64.StdEncoding.EncodeToString(imgBuffer.Bytes())
	return imgBase64
}
