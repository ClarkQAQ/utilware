package resize

import (
	"bytes"
	"image/jpeg"
)

func PicByteZip(data []byte, w, h uint) ([]byte, error) {
	img, e := jpeg.Decode(bytes.NewBuffer(data))
	if e != nil {
		return data, e
	}
	data_out := bytes.NewBuffer(nil)
	e = jpeg.Encode(data_out, Resize(1000, 0, img, Bilinear), nil)
	if e != nil {
		return data, e
	}
	return data_out.Bytes(), nil
}
