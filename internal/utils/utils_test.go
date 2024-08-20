package utils_test

import (
	"crypto/sha256"
	"io"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	config "github.com/nextlag/keeper/config/server"
	"github.com/nextlag/keeper/internal/utils"
)

const phrase = "This is top secret"

func TestCrypto(t *testing.T) {
	secretKey := "secretKey"
	encryptedString := utils.Encrypt(secretKey, phrase)
	decryptedString := utils.Decrypt(secretKey, encryptedString)

	if phrase != decryptedString {
		t.Errorf("got %q, wanted %q", decryptedString, phrase)
	}
}

func TestHash(t *testing.T) {
	password := "TestPassword"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Errorf("got %v error", err)
	}
	if utils.VerifyPassword(hashedPassword, password) != nil {
		t.Errorf("got %v error", err)
	}
}

func TestToken(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		require.Error(t, err)
	}
	testPrivateKey := cfg.Security.AccessTokenPrivateKey
	testPublicKey := cfg.Security.AccessTokenPublicKey
	testToken, err := utils.CreateToken(
		time.Hour,
		uuid.New(),
		testPrivateKey)
	if err != nil {
		t.Errorf("got %v error", err)
	}
	if _, err = utils.ValidToken(testToken, testPublicKey); err != nil {
		t.Errorf("got %v error", err)
	}
}

func TestCryptoFile(t *testing.T) {
	inputFilePath := "../../README.md"
	outputEncryptedFilePath := "../../encrypted_README.md"
	outputDecryptedFilePath := "../../decrypted_README.md"
	err := utils.EncryptFile(phrase, inputFilePath, outputEncryptedFilePath)
	if err != nil {
		t.Errorf("got %v error", err)
	}
	err = utils.DecryptFile(phrase, outputEncryptedFilePath, outputDecryptedFilePath)
	if err != nil {
		t.Errorf("got %v error", err)
	}

	hashes := make([]string, 2)
	for index, filePath := range []string{inputFilePath, outputDecryptedFilePath} {
		hashes[index], err = getFileHash(filePath)
		if err != nil {
			t.Errorf("got %v error", err)
		}
	}

	if !(hashes[0] == hashes[1]) {
		t.Errorf("files hashes are different:\n%s\n%s\n", hashes[0], hashes[1])
	}

	for _, filePath := range []string{outputEncryptedFilePath, outputDecryptedFilePath} {
		if err = os.Remove(filePath); err != nil {
			t.Errorf("got %v error", err)
		}
	}
}

func getFileHash(filePath string) (string, error) {
	inputData, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer inputData.Close()

	hash := sha256.New()
	if _, err = io.Copy(hash, inputData); err != nil {
		return "", err
	}

	return string(hash.Sum(nil)), nil
}
