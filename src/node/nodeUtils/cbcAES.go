package nodeUtils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
)


func deriveAES128KeyFromUint(num uint64) []byte {
	s := strconv.FormatUint(num, 10)
	hash := sha256.Sum256([]byte(s))
	return hash[:16] 
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	if padding == 0 {
		padding = blockSize
	}
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded data size")
	}

	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize || padding > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}

	for i := 0; i < padding; i++ {
		if data[len(data)-1-i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding bytes")
		}
	}

	return data[:len(data)-padding], nil
}

func EncryptStringAES128CBC(keyInt uint64, plaintext string) (string, error) {
	key := deriveAES128KeyFromUint(keyInt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("NewCipher: %w", err)
	}

	blockSize := block.BlockSize() 
	plainBytes := pkcs7Pad([]byte(plaintext), blockSize)

	cipherBytes := make([]byte, blockSize+len(plainBytes))

	iv := cipherBytes[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("IV generation: %w", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherBytes[blockSize:], plainBytes)

	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

func DecryptStringAES128CBC(keyInt uint64, b64Cipher string) (string, error) {
	key := deriveAES128KeyFromUint(keyInt)
	data, err := base64.StdEncoding.DecodeString(b64Cipher)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("NewCipher: %w", err)
	}

	blockSize := block.BlockSize()
	if len(data) < blockSize || len(data)%blockSize != 0 {
		return "", fmt.Errorf("invalid ciphertext length")
	}

	iv := data[:blockSize]
	ciphertext := data[blockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext) 

	plainBytes, err := pkcs7Unpad(ciphertext, blockSize)
	if err != nil {
		return "", fmt.Errorf("unpad: %w", err)
	}

	return string(plainBytes), nil
}


