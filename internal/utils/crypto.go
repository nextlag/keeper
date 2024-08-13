package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const keyLength = 32

var errEmptyFile = errors.New("empty file has been given")

// getKeyFromPass generates a 32-byte key from a provided string.
// If the input string is shorter than 32 bytes, it is padded.
// If it is longer, it is truncated to 32 bytes.
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

// Encrypt encrypts a string using AES-GCM with a key derived from keyString.
// The encrypted string is returned in base64 URL encoding format.
// If the input string is empty, it is returned as-is.
func Encrypt(keyString, stringToEncrypt string) string {
	if stringToEncrypt == "" {
		return stringToEncrypt
	}
	cipherBlock, err := aes.NewCipher(getKeyFromPass(keyString))
	if err != nil {
		log.Printf("Encrypt - aes.NewCipher - %v", err)
		return ""
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		log.Printf("Encrypt - cipher.NewGCM - %v", err)
		return ""
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("Encrypt - io.ReadFull(rand.Reader, nonce) - %v", err)
		return ""
	}

	return base64.URLEncoding.EncodeToString(aead.Seal(nonce, nonce, []byte(stringToEncrypt), nil))
}

// Decrypt decrypts a base64 URL encoded AES-GCM encrypted string.
// The decrypted string is returned. If the input string is empty, it is returned as-is.
func Decrypt(keyString, encryptedString string) (decryptedString string) {
	if encryptedString == "" {
		return encryptedString
	}
	encryptData, err := base64.URLEncoding.DecodeString(encryptedString)
	if err != nil {
		log.Println(err)
		return
	}

	cipherBlock, err := aes.NewCipher(getKeyFromPass(keyString))
	if err != nil {
		log.Printf("Decrypt - aes.NewCipher - %v", err)
		return ""
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		log.Printf("Decrypt - cipher.NewGCM - %v", err)
		return ""
	}

	nonceSize := aead.NonceSize()
	if len(encryptData) < nonceSize {
		log.Printf("Decrypt - aead.NonceSize - %v", err)
		return ""
	}

	nonce, cipherText := encryptData[:nonceSize], encryptData[nonceSize:]
	plainData, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Printf("Decrypt - aead.Open - %v", err)
		return ""
	}

	return string(plainData)
}

// EncryptFile encrypts the contents of a file using AES-GCM and writes the result
// to a new file. The key for encryption is derived from keyString.
// If the input file is empty, an error is returned.
func EncryptFile(keyString, inputFilePath, outputFilePath string) error {
	fi, err := os.Stat(inputFilePath)
	if err != nil {
		return fmt.Errorf("EncryptFile - os.Stat - %w", err)
	}
	if size := fi.Size(); size == 0 {
		return fmt.Errorf("EncryptFile - fi.Size - %w", errEmptyFile)
	}
	cipherBlock, err := aes.NewCipher(getKeyFromPass(keyString))
	if err != nil {
		return fmt.Errorf("EncryptFile - aes.NewCipher - %w", err)
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return fmt.Errorf("EncryptFile - cipher.NewGCM - %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("EncryptFile - io.ReadFull(rand.Reader, nonce) - %w", err)
	}
	inputData, err := os.ReadFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("EncryptFile - os.ReadFile - %w", err)
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("EncryptFile - os.Create - %w", err)
	}

	defer outputFile.Close()

	_, err = outputFile.WriteString(
		base64.URLEncoding.EncodeToString(aead.Seal(nonce, nonce, inputData, nil)))
	if err != nil {
		return fmt.Errorf("EncryptFile - base64.URLEncoding.EncodeToString - %w", err)
	}

	return nil
}

// DecryptFile decrypts a base64 URL encoded AES-GCM encrypted file and writes
// the result to a new file. The key for decryption is derived from keyString.
func DecryptFile(keyString, encryptedPath, decryptedFilePath string) error {
	encryptedData, err := os.ReadFile(encryptedPath)
	if err != nil {
		return fmt.Errorf("DecryptFile - os.ReadFile - %w", err)
	}

	encryptData, err := base64.URLEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return fmt.Errorf("DecryptFile - base64.URLEncoding.DecodeString - %w", err)
	}

	cipherBlock, err := aes.NewCipher(getKeyFromPass(keyString))
	if err != nil {
		return fmt.Errorf("DecryptFile - aes.NewCipher - %w", err)
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return fmt.Errorf("DecryptFile - cipher.NewGCM - %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(encryptData) < nonceSize {
		return fmt.Errorf("DecryptFile - aead.NonceSize - %w", err)
	}

	nonce, cipherText := encryptData[:nonceSize], encryptData[nonceSize:]
	plainData, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return fmt.Errorf("DecryptFile - aead.Open - %w", err)
	}
	outputFile, err := os.Create(decryptedFilePath)
	if err != nil {
		return fmt.Errorf("DecryptFile - os.Create - %w", err)
	}

	defer outputFile.Close()
	_, err = outputFile.WriteString(
		string(plainData))
	if err != nil {
		return fmt.Errorf("DecryptFile - outputFile.WriteString - %w", err)
	}

	return nil
}
