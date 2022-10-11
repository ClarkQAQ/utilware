package gimg

import (
	"bytes"
	"image"
)

// MergeImage 合并图层
func (g *Gimg) DoImageBytes(x, y int, data []byte) error {
	img, _, e := image.Decode(bytes.NewBuffer(data))
	if e != nil {
		return e
	}

	MergeImage(g.Carrier, img, img.Bounds().Min.Sub(image.Point{x, y}))
	return nil
}

// MergeImage 合并图层
func (g *Gimg) DoImage(x, y int, img image.Image) {
	MergeImage(g.Carrier, img, img.Bounds().Min.Sub(image.Point{x, y}))
}
