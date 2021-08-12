package go-pkg

import (
	"os/exec"
)

// 获取命令的路径
func FindCommandPath(str string) (string, error) {
	path, err := exec.LookPath(str)
	if err != nil {
		return "", err
	}
	return path, nil

}

// 获取标准正确输出
func RunCmd(str string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", str)
	result, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(result), nil

}

// 标准正确错误输出到标准正确输出
func RunCMD(str string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", str)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return string(result), err
	}
	return string(result), nil
}
