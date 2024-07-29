package utils_test

import (
	"crypto/sha256"
	"io"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

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
	testPrivateKey := "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV1d0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktVd2dnU2hBZ0VBQW9JQkFRQ050YklMTlBJRmZ5SXUKa3FmMnZkSmZ4aXdFZi8rMFVrMkRrNkNCMDFoTWx4alE3YUF0dTBiS0RZV3VRZXNuYkE3bUt5R2VvQ3l6T2J0ago0a3A0dkdxYzB3NjZCU1VxUDBLeEJOSk5uNjJzMW1PVjNnbHY0TFNWVzRsZzJrZXpYdDAxdFZ4RFFEQnRkVWhHCm45Slh6NkN5RldNMkJ2UmdxNHBHSlp1bEJPQ0ladU0zdEIrcXA4L2ZiYzVTRktYRXBUaTNGYitiaCszQVp1YnYKVUZWTTlQOHFYY2cwUkdXVTVXYXpWbkY5Y3AyNGdFdEEyY3pCaW5VcWdjRFZqTTBOWWxjKzVvVzFhVXNvYmhicwo4UXEybzFRNUd1djZ5TitFam9kaXRUTFhlbFVnK201TEoyZFhGa3p4R0g2UlhIdnlrMWVLSVJ2cG9jL0pVTEpNCjdKbGk0T3NGQWdNQkFBRUNnZjlzNlBldEpVUGNkWmtQc2lia3UzNnpuTnEzbXFnckxoWGt5ZERSOWx3bWdQblIKbU05Q1ZteFJYWk1nR2dsZ2d1dndlYldjOC8xbXdUZ0R6Q2J3STk3TUtHbHBEZ3RDTE54VXNCL3hDSWV5RGhMNwpXMnBsVVkxNFBLR1lqaW9NOFJ1UjY1QzlIdGdaUjhvRWZWQnJyR2NVZHR1STZrOW0vRytJK2Q1bE5ScGJ1Wkp2CjMvMkxHNU0yMkhXRlVFYVYybGhkdzZIV0lhODhMZnpWSzltK0tCWG1YZnRYcDFTRUtxUXJmYWZnUzdDK3lYcnYKMWtaNHpjbGdYczd4M1dyR2NhSzVPYXY5T3d2dVhXRFJHb0ZwMGVicTU5UFRCZHUzRHBRbC9YMHdZY2RBdTVYUwoxUFNlckkwd1NjcEp0Y0FNMlprZXZGZXdCWk90dnp2eUdxLzlOTEVDZ1lFQXhaeS9hcUExZDlESW5lQzdYK2FkCjhiV1cxcUkwdEViZkNkbTBRMUJocTJ6U0VmbjlCVlFISmtaZkx3WStLT1pUenJuMjMzenlXVEg2czdkR0VSYUcKeVBEamFVTlUzRlpSWDBlSERqc0F5SlFxZ2MxaE53Tkh6cHlIbEpMUG40SjFxcldYT3JCK0JnNWgrODJzL1RzRApUbGtZWWVPbEwvMmZuT1VPSFVFUmx2VUNnWUVBdDVTRTV6MjIwejgzVjZhcWlLNngwVXJlT0ltNzhIVjZqL25BCjZoYVd6YVl2ejFQTTJCQVVYcENUVEdXejRiMFBQaHNoR3c1K1pvMU40cjhSUGVWNUhEeWdyNFZtQ0hoeTEwbTEKQ0g4ZGF5SEdJbEM3Qk9iZ015UEZXUk1RNUVqME9QQnRBWE5iT2JwQ3ZKRHcwKzNCU1UvVHdQUkZ5M0xLNVFNOQpudER6MmRFQ2dZRUFwWk13ejZadEpuZEpvUDhzQUs2NnFFditsdGhTVUxzUkpxL0MycVAvTWlONzRKUVY3T0Q3CkhKYmFLZ3lSQ0xQMGhNSk1sL1daR2lOR2JFNmo1cTE1UWVTVXB4NURmRnJXMDM2Ykt0RkZWc3JPMHZQREFOVSsKMVY5U09xcklURjZET1FYdU1MNncyV0l6dDBnZUtnL0lOVjF4a0pPdFZRaXRORWk0Q3NyNmNnRUNnWUJkOGhFNQpUU204WFVOekJZV0x3T3Fha2xlNlV6SHNVaEpRajUwYnJrZXFJZnVoZTk3K1N1eEJvSGJneDhNUUtISWVkRCt5CjJ1M3dpU3RzZUI3WXNCQVVWU3BkNkVSWTNWclh0WTZCTkp2WGNVYzExRjZBbWEzdVBjWUdXVzF4aGF3RlgzUkoKSThGeGYxSWJzWWlzeTZUNFlYT1o4T2V1djZYNUlIbHVScndqb1FLQmdCNVVsRjZUczF5YytuSG5id3Jvbjl0QwpXa1RIeFU4RWxoeVlzM3V3bkFmaGc3NTdwODVjWFUyclEydStiL3lzR2ovVzArUDZUanRWM1dDOWZHV2VTT1NWCm4vRGxpelJSUkw2eS9xOGpqRFZpV0dNdW92dWNZOTJJUlhNc0hiSUVUU0ZvSHZIalNrV0VuRDljWU82Z2tRVWcKejlQL0lUUm1EL3ZKVk9DaGhjbXEKLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQo="
	testPublicKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUFqYld5Q3pUeUJYOGlMcEtuOXIzUwpYOFlzQkgvL3RGSk5nNU9nZ2ROWVRKY1kwTzJnTGJ0R3lnMkZya0hySjJ3TzVpc2hucUFzc3ptN1krSktlTHhxCm5OTU91Z1VsS2o5Q3NRVFNUWit0ck5aamxkNEpiK0MwbFZ1SllOcEhzMTdkTmJWY1EwQXdiWFZJUnAvU1Y4K2cKc2hWak5nYjBZS3VLUmlXYnBRVGdpR2JqTjdRZnFxZlAzMjNPVWhTbHhLVTR0eFcvbTRmdHdHYm03MUJWVFBULwpLbDNJTkVSbGxPVm1zMVp4ZlhLZHVJQkxRTm5Nd1lwMUtvSEExWXpORFdKWFB1YUZ0V2xMS0c0VzdQRUt0cU5VCk9ScnIrc2pmaEk2SFlyVXkxM3BWSVBwdVN5ZG5WeFpNOFJoK2tWeDc4cE5YaWlFYjZhSFB5VkN5VE95Wll1RHIKQlFJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="
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
