package util

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// 获取 ed25519 算法的签名字符串 base64编码
func GetED25519Base64SignStr(message []byte, privateKey ed25519.PrivateKey) (sgin string) {
	return base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, message))
}

// 获取 hamc sha256 算法签名的字符串 16进制编码
func GetHamcSha256HexEncodeSignStr(message, secretKey string) (sgin string) {
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return ""
	}
	sgin = hex.EncodeToString(mac.Sum(nil))
	return
}

func DecodePemPrivateKey(ed25519PrivateKey string) ed25519.PrivateKey {
	block, _ := pem.Decode([]byte(ed25519PrivateKey))
	if block == nil {
		panic("decode ed25519PrivateKey error")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(errors.New("ParsePKCS8PrivateKey err:" + err.Error()))
	}
	return priv.(ed25519.PrivateKey)
}

func DecodePrivateKey(ed25519PrivateKey string) (ed25519.PrivateKey, error) {
	keyStr := strings.ReplaceAll(ed25519PrivateKey, "\n", "")
	keyStr = strings.ReplaceAll(keyStr, " ", "")
	decoded, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, fmt.Errorf("ed25519 private key base64 decode failed: %w", err)
	}
	pri, err := x509.ParsePKCS8PrivateKey(decoded)
	if err != nil {
		return nil, fmt.Errorf("parse pkcs8 private key failed: %w", err)
	}
	ed25519Pri, ok := pri.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not ed25519 type")
	}
	return ed25519Pri, nil
}
