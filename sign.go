package gopkg

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

// GetHmacSHA256Sign 获取hmac sha256签名
func GetHmacSHA256Sign(secret, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// GetHmacSHA512Sign 获取hmac sha512签名
func GetHmacSHA512Sign(secret, params string) (string, error) {
	mac := hmac.New(sha512.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// GetHmacSHA1Sign 获取hmac sha1签名
func GetHmacSHA1Sign(secret, params string) (string, error) {
	mac := hmac.New(sha1.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// GetHmacMD5Sign 获取hmac md5签名
func GetHmacMD5Sign(secret, params string) (string, error) {
	mac := hmac.New(md5.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// GetHmacSha384Sign 获取hmac sha384签名
func GetHmacSha384Sign(secret, params string) (string, error) {
	mac := hmac.New(sha512.New384, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", nil
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// GetHmacSHA256Base64Sign 获取hmac sha256base64签名
func GetHmacSHA256Base64Sign(secret, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	signByte := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signByte), nil
}

// GetHmacSHA512Base64Sign 获取hmac sha512base64签名
func GetHmacSHA512Base64Sign(hmac_key string, hmac_data string) string {
	hmh := hmac.New(sha512.New, []byte(hmac_key))
	hmh.Write([]byte(hmac_data))

	hex_data := hex.EncodeToString(hmh.Sum(nil))
	hash_hmac_bytes := []byte(hex_data)
	hmh.Reset()

	return base64.StdEncoding.EncodeToString(hash_hmac_bytes)
}
