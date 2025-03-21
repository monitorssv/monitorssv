package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/scrypt"
	"io"
)

func GenerateEncryptKey(data []byte) []byte {
	sk, err := scrypt.Key(data, nil, 32768, 8, 1, 32)
	if err != nil {
		panic(err)
	}

	sum := sha256.Sum256(append(data, sk...))
	return sum[:]
}

func Hash256(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

func EncryptData(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

func DecryptData(encryptedData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	data := make([]byte, len(encryptedData[aes.BlockSize:]))

	iv := encryptedData[:aes.BlockSize]
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(data, encryptedData[aes.BlockSize:])

	return data, nil
}
