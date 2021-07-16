package crypc

import (
	"crypto/md5"
	"fmt"
)

func Md5Sum(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
