package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// KeySize is the size of the encryption key in bytes (AES-256)
	KeySize = 32
	// NonceSize is the size of the nonce for GCM (12 bytes recommended)
	NonceSize = 12
	// SaltSize is the size of the salt for key derivation
	SaltSize = 16
)

// Crypto handles encryption and decryption of VPN packets
type Crypto struct {
	aead cipher.AEAD
}

// NewCrypto creates a new crypto instance from a password
func NewCrypto(password string) (*Crypto, error) {
	// Derive key from password using PBKDF2
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(password), salt, 4096, KeySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Crypto{aead: aead}, nil
}

// NewCryptoFromKey creates a new crypto instance from a raw key
func NewCryptoFromKey(key []byte) (*Crypto, error) {
	if len(key) != KeySize {
		return nil, errors.New("invalid key size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Crypto{aead: aead}, nil
}

// Encrypt encrypts plaintext and returns ciphertext with nonce prepended
func (c *Crypto) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := c.aead.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext (with nonce prepended) and returns plaintext
func (c *Crypto) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < NonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:NonceSize], ciphertext[NonceSize:]
	plaintext, err := c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Overhead returns the encryption overhead (nonce + GCM tag)
func (c *Crypto) Overhead() int {
	return NonceSize + c.aead.Overhead()
}
