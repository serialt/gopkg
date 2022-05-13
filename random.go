package gopkg

import (
	cryRand "crypto/rand"
	"math/big"
	"math/rand"
	"time"
)

// RandRangeInt 获取min和max之前的随机数
func RandRangeInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

// RandInt64Crypto 通过crypto库生成int64随机数
func RandInt64Crypto() int64 {
	n, _ := cryRand.Int(cryRand.Reader, big.NewInt(100))
	return n.Int64()
}

// RandInt64 生成int64随机数
func RandInt64(min int64, max int64) int64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Int63n(max-min)
}

// GetRandomString 获取随机字符串
func GetRandomString(l int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// GetRandomNumeral 获取随机自然数字符串
func GetRandomNumeral(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
