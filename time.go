package gopkg

import (
	"fmt"
	"os"
	"time"
)

/**
 * @description: Timestamp2String 时间戳转字符串
 * @author:
 * @param timestampSec {int64 } 时间戳秒值
 * @param timestampNSec {int64} 时间戳纳秒值
 * @param format {string} 时间字符串格式化类型 如：2006/01/02 15:04:05
 * @param zone {*time.Location}
 * @return {string}
 */
func Timestamp2String(timestampSec, timestampNSec int64, format string, zone *time.Location) string {
	switch zone {
	default:
		return time.Unix(timestampSec, timestampNSec).Local().Format(format) //设置时间戳 使用模板格式化为日期字符串
	case time.UTC:
		return time.Unix(timestampSec, timestampNSec).UTC().Format(format) //设置时间戳 使用模板格式化为日期字符串
	}
}

/**
 * @description: String2Timestamp 字符串转时间戳
 * @author:
 * @param date {string} 待转换时间字符串 如：2019/09/17 10:16:56
 * @param format {string} 时间字符串格式化类型 如：2006/01/02 15:04:05
 * @param zone {*time.Location} zone 时区 如：time.Local / time.UTC
 * @return {int64,error}
 */
func String2Timestamp(date, format string, zone *time.Location) (int64, error) {
	var (
		theTime time.Time
		err     error
	)
	if theTime, err = time.ParseInLocation(format, date, zone); nil != err {
		return 0, err
	}
	return theTime.Unix(), nil
}

// GetDate 返回当前时间
func GetDate() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 03:04:05")
}

// GetRunTime 获取当前系统环境
func GetRunTime() string {
	//获取系统环境变量
	RUN_TIME := os.Getenv("RUN_TIME")
	if RUN_TIME == "" {
		fmt.Println("No RUN_TIME Can't start")
	}
	return RUN_TIME
}
