package gob4aes

import (
	"utilware/util/codec"
	"utilware/util/crypc"

	"utilware/arpc"
)

type mpkv struct {
	Body   []byte
	Values map[interface{}]interface{}
}

// MsgPack represents a gzip coding middleware.
type GobForAes struct {
	aesPassword []byte
}

// Encode implements arpc MessageCoder.
func (mp *GobForAes) Encode(client *arpc.Client, msg *arpc.Message) *arpc.Message {
	body := msg.Data()
	v := &mpkv{
		Body:   body,
		Values: msg.Values(),
	}

	b, e := func(val *mpkv) ([]byte, error) {
		b, e := codec.GobEncode(val)
		if e != nil {
			return nil, e
		}

		b, e = crypc.AesCBCEncrypter(b, mp.aesPassword)
		if e != nil {
			return nil, e
		}

		return b, nil
	}(v)

	if e != nil {
		msg.SetError(true)
		b = []byte(e.Error())
	}

	ml := msg.MethodLen()
	msg.Buffer = append(msg.Buffer[:arpc.HeadLen+ml], b...)
	msg.SetBodyLen(ml + len(b))

	return msg
}

// Decode implements arpc MessageCoder.
func (mp *GobForAes) Decode(client *arpc.Client, msg *arpc.Message) *arpc.Message {

	v := &mpkv{
		Body:   msg.Data(),
		Values: msg.Values(),
	}

	if e := func(data []byte, val *mpkv) error {
		crypted, e := crypc.AesCBCDecrypter(data, mp.aesPassword)
		if e != nil {
			return e
		}

		return codec.GobDecode(crypted, val)
	}(msg.Data(), v); e != nil {
		msg.SetError(true)
		v.Body = []byte(e.Error())
	}

	ml := msg.MethodLen()
	msg.Buffer = append(msg.Buffer[:arpc.HeadLen+ml], v.Body...)
	msg.SetBodyLen(ml + len(v.Body))
	for k, v := range v.Values {
		msg.Set(k, v)
	}

	return msg
}

// New returns the MsgPack coding middleware.
func New(k []byte) *GobForAes {
	return &GobForAes{
		aesPassword: k,
	}
}
