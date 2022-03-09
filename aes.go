package sugar

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// AESCommon AES工具
//
// 对称加密的四种模式(ECB、CBC、CFB、OFB、GCM)
//
// ECB模式——电码本模式（Electronic Codebook Book (ECB)）
//
// 优点:1.简单;2.有利于并行计算;3.误差不会被传送；
//
// 缺点:1.不能隐藏明文的模式;2.可能对明文进行主动攻击；
//
// ======================================
//
// CBC模式——密码分组链接模式（Cipher Block Chaining (CBC)）
//
// 优点:1.不容易主动攻击,安全性好于ECB,适合传输长度长的报文,是SSL、IPSec的标准。
//
// 缺点:1.不利于并行计算;2.误差传递;3.需要初始化向量IV
//
// ======================================
//
// CFB模式——密码反馈模式（Cipher FeedBack (CFB)）
//
// 优点:1.隐藏了明文模式;2.分组密码转化为流模式;3.可以及时加密传送小于分组的数据;
//
// 缺点:1.不利于并行计算;2.误差传送:一个明文单元损坏影响多个单元;3.唯一的IV;
//
// ======================================
//
// OFB模式——输出反馈模式（Output FeedBack (OFB)）
//
// 优点:1.隐藏了明文模式;2.分组密码转化为流模式;3.可以及时加密传送小于分组的数据;
//
// 缺点:1.不利于并行计算;2.对明文的主动攻击是可能的;3.误差传送：一个明文单元损坏影响多个单元;
//
// ======================================
// GCM加密
//
// GCM中的G就是指GMAC，C就是指CTR。
// GCM可以提供对消息的加密和完整性校验，另外，它还可以提供附加消息的完整性校验。在实际应用场景中，
// 有些信息是我们不需要保密，但信息的接收者需要确认它的真实性的，例如源IP，源端口，目的IP，IV，
// 等等。因此，我们可以将这一部分作为附加消息加入到MAC值的计算当中。下图的Ek表示用对称秘钥k对输入
// 做AES运算。最后，密文接收者会收到密文、IV（计数器CTR的初始值）
//

//--------------------------------------------------------------------------------------------------------------------

// AESEncryptCBC AES CBC模式加密
func AESEncryptCBC(data []byte, key []byte) (encrypted []byte) {
	// 分组密钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	blockMode, blockSize := aesBlockMode(key, true) // 加密模式
	data = aesPkcs5Padding(data, blockSize)         // 补全码
	encrypted = make([]byte, len(data))             // 创建数组
	blockMode.CryptBlocks(encrypted, data)          // 加密
	return encrypted
}

// AESDecryptCBC AES CBC模式解密
func AESDecryptCBC(encrypted []byte, key []byte) []byte {
	blockMode, _ := aesBlockMode(key, false)    // 加密模式
	decrypted := make([]byte, len(encrypted))   // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted) // 解密
	decrypted = aesPkcs5UnPadding(decrypted)    // 去除补全码
	return decrypted
}

// aesBlockMode 加密模式
func aesBlockMode(key []byte, encrypt bool) (cipher.BlockMode, int) {
	block, _ := aes.NewCipher(key) // 分组密钥
	blockSize := block.BlockSize() // 获取密钥块的长度
	if encrypt {
		return cipher.NewCBCEncrypter(block, key[:blockSize]), blockSize // 加密模式
	}
	return cipher.NewCBCDecrypter(block, key[:blockSize]), blockSize // 加密模式
}

// aesPkcs5Padding 明文补码算法
func aesPkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// aesPkcs5UnPadding 明文减码算法
func aesPkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

//--------------------------------------------------------------------------------------------------------------------

// AESEncryptECB AES ECB模式加密
func AESEncryptECB(data []byte, key []byte) (encrypted []byte) {
	cpr, _ := aes.NewCipher(aesGenerateKey(key))
	length := (len(data) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, data)
	pad := byte(len(plain) - len(data))
	for i := len(data); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cpr.BlockSize(); bs <= len(data); bs, be = bs+cpr.BlockSize(), be+cpr.BlockSize() {
		cpr.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted
}

// AESDecryptECB AES ECB模式解密
func AESDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cpr, _ := aes.NewCipher(aesGenerateKey(key))
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, cpr.BlockSize(); bs < len(encrypted); bs, be = bs+cpr.BlockSize(), be+cpr.BlockSize() {
		cpr.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim]
}

func aesGenerateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

//--------------------------------------------------------------------------------------------------------------------

// AESEncryptCFB AES CFB模式加密
func AESEncryptCFB(data []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(data))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], data)
	return encrypted
}

// AESDecryptCFB AES CFB模式解密
func AESDecryptCFB(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("cipher text too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

//--------------------------------------------------------------------------------------------------------------------

// AESEncryptOFB AES OFB模式加密
func AESEncryptOFB(data []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(data))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], data)
	return encrypted
}

// AESDecryptOFB AES OFB模式解密
func AESDecryptOFB(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

//--------------------------------------------------------------------------------------------------------------------

// AssertKey  判断AES密钥长度是否合法,密钥长度16/24/32字节
func AssertKey(key []byte) {
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		panic("key length must be 16/24/32 bytes")
	}
}

// AssertIV 判断AES向量长度是否合法
func AssertIV(iv []byte) {
	keyLen := len(iv)
	if keyLen < 16 {
		panic("iv length must >= 16 bytes")
	}
}

// encodeBase64 base64编码
func encodeBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// decodeBase64 base64解码
func decodeBase64(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

// AESEncryptGCM AES GCM模式加密
func AESEncryptGCM(plainText []byte, cipherKey []byte) (result []byte, err error) {
	var (
		aesBlock cipher.Block
		gcm      cipher.AEAD
	)
	aesBlock, err = aes.NewCipher(cipherKey)
	if err != nil {
		return
	}
	gcm, err = cipher.NewGCM(aesBlock)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	result = gcm.Seal(nonce, nonce, plainText, nil)
	return
}

// AESDecryptGCM AESDecryptGCM AESDecryptGCM模式解密
func AESDecryptGCM(encryptText []byte, cipherKey []byte) (result []byte, err error) {
	var (
		aesBlock cipher.Block
		gcm      cipher.AEAD
	)
	aesBlock, err = aes.NewCipher([]byte(cipherKey))
	if err != nil {
		return
	}
	gcm, err = cipher.NewGCM(aesBlock)
	if err != nil {
		return
	}
	nonceSize := gcm.NonceSize()
	if len(encryptText) < nonceSize {
		err = errors.New("cipher: incorrect size given to GCM")
		return
	}
	nonce, cipherText := encryptText[:nonceSize], encryptText[nonceSize:]
	return gcm.Open(nil, nonce, cipherText, nil)
}

// Encrypt gcm-aes base64加密字符串 秘钥长度为16/24/32位
func AESEncryptGCMBase64String(src, cipherKey string) (dst string, err error) {
	var encryptResult []byte
	encryptResult, err = AESEncryptGCM([]byte(src), []byte(cipherKey))
	if err != nil {
		return
	}
	dst = encodeBase64(encryptResult)
	return
}

// Decrypt gcm-aes base64解密为字符串 秘钥长度为16/24/32位
func AESDecryptGCMBase64String(src string, cipherKey string) (dst string, err error) {
	var (
		encryptResult []byte
		decryptResult []byte
	)
	encryptResult, err = decodeBase64(src)
	if err != nil {
		return
	}
	decryptResult, err = AESDecryptGCM(encryptResult, []byte(cipherKey))
	if err != nil {
		return
	}
	dst = string(decryptResult)
	return
}
