package gopkg

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestCryptoAES(t *testing.T) {
	data := []byte("Hello World")     // 待加密的数据
	key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
	t.Log("原文：", string(data))

	t.Log("------------------ CBC模式 --------------------")
	encrypted := AESEncryptCBC(data, key)
	t.Log("密文(hex)：", hex.EncodeToString(encrypted))
	t.Log("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := AESDecryptCBC(encrypted, key)
	t.Log("解密结果：", string(decrypted))

	t.Log("------------------ ECB模式 --------------------")
	encrypted = AESEncryptECB(data, key)
	t.Log("密文(hex)：", hex.EncodeToString(encrypted))
	t.Log("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = AESDecryptECB(encrypted, key)
	t.Log("解密结果：", string(decrypted))

	t.Log("------------------ CFB模式 --------------------")
	encrypted = AESEncryptCFB(data, key)
	t.Log("密文(hex)：", hex.EncodeToString(encrypted))
	t.Log("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = AESDecryptCFB(encrypted, key)
	t.Log("解密结果：", string(decrypted))

	t.Log("------------------ OFB模式 --------------------")
	encrypted = AESEncryptOFB(data, key)
	t.Log("密文(hex)：", hex.EncodeToString(encrypted))
	t.Log("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = AESDecryptOFB(encrypted, key)
	t.Log("解密结果：", string(decrypted))
}
