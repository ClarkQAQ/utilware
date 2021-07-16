package codec

import (
	"bytes"
	"encoding/gob"
)

func GobEncode(data interface{}) ([]byte, error) {
	var b bytes.Buffer
	e := gob.NewEncoder(&b).Encode(data)
	return b.Bytes(), e
}

func GobDecode(code []byte, data interface{}) error {
	var b bytes.Buffer
	b.Write(code)
	return gob.NewEncoder(&b).Encode(&data)
}

func GobEnDeCode(v interface{}) (data interface{}, e error) {
	var buf bytes.Buffer
	if e = gob.NewEncoder(&buf).Encode(v); e != nil {
		return nil, e
	}
	e = gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(data)
	return
}
