package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

const noneLen = 12

type Encrypter interface {
	Encrypt(text string) (string, error)
	Decrypt(text string) (string, error)
}

type encryption struct {
	key []byte
}

func NewEncryption(key []byte) Encrypter {
	return &encryption{key: key}
}

func (e *encryption) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, noneLen)
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(text), nil)
	ciphertext = append(ciphertext, nonce...)

	return hex.EncodeToString(ciphertext), nil
}

func (e *encryption) Decrypt(text string) (string, error) {
	if text == "" {
		return "", nil
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	plain, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}

	nonce := plain[len(plain)-noneLen:]
	plain = plain[:len(plain)-noneLen]

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, plain, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
