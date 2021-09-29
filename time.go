package gopkg

import (
	"fmt"
	"os"
	"time"
)

// Timestamp2String 时间戳转字符串
//
// timestampSec 时间戳秒值
//
// timestampNSec 时间戳纳秒值
//
// format 时间字符串格式化类型 如：2006/01/02 15:04:05
func Timestamp2String(timestampSec, timestampNSec int64, format string, zone *time.Location) string {
	switch zone {
	default:
		return time.Unix(timestampSec, timestampNSec).Local().Format(format) //设置时间戳 使用模板格式化为日期字符串
	case time.UTC:
		return time.Unix(timestampSec, timestampNSec).UTC().Format(format) //设置时间戳 使用模板格式化为日期字符串
	}
}

// 返回当前时间
func GetDate() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 03:04:05")
}

// 获取当前系统环境
func GetRunTime() string {
	//获取系统环境变量
	RUN_TIME := os.Getenv("RUN_TIME")
	if RUN_TIME == "" {
		fmt.Println("No RUN_TIME Can't start")
	}
	return RUN_TIME
}
