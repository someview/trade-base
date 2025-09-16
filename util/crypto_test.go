package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEncryptDecrypt 测试加密后再解密是否能正确还原原始数据
func TestEncryptDecrypt(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		key      string
		wantFail bool
	}{
		{
			name:     "基本示例",
			input:    "Hello, World!",
			key:      "test",
			wantFail: false,
		},
		{
			name:     "空字符串",
			input:    "",
			key:      "test",
			wantFail: false,
		},
		{
			name:     "中文字符",
			input:    "你好，世界！",
			key:      "test",
			wantFail: false,
		},
		{
			name:     "长文本",
			input:    strings.Repeat("Lorem ipsum dolor sit amet ", 100),
			key:      "test",
			wantFail: false,
		},
		{
			name:     "特殊字符",
			input:    "!@#$%^&*()_+<>?:\"{}|~`-=[]\\;',./",
			key:      "test",
			wantFail: false,
		},
		{
			name:     "二进制数据（私钥格式）",
			input:    "MC4CAQAwBQYDK2VwBCIEIFuewfB2QIm52nz9W2AR9gkrvvyzAapnz6YQBI0D6FrC",
			key:      "test",
			wantFail: false,
		},
		{
			name:     "32字节密钥",
			input:    "Hello, World!",
			key:      "12345678901234567890123456789012",
			wantFail: false,
		},
		{
			name:     "超长密钥（应该失败）",
			input:    "Hello, World!",
			key:      "1234567890123456789012345678901234567890",
			wantFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 加密
			encrypted, err := EncryptString(tc.input, tc.key)

			if tc.wantFail {
				assert.Error(t, err, "EncryptString() 期望失败但成功了")
				return
			}

			require.NoError(t, err, "EncryptString() 不应该返回错误")
			assert.Contains(t, encrypted, ":", "加密结果应包含分隔符")

			// 解密
			decrypted, err := DecryptString(encrypted, tc.key)
			require.NoError(t, err, "DecryptString() 不应该返回错误")

			// 验证解密后的结果是否与原始输入匹配
			assert.Equal(t, tc.input, decrypted, "解密结果应与原始输入匹配")
		})
	}
}

// TestPythonCompatibility 测试与Python实现的兼容性
func TestPythonCompatibility(t *testing.T) {
	// 以下是Python encrypt_string函数加密的示例
	// 注意：实际测试时应替换为使用Python代码生成的真实加密值
	pythonEncrypted := "aXuiYwUiMLwl+lnHR7ELNw==:ZOSJ8/UDo/M8cQbIgpFCKw=="
	originalText := "test message"
	key := "test"

	// 测试解密Python加密的数据
	decrypted, err := DecryptString(pythonEncrypted, key)
	require.NoError(t, err, "解密Python加密数据失败")
	assert.Equal(t, originalText, decrypted, "Python兼容性测试失败")

	// 测试Go加密后Python可以解密（需要手动验证）
	goEncrypted, err := EncryptString(originalText, key)
	require.NoError(t, err, "Go加密失败")
	t.Logf("Go加密结果 (用于Python解密验证): %s", goEncrypted)
}

// TestInvalidDecrypt 测试无效加密数据的解密
func TestInvalidDecrypt(t *testing.T) {
	testCases := []struct {
		name        string
		encrypted   string
		key         string
		expectedErr string
	}{
		{
			name:        "格式错误",
			encrypted:   "invalid_format_without_separator",
			key:         "test",
			expectedErr: "invalid encrypted data format",
		},
		{
			name:        "无效的IV Base64",
			encrypted:   "invalid!base64:YWJjZGVmZ2g=",
			key:         "test",
			expectedErr: "failed to decode IV",
		},
		{
			name:        "无效的密文 Base64",
			encrypted:   "YWJjZGVmZ2g=:invalid!base64",
			key:         "test",
			expectedErr: "failed to decode ciphertext",
		},
		{
			name:        "IV长度错误",
			encrypted:   "YWJj:YWJjZGVmZ2g=", // IV太短
			key:         "test",
			expectedErr: "invalid IV length",
		},
		{
			name:        "密文不是块大小的倍数",
			encrypted:   "YWJjZGVmZ2hpamtsbW5vcA==:YWJj", // 密文太短
			key:         "test",
			expectedErr: "ciphertext is not a multiple",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecryptString(tc.encrypted, tc.key)
			assert.Error(t, err, "应该返回错误")
			assert.Contains(t, err.Error(), tc.expectedErr, "错误消息应包含预期内容")
		})
	}
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	// 极短文本
	shortText := "a"
	key := "test"

	encrypted, err := EncryptString(shortText, key)
	require.NoError(t, err, "加密短文本失败")

	decrypted, err := DecryptString(encrypted, key)
	require.NoError(t, err, "解密短文本失败")
	assert.Equal(t, shortText, decrypted, "短文本解密结果不匹配")

	// 使用不同密钥解密
	wrongKey := "wrong"
	_, err = DecryptString(encrypted, wrongKey)
	assert.Error(t, err, "使用错误密钥应该解密失败")

	// 空密钥（应该成功，使用全0填充）
	emptyKey := ""
	encrypted, err = EncryptString(shortText, emptyKey)
	require.NoError(t, err, "使用空密钥加密失败")

	decrypted, err = DecryptString(encrypted, emptyKey)
	require.NoError(t, err, "使用空密钥解密失败")
	assert.Equal(t, shortText, decrypted, "使用空密钥解密结果不匹配")
}

// func TestEncryptAndDecryptString(t *testing.T) {
// 	val, err := EncryptString(testdata.GateAPIKey, "betago")
// 	assert.Nil(t, err)
// 	apiKey, err := DecryptString(val, "betago")
// 	assert.Nil(t, err)
// 	assert.Equal(t, testdata.GateAPIKey, apiKey)
// 	priveteVal, err := EncryptString(testdata.GatePrivateKey, "betago")
// 	assert.Nil(t, err)
// 	privateKey, err := DecryptString(priveteVal, "betago")
// 	assert.Nil(t, err)
// 	assert.Equal(t, testdata.GatePrivateKey, privateKey)

// 	cryptoPriKey := "/pjqmFeBkTZxGYywGplVAQ==:BZ1+83487uFgz22qdo6lC69hUI/X0GEAtedO4QiIwB9zUlq/RX7hqpcRb1kro/gZRwOrmCyUJs1mp1H5b8DkYrmcrbpjk5Uhsr5wA4SwvwE="
// 	realPriKey := "MC4CAQAwBQYDK2VwBCIEINiPj07y/DBWgSgxCH8+1m7fZjUtdKv6WHpyiQdkAtDb"
// 	realTestKey, err := DecryptString(cryptoPriKey, "betago")
// 	assert.Nil(t, err)
// 	assert.Equal(t, realTestKey, realPriKey)
// }
