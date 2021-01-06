package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// Encrypting performs the encryption from plain text to encrypted text
func Encrypting(stringToEncrypt string, keyString string) (encryptedString string, err error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	plaintext := []byte(stringToEncrypt)
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	encryptedString = fmt.Sprintf("%x", ciphertext)
	return
}

// Decrypting performs the decryption from encrypted text to plain text
func Decrypting(encryptedString string, keyString string) (decryptedString string, err error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return
	}

	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return
	}

	decryptedString = fmt.Sprintf("%s", plaintext)
	return
}
