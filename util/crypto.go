package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

// DecryptString 解密字符串
func DecryptString(encryptedData string, key string) (string, error) {
	// 分割IV和密文
	parts := strings.Split(encryptedData, ":")
	if len(parts) != 2 {
		return "", errors.New("invalid encrypted data format")
	}

	// 解码base64的IV和密文
	iv, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("failed to decode IV: %v", err)
	}

	if len(iv) != aes.BlockSize {
		return "", fmt.Errorf("invalid IV length: got %d, want %d", len(iv), aes.BlockSize)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %v", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	// 使用与Python相同的密钥处理方式（用ASCII '0' 填充到32字节）
	key_bytes := []byte(key)
	if len(key_bytes) > 32 {
		return "", fmt.Errorf("key is too long: got %d bytes, maximum is 32", len(key_bytes))
	}

	// 用ASCII字符'0'填充到32字节，与Python的ljust(32, b'0')相同
	padded_key := make([]byte, 32)
	copy(padded_key, key_bytes)
	for i := len(key_bytes); i < 32; i++ {
		padded_key[i] = '0' // ASCII字符'0'，值为48
	}

	// 创建cipher
	block, err := aes.NewCipher(padded_key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// 创建CBC解密器
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 检查解密后的数据长度
	if len(plaintext) == 0 {
		return "", errors.New("decrypted data is empty")
	}

	// 获取填充值
	padding := int(plaintext[len(plaintext)-1])

	// 验证填充值是否合理
	if padding <= 0 || padding > aes.BlockSize {
		return "", fmt.Errorf("invalid padding size: %d (should be 1-%d)", padding, aes.BlockSize)
	}

	if padding > len(plaintext) {
		return "", fmt.Errorf("padding size %d exceeds data length %d", padding, len(plaintext))
	}

	// 计算原始数据的结束位置
	originalLength := len(plaintext) - padding
	if originalLength < 0 {
		return "", fmt.Errorf("invalid padding: calculated length is negative (%d)", originalLength)
	}

	// 验证所有填充字节
	for i := originalLength; i < len(plaintext); i++ {
		if plaintext[i] != byte(padding) {
			return "", fmt.Errorf("invalid padding value at position %d: got %d, expected %d", i, plaintext[i], padding)
		}
	}

	// 移除填充
	return string(plaintext[:originalLength]), nil
}

// EncryptString 加密字符串
// 使用AES-CBC模式加密，与Python实现兼容
func EncryptString(plaintext string, key string) (string, error) {
	// 处理密钥（将密钥填充到32字节）
	key_bytes := []byte(key)
	if len(key_bytes) > 32 {
		return "", fmt.Errorf("key is too long: got %d bytes, maximum is 32", len(key_bytes))
	}

	// 用ASCII字符'0'填充密钥到32字节，与Python的ljust(32, b'0')相同
	padded_key := make([]byte, 32)
	copy(padded_key, key_bytes)
	for i := len(key_bytes); i < 32; i++ {
		padded_key[i] = '0'
	}

	// 创建AES加密块
	block, err := aes.NewCipher(padded_key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// 对明文进行PKCS#7填充，与Python的pad()函数行为相同
	plaintext_bytes := []byte(plaintext)
	padding := aes.BlockSize - (len(plaintext_bytes) % aes.BlockSize)
	padded_plaintext := make([]byte, len(plaintext_bytes)+padding)
	copy(padded_plaintext, plaintext_bytes)

	// 填充值等于填充的字节数（PKCS#7标准）
	for i := len(plaintext_bytes); i < len(padded_plaintext); i++ {
		padded_plaintext[i] = byte(padding)
	}

	// 生成随机初始化向量(IV)，与Python中AES.new()自动生成IV相同
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %v", err)
	}

	// 加密填充后的明文
	ciphertext := make([]byte, len(padded_plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded_plaintext)

	// 将IV和密文分别编码为base64并用冒号连接，与Python实现的格式相同
	iv_base64 := base64.StdEncoding.EncodeToString(iv)
	ciphertext_base64 := base64.StdEncoding.EncodeToString(ciphertext)

	return iv_base64 + ":" + ciphertext_base64, nil
}
