package sugar

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mattn/go-isatty"
	"github.com/mitchellh/go-homedir"
)

// 获取命令的路径
func FindCommandPath(str string) (string, error) {
	return exec.LookPath(str)

}

// HasExecutable in the system
//
// Usage:
// 	HasExecutable("bash")
func HasExecutable(binName string) bool {
	_, err := exec.LookPath(binName)
	return err == nil
}

// 获取标准正确输出
func RunCmd(str string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", str)
	result, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(result), nil

}

// 标准正确错误输出到标准正确输出
func RunCMD(str string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", str)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return string(result), err
	}
	return string(result), nil
}

// ExecCmd an command and return output.
// Usage:
// 	ExecCmd("ls", []string{"-al"})
func ExecCmd(binName string, args []string, workDir ...string) (string, error) {
	// create a new Cmd instance
	cmd := exec.Command(binName, args...)
	if len(workDir) > 0 {
		cmd.Dir = workDir[0]
	}

	bs, err := cmd.Output()
	return string(bs), err
}

func ShellExec(cmdLine string, shells ...string) (string, error) {
	// shell := "/bin/sh"
	shell := "bash"
	if len(shells) > 0 {
		shell = shells[0]
	}

	var out bytes.Buffer

	cmd := exec.Command(shell, "-c", cmdLine)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}

// // ExecLine quick exec an command line string
// func ExecLine(cmdLine string, workDir ...string) (string, error) {
// 	p := cmdline(cmdLine)

// 	// create a new Cmd instance
// 	cmd := p.NewExecCmd()
// 	if len(workDir) > 0 {
// 		cmd.Dir = workDir[0]
// 	}

// 	bs, err := cmd.Output()
// 	return string(bs), err
// }

// // QuickExec quick exec an simple command line
// func QuickExec(cmdLine string, workDir ...string) (string, error) {
// 	return ExecLine(cmdLine, workDir...)
// }

// MustFindUser must find an system user by name
func MustFindUser(uname string) *user.User {
	u, err := user.Lookup(uname)
	if err != nil {
		panic(err)
	}
	return u
}

// LoginUser must get current user
func LoginUser() *user.User {
	return CurrentUser()
}

// CurrentUser must get current user
func CurrentUser() *user.User {
	// check $HOME/.terminfo
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return u
}

// UserHomeDir is alias of os.UserHomeDir, but ignore error
func UserHomeDir() string {
	dir, _ := os.UserHomeDir()
	return dir
}

// UHomeDir get user home dir path.
func UHomeDir() string {
	// check $HOME/.terminfo
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return u.HomeDir
}

// HomeDir get user home dir path.
func HomeDir() string {
	dir, _ := homedir.Dir()
	return dir
}

// UserDir will prepend user home dir to subPath
func UserDir(subPath string) string {
	dir, _ := homedir.Dir()

	return dir + "/" + subPath
}

// UserCacheDir will prepend user `$HOME/.cache` to subPath
func UserCacheDir(subPath string) string {
	dir, _ := homedir.Dir()

	return dir + "/.cache/" + subPath
}

// UserConfigDir will prepend user `$HOME/.config` to subPath
func UserConfigDir(subPath string) string {
	dir, _ := homedir.Dir()

	return dir + "/.config/" + subPath
}

// Hostname is alias of os.Hostname, but ignore error
func Hostname() string {
	name, _ := os.Hostname()
	return name
}

// IsMSys msys(MINGW64) env，不一定支持颜色
func IsMSys() bool {
	// "MSYSTEM=MINGW64"
	if len(os.Getenv("MSYSTEM")) > 0 {
		return true
	}

	return false
}

// IsConsole check out is in stderr/stdout/stdin
//
// Usage:
// 	sysutil.IsConsole(os.Stdout)
func IsConsole(out io.Writer) bool {
	o, ok := out.(*os.File)
	if !ok {
		return false
	}

	fd := o.Fd()

	// fix: cannot use 'o == os.Stdout' to compare
	return fd == uintptr(syscall.Stdout) || fd == uintptr(syscall.Stdin) || fd == uintptr(syscall.Stderr)
}

// IsTerminal isatty check
//
// Usage:
// 	sysutil.IsTerminal(os.Stdout.Fd())
func IsTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd)
}

// StdIsTerminal os.Stdout is terminal
func StdIsTerminal() bool {
	return IsTerminal(os.Stdout.Fd())
}

var curShell string

// CurrentShell get current used shell env file.
// eg "/bin/zsh" "/bin/bash".
// if onlyName=true, will return "zsh", "bash"
func CurrentShell(onlyName bool) (path string) {
	var err error
	if curShell == "" {
		path, err = ShellExec("echo $SHELL")
		if err != nil {
			return ""
		}

		path = strings.TrimSpace(path)
		// cache result
		curShell = path
	} else {
		path = curShell
	}

	if onlyName && len(path) > 0 {
		path = filepath.Base(path)
	}
	return
}

// HasShellEnv has shell env check.
//
// Usage:
// 	HasShellEnv("sh")
// 	HasShellEnv("bash")
func HasShellEnv(shell string) bool {
	// can also use: "echo $0"
	out, err := ShellExec("echo OK", shell)
	if err != nil {
		return false
	}

	return strings.TrimSpace(out) == "OK"
}

// IsShellSpecialVar reports whether the character identifies a special
// shell variable such as $*.
func IsShellSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

// Workdir get
func Workdir() string {
	dir, _ := os.Getwd()
	return dir
}
