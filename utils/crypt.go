package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

func AesEncryptCFB(origData []byte, key []byte) (encrypted string, err error) {
	paddedOrigData := string(origData)
	remainder := len(paddedOrigData) % 16
	if remainder > 0 {
		paddedOrigData = paddedOrigData + strings.Repeat(" ", 16-remainder)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	tmpBs := make([]byte, len(paddedOrigData))
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(tmpBs, []byte(paddedOrigData))

	encryptedBs := append(iv, tmpBs[:len(origData)]...)
	encrypted = base64.StdEncoding.EncodeToString(encryptedBs)
	return
}

func AesDecryptCFB(encrypted string, key []byte) (decrypted string, err error) {
	encryptedBs, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	origData := encryptedBs[aes.BlockSize:]
	origDataLen := len(origData)
	paddedOrigData := origData
	remainder := len(paddedOrigData) % 16
	if remainder > 0 {
		paddedOrigData = append(paddedOrigData, []byte(strings.Repeat(" ", 16-remainder))...)
	}

	block, _ := aes.NewCipher(key)
	if len(paddedOrigData) < aes.BlockSize {
		err = fmt.Errorf("aes cipher text too short: %s", string(encrypted))
		return
	}

	iv := make([]byte, 0)
	iv = append(iv, encryptedBs[:aes.BlockSize]...)
	tmpBs := make([]byte, len(paddedOrigData))

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(tmpBs, paddedOrigData)

	decryptedBs := tmpBs[0:origDataLen]
	decrypted = string(decryptedBs)
	return
}
