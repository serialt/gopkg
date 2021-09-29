package gopkg

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
	Md5File

	strFilename: 计算md5值的文件名

*/

// Md5File 计算文件的md5值
func Md5File(strFilename string) string {
	f, err := os.Open(strFilename)
	if err != nil {
		return ""
	}

	defer f.Close()

	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, f); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}

// Md5Sum 计算md5值
func Md5Sum(data []byte) string {
	md5Hash := md5.New()
	md5Hash.Write(data)
	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}

// Md5SumUpper 计算md5值的大写字母
func Md5SumUpper(data []byte) string {
	return strings.ToUpper(Md5Sum(data))
}

// Md5String 加密字符串
func Md5String(data string) string {
	return Md5Sum([]byte(data))
}

// Md5StringUpper 加密字符串变大写
func Md5StringUpper(data string) string {
	return strings.ToUpper(Md5String(data))
}
