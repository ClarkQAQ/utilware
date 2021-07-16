package crypc

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func Sha1Sum(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))
}

func Sha256Sum(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func Sha256Sum224(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum224(data))
}

func Sha512Sum(data []byte) string {
	return fmt.Sprintf("%x", sha512.Sum512(data))
}

func Sha512Sum384(data []byte) string {
	return fmt.Sprintf("%x", sha512.Sum384(data))
}
