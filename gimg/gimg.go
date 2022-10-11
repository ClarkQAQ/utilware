package gimg

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"utilware/dep/uuid"
	"utilware/dep/x/image/bmp"
)

type Gimg struct {
	x       int
	y       int
	Carrier *image.RGBA
}

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func New(x, y int, sourcePath string) (*Gimg, error) {
	// 创建结构
	g := &Gimg{
		x:       x,
		y:       y,
		Carrier: NewIMG(0, 0, x, y),
	}

	// 检查源地址
	if sourcePath == "" {
		return g, nil
	}

	// 解析图片
	bgFile, e := os.Open(sourcePath)
	if e != nil {
		return nil, fmt.Errorf("open file err: %v", e)
	}
	bgImage, e := imgDecode(bgFile)
	if e != nil {
		return nil, fmt.Errorf("decode image err: %v", e)
	}

	// 重新判断画布大小
	if x == 0 || y == 0 {
		g.Carrier = NewIMG(0, 0, bgImage.Bounds().Size().X, bgImage.Bounds().Size().Y)
	}

	// 合并
	MergeImage(g.Carrier, bgImage, bgImage.Bounds().Min.Sub(image.Point{0, 0}))
	return g, nil
}

func NewFromByte(x, y int, sourceData []byte) (*Gimg, error) {
	// 创建结构
	g := &Gimg{
		x:       x,
		y:       y,
		Carrier: NewIMG(0, 0, x, y),
	}

	// 检查源地址
	if sourceData == nil {
		return g, nil
	}

	bgImage, e := imgDecode(bytes.NewReader(sourceData))
	if e != nil {
		return nil, e
	}

	// 重新判断画布大小
	if x == 0 || y == 0 {
		g.Carrier = NewIMG(0, 0, bgImage.Bounds().Size().X, bgImage.Bounds().Size().Y)
	}

	// 合并
	MergeImage(g.Carrier, bgImage, bgImage.Bounds().Min.Sub(image.Point{0, 0}))
	return g, nil
}

func imgDecode(r io.Reader) (image.Image, error) {
	img, _, e := image.Decode(r)
	if e == nil {
		return img, nil
	}

	img, e = png.Decode(r)
	if e == nil {
		return img, nil
	}

	img, e = jpeg.Decode(r)
	if e == nil {
		return img, nil
	}

	img, e = bmp.Decode(r)
	if e == nil {
		return img, nil
	}

	img, e = gif.Decode(r)
	if e == nil {
		return img, nil
	}

	return nil, e
}

func (g *Gimg) Save(fileName ...string) error {
	if len(fileName) <= 0 || fileName[0] == "" {
		fileName = append(fileName, uuid.New().String()+".png")
	}

	data, e := g.SaveToBytes(nil)
	if e != nil {
		return e
	}

	ioutil.WriteFile(fileName[0], data, os.ModePerm)
	return nil
}

func (g *Gimg) SaveToBytes(data []byte) ([]byte, error) {
	merged := NewMerged(data)

	// 合并
	if e := Merge(g.Carrier, merged); e != nil {
		return nil, fmt.Errorf("merge image err: %v", e)
	}
	return merged.Bytes(), nil
}

func (g *Gimg) Clone() *Gimg {
	m := &Gimg{
		x: g.x,
		y: g.y,
	}

	m.Carrier = NewIMG(0, 0, g.Carrier.Bounds().Size().X, g.Carrier.Bounds().Size().Y)
	// 合并
	MergeImage(m.Carrier, g.Carrier, g.Carrier.Bounds().Min.Sub(image.Point{0, 0}))

	return m
}
