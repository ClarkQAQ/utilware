package gimg

import (
	"fmt"
	"strconv"

	"utilware/dep/freetype/truetype"
)

type TextHandler struct {
	g     *Gimg
	dSize float64
	dRGB  Color
	dFont *truetype.Font
	dText string
}

type Text struct {
	X    int
	Y    int
	Size float64
	RGB  string
	Text string
	Font *truetype.Font
}

func hex2rgb(color string) (c Color, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("%v", r)
		}
	}()

	ri, _ := strconv.ParseInt(color[:2], 16, 10)
	gi, _ := strconv.ParseInt(color[2:4], 16, 18)
	bi, _ := strconv.ParseInt(color[4:], 16, 10)
	return Color{uint8(ri), uint8(gi), uint8(bi), 0}, nil
}

// 创建文字工作间
func (g *Gimg) DoText(size float64, color, text string, font *truetype.Font) *TextHandler {
	if size == 0 {
		size = 26
	}

	c, e := hex2rgb(color)
	if e != nil {
		c = Color{100, 100, 100, 0}
	}

	return &TextHandler{
		g:     g,
		dSize: size,
		dRGB:  c,
		dFont: font,
		dText: text,
	}
}

// 合并文字
func (h *TextHandler) Text(t Text) (e error) {
	if t.Size == 0 {
		t.Size = h.dSize
	}

	if t.Text == "" {
		t.Text = h.dText
	}

	// 处理颜色
	color := h.dRGB
	if t.RGB != "" {
		if color, e = hex2rgb(t.RGB); e != nil {
			return fmt.Errorf("hex2rgb err: %v", e)
		}
	}

	// 处理字体
	if t.Font == nil {
		t.Font = h.dFont
	}

	dText := NewDrawText(h.g.Carrier)

	//设置颜色
	dText.SetColor(color.R, color.G, color.B)

	// 合并
	e = dText.MergeText(t.Text, t.Size, t.Font, t.X, t.Y)
	if e != nil {
		return fmt.Errorf("merge text err: %v", e)
	}
	return nil
}
