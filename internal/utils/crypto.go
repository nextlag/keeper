package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
)

const keyLength = 32

func getKeyFromPass(keyString string) []byte {
	key := []byte(keyString)

	if len(key) < keyLength {
		for {
			key = append(key, key[0])
			if len(key) == keyLength {
				break
			}
		}
	} else if len(key) > keyLength {
		key = key[:keyLength]
	}

	return key
}

func Encrypt(keyString, stringToEncrypt string) string {
	if stringToEncrypt == "" {
		return stringToEncrypt
	}
	cipherBlock, err := aes.NewCipher(getKeyFromPass(keyString))
	if err != nil {
		log.Fatalf("Encrypt - aes.NewCipher - %v", err)
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		log.Fatalf("Encrypt - cipher.NewGCM - %v", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("Encrypt - io.ReadFull(rand.Reader, nonce) - %v", err)
	}

	return base64.URLEncoding.EncodeToString(aead.Seal(nonce, nonce, []byte(stringToEncrypt), nil))
}

func Decrypt(keyString, encryptedString string) (decryptedString string) {
	if encryptedString == "" {
		return encryptedString
	}
	encryptData, err := base64.URLEncoding.DecodeString(encryptedString)
	if err != nil {
		log.Fatal(err)
	}

	cipherBlock, err := aes.NewCipher(getKeyFromPass(keyString))
	if err != nil {
		log.Fatalf("Decrypt - aes.NewCipher - %v", err)
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		log.Fatalf("Decrypt - cipher.NewGCM - %v", err)
	}

	nonceSize := aead.NonceSize()
	if len(encryptData) < nonceSize {
		log.Fatalf("Decrypt - aead.NonceSize - %v", err)
	}

	nonce, cipherText := encryptData[:nonceSize], encryptData[nonceSize:]
	plainData, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Fatalf("Decrypt - aead.Open - %v", err)
	}

	return string(plainData)
}
