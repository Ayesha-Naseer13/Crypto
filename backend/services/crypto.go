package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
	"os"
)

type CryptoService struct {
	aesKey []byte
}

func NewCryptoService() *CryptoService {
	key := os.Getenv("AES_KEY")
	if len(key) < 32 {
		key = "default-32-byte-key-for-aes-enc!"
	}
	return &CryptoService{aesKey: []byte(key[:32])}
}

// GenerateKeyPair generates ECDSA key pair
func (s *CryptoService) GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// PublicKeyToString converts public key to hex string
func (s *CryptoService) PublicKeyToString(pub *ecdsa.PublicKey) string {
	return hex.EncodeToString(elliptic.Marshal(pub.Curve, pub.X, pub.Y))
}

// PrivateKeyToString converts private key to hex string
func (s *CryptoService) PrivateKeyToString(priv *ecdsa.PrivateKey) string {
	return hex.EncodeToString(priv.D.Bytes())
}

// StringToPublicKey converts hex string to public key
func (s *CryptoService) StringToPublicKey(pubStr string) (*ecdsa.PublicKey, error) {
	pubBytes, err := hex.DecodeString(pubStr)
	if err != nil {
		return nil, err
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), pubBytes)
	if x == nil {
		return nil, errors.New("invalid public key")
	}
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}

// StringToPrivateKey converts hex string to private key
func (s *CryptoService) StringToPrivateKey(privStr string, pubKey *ecdsa.PublicKey) (*ecdsa.PrivateKey, error) {
	privBytes, err := hex.DecodeString(privStr)
	if err != nil {
		return nil, err
	}
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey = *pubKey
	priv.D = new(big.Int).SetBytes(privBytes)
	return priv, nil
}

// GenerateWalletID creates wallet ID from public key hash
func (s *CryptoService) GenerateWalletID(publicKey string) string {
	hash := sha256.Sum256([]byte(publicKey))
	return hex.EncodeToString(hash[:])[:40] // 40 char wallet ID
}

// SignData signs data with private key
func (s *CryptoService) SignData(privateKey *ecdsa.PrivateKey, data string) (string, error) {
	hash := sha256.Sum256([]byte(data))
	r, ss, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}
	signature := append(r.Bytes(), ss.Bytes()...)
	return hex.EncodeToString(signature), nil
}

// VerifySignature verifies signature with public key
func (s *CryptoService) VerifySignature(publicKey *ecdsa.PublicKey, data string, signatureHex string) bool {
	hash := sha256.Sum256([]byte(data))
	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil || len(sigBytes) != 64 {
		return false
	}
	r := new(big.Int).SetBytes(sigBytes[:32])
	ss := new(big.Int).SetBytes(sigBytes[32:])
	return ecdsa.Verify(publicKey, hash[:], r, ss)
}

// EncryptPrivateKey encrypts private key with AES
func (s *CryptoService) EncryptPrivateKey(privateKey string) (string, error) {
	block, err := aes.NewCipher(s.aesKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(privateKey), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPrivateKey decrypts private key with AES
func (s *CryptoService) DecryptPrivateKey(encrypted string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.aesKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// HashSHA256 computes SHA-256 hash
func HashSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
