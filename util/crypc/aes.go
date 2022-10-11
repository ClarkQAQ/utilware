package crypc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func AesCBCEncrypter(data, key []byte) ([]byte, error) {
	block, e := aes.NewCipher(key)
	if e != nil {
		return nil, e
	}

	length := len(data)
	if length < 1 {
		return nil, errors.New("data is empty")
	}

	blockSize := block.BlockSize()

	padding := blockSize - length%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	encryptBytes := append(data, padText...)

	crypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	blockMode.CryptBlocks(crypted, encryptBytes)

	return crypted, nil
}

func AesCBCDecrypter(data, key []byte) ([]byte, error) {
	block, e := aes.NewCipher(key)
	if e != nil {
		return nil, e
	}

	blockSize := block.BlockSize()

	crypted := make([]byte, len(data))
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	blockMode.CryptBlocks(crypted, data)

	length := len(crypted)
	if length < 1 {
		return nil, errors.New("data is empty")
	}

	offset := int(crypted[length-1])
	if offset > length {
		return nil, errors.New("data is invalid")
	}

	return crypted[:(length - offset)], nil
}
