package gopkg

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/go-homedir"
)

var (
	// perm and flags for create log file
	DefaultDirPerm  os.FileMode = 0775
	DefaultFilePerm os.FileMode = 0665

	DefaultFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

// alias methods
var (
	DirExist  = IsDir
	FileExist = IsFile
)

// 当前项目根目录
var API_ROOT string

// 获取项目路径
func GetRootPath() string {

	if API_ROOT != "" {
		return API_ROOT
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		print(err.Error())
	}

	API_ROOT = strings.Replace(dir, "\\", "/", -1)
	return API_ROOT
}

// 判断文件目录否存在
func IsDirExists(path string) bool {
	if path == "" {
		return false
	}
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

}

// 判断文件目录否存在
func IsDir(path string) bool {
	if path == "" {
		return false
	}
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

}

// Mode get file mode
func Mode(path string) os.FileMode {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0755
	}
	return fileInfo.Mode()
}

// CreateDir 创建目录
func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist := IsDirExists(v)

		if !exist {

			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return err
}

// FileExt get filename ext. alias of path.Ext()
func FileExt(fpath string) string {
	return path.Ext(fpath)
}

// Suffix get filename ext. alias of path.Ext()
func Suffix(fpath string) string {
	return path.Ext(fpath)
}

// 创建文件
func MkdirFile(path string) error {

	err := os.Mkdir(path, os.ModePerm) //在当前目录下生成md目录
	if err != nil {
		return err
	}
	return nil
}

// FileExists reports whether the named file or directory exists.
func FileExists(path string) bool {
	return IsFile(path)
}

// IsFile reports whether the named file or directory exists.
func IsFile(path string) bool {
	if path == "" {
		return false
	}

	if fi, err := os.Stat(path); err == nil {
		return !fi.IsDir()
	}
	return false
}

// IsAbsPath is abs path.
func IsAbsPath(aPath string) bool {
	return path.IsAbs(aPath)
}

// DisableCache will disable caching of the home directory. Caching is enabled
// by default.
var DisableCache bool

var homedirCache string
var cacheLock sync.RWMutex

// Dir returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func HomeDir() (string, error) {
	if !DisableCache {
		cacheLock.RLock()
		cached := homedirCache
		cacheLock.RUnlock()
		if cached != "" {
			return cached, nil
		}
	}

	cacheLock.Lock()
	defer cacheLock.Unlock()

	var result string
	var err error
	if runtime.GOOS == "windows" {
		result, err = dirWindows()
	} else {
		// Unix-like system, so just assume Unix
		result, err = dirUnix()
	}

	if err != nil {
		return "", err
	}
	homedirCache = result
	return result, nil
}

// Expand expands the path to include the home directory if the path
// is prefixed with `~`. If it isn't prefixed with `~`, the path is
// returned as-is.
func Expand(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := HomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}

func dirUnix() (string, error) {
	homeEnv := "HOME"
	if runtime.GOOS == "plan9" {
		// On plan9, env vars are lowercase.
		homeEnv = "home"
	}

	// First prefer the HOME environmental variable
	if home := os.Getenv(homeEnv); home != "" {
		return home, nil
	}

	var stdout bytes.Buffer

	// If that fails, try OS specific commands
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("sh", "-c", `dscl -q . -read /Users/"$(whoami)" NFSHomeDirectory | sed 's/^[^ ]*: //'`)
		cmd.Stdout = &stdout
		if err := cmd.Run(); err == nil {
			result := strings.TrimSpace(stdout.String())
			if result != "" {
				return result, nil
			}
		}
	} else {
		cmd := exec.Command("getent", "passwd", strconv.Itoa(os.Getuid()))
		cmd.Stdout = &stdout
		if err := cmd.Run(); err != nil {
			// If the error is ErrNotFound, we ignore it. Otherwise, return it.
			if err != exec.ErrNotFound {
				return "", err
			}
		} else {
			if passwd := strings.TrimSpace(stdout.String()); passwd != "" {
				// username:password:uid:gid:gecos:home:shell
				passwdParts := strings.SplitN(passwd, ":", 7)
				if len(passwdParts) > 5 {
					return passwdParts[5], nil
				}
			}
		}
	}

	// If all else fails, try the shell
	stdout.Reset()
	cmd := exec.Command("sh", "-c", "cd && pwd")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func dirWindows() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

// Dir get dir path, without last name. 获取路径的目录
func PathDir(fpath string) string {
	return filepath.Dir(fpath)
}

// Name get file/dir name 获取路径的文件名
func Name(fpath string) string {
	// return path.Base(fpath)
	return filepath.Base(fpath)
}

// DeleteFile 删除文件或目录
func DeleteFile(filePath string) error {
	return os.RemoveAll(filePath)
}

//创建文件夹,支持x/a/a  多层级
func MkDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			//文件夹不存在，创建
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// FilePathExists 判断路径是否存在
func FilePathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FileReadFirstLine 从文件中读取第一行并返回字符串数组
func FileReadFirstLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()
	finReader := bufio.NewReader(file)
	inputString, _ := finReader.ReadString('\n')
	return StringTrimN(inputString), nil
}

// FileReadPointLine 从文件中读取指定行并返回字符串数组
func FileReadPointLine(filePath string, line int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()
	finReader := bufio.NewReader(file)
	lineCount := 1
	for {
		inputString, err := finReader.ReadString('\n')
		//fmt.Println(inputString)
		if err == io.EOF {
			if lineCount == line {
				return inputString, nil
			}
			return "", errors.New("index out of line count")
		}
		if lineCount == line {
			return inputString, nil
		}
		lineCount++
	}
}

// FileReadLines 从文件中逐行读取并返回字符串数组
func FileReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	finReader := bufio.NewReader(file)
	var fileList []string
	for {
		inputString, err := finReader.ReadString('\n')
		//fmt.Println(inputString)
		if err == io.EOF {
			fileList = append(fileList, StringTrimN(inputString))
			break
		}
		fileList = append(fileList, StringTrimN(inputString))
	}
	//fmt.Println("fileList",fileList)
	return fileList, nil
}

// FileParentPath 文件父路径
func FileParentPath(filePath string) string {
	return filePath[0:strings.LastIndex(filePath, "/")]
}

// WriteStringToFile write string to file
func WriteStringToFile(content, path string, mode os.FileMode) (err error) {
	bytes := []byte(content)
	return ioutil.WriteFile(path, bytes, mode)
}

// FileAppend 追加内容到文件中
//
// filePath 文件地址
//
// data 内容
//
// force 如果文件已存在，会将文件清空
//
// It returns the number of bytes written and an error
func FileAppend(filePath string, data []byte, force bool) (int, error) {
	var (
		file *os.File
		n    int
		err  error
	)
	exist := FilePathExists(filePath)
	if exist {
		if force {
			// 创建文件，如果文件已存在，会将文件清空
			if file, err = os.Create(filePath); err != nil {
				return 0, err
			}
		} else {
			if file, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644); nil != err {
				return 0, err
			}
		}
	} else {
		parentPath := FileParentPath(filePath)
		if err = os.MkdirAll(parentPath, os.ModePerm); nil != err {
			return 0, err
		}
		if file, err = os.Create(filePath); err != nil {
			return 0, err
		}
	}
	defer func() {
		_ = file.Close()
	}()
	// 将数据写入文件中
	//file.WriteString(string(data)) //写入字符串
	if n, err = file.Write(data); nil != err { // 写入byte的slice数据
		return 0, err
	}
	return n, nil
}

// Filter file filter
type Filter func(os.FileInfo) bool

// GetDirListWithFilter get directory list with filter
func GetDirListWithFilter(path string, filter Filter) ([]string, error) {
	var dirList []string

	paths, err := filepath.Glob(filepath.Join(path, "*"))

	log.Printf("paths: %v", paths)

	for _, value := range paths {
		f, err := os.Stat(value)
		if err != nil {
			return dirList, err
		}
		if filter != nil && !filter(f) {
			continue
		}
		if f.IsDir() {
			dir := strings.Replace(value, path, "", 1)
			if strings.HasPrefix(dir, "/") {
				dir = strings.Replace(dir, "/", "", 1)
			}
			dirList = append(dirList, dir)
		}
	}

	return dirList, err
}

// RecreateDir recreate dir
func RecreateDir(dir string) error {
	mode := Mode(dir)
	_ = os.RemoveAll(dir)
	return os.MkdirAll(dir, mode)
}

// GetFilepaths get all filepaths in a directory tree
func GetFilepaths(dir string) ([]string, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}

// File template file
type File struct {
	Path    string
	Content string
}

// GetFiles get files
func GetFiles(dir string) ([]*File, error) {
	var files []*File

	paths, err := GetFilepaths(dir)
	if err != nil {
		return files, err
	}

	for _, path := range paths {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return files, err
		}
		content := string(bytes)
		file := &File{
			Path:    path,
			Content: content,
		}
		files = append(files, file)
	}

	return files, nil
}

// FileModify 修改文件中指定位置的内容
//
// filePath 文件地址
//
// offset 以0为起始坐标的偏移量
//
// data 内容
//
// force 如果文件已存在，会将文件清空
//
// It returns the number of bytes written and an error
func FileModify(filePath string, offset int64, data []byte, force bool) (int, error) {
	var (
		file *os.File
		n    int
		err  error
	)
	exist := FilePathExists(filePath)
	if exist {
		if force {
			// 创建文件，如果文件已存在，会将文件清空
			if file, err = os.Create(filePath); err != nil {
				return 0, err
			}
		} else {
			if file, err = os.OpenFile(filePath, os.O_RDWR, 0644); nil != err {
				return 0, err
			}
		}
	} else {
		parentPath := FileParentPath(filePath)
		if err = os.MkdirAll(parentPath, os.ModePerm); nil != err {
			return 0, err
		}
		if file, err = os.Create(filePath); err != nil {
			return 0, err
		}
	}
	defer func() {
		_ = file.Close()
	}()
	// 以0为起始坐标偏移指标到指定位置
	if _, err = file.Seek(offset, io.SeekStart); nil != err {
		return 0, err
	}
	// 将数据写入文件中
	//file.WriteString(string(data)) //写入字符串
	if n, err = file.Write(data); nil != err { // 写入byte的slice数据
		return 0, err
	}
	return n, nil
}

// FileLoopDirs 遍历目录下的所有子目录，即返回pathname下面的所有目录，目录为绝对路径
func FileLoopDirs(pathname string) ([]string, error) {
	var s []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

// FileLoopOneDirs 遍历目录下的所有子目录，即返回pathname下面的所有目录，目录为相对路径
func FileLoopOneDirs(pathname string) ([]string, error) {
	var s []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			s = append(s, fi.Name())
		}
	}
	return s, nil
}

// FileLoopFiles 遍历文件夹及子文件夹下的所有文件，即返回pathname目录下所有的文件，文件名为绝对路径
func FileLoopFiles(pathname string) ([]string, error) {
	var s []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := path.Join(pathname, fi.Name())
			sNew, err := FileLoopFiles(fullDir)
			if err != nil {
				return s, err
			}
			s = append(s, sNew...)
		} else {
			fullName := filepath.Join(pathname, fi.Name())
			s = append(s, fullName)
		}
	}
	return s, nil
}

// FileLoopFileNames 遍历文件夹及子文件夹下的所有文件名，即返回pathname目录下所有的文件，文件名为相对路径
func FileLoopFileNames(pathname string) ([]string, error) {
	var s []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := path.Join(pathname, fi.Name())
			sNew, err := FileLoopFileNames(fullDir)
			if err != nil {
				return s, err
			}
			s = append(s, sNew...)
		} else {
			s = append(s, fi.Name())
		}
	}
	return s, nil
}

// FileMove 移动文件
func FileMove(src string, dst string) (err error) {
	if dst == "" {
		return nil
	}
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}
	var revoke = false
	dir := filepath.Dir(dst)
Redirect:
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
		if !revoke {
			revoke = true
			goto Redirect
		}
	}
	return os.Rename(src, dst)
}

// TrimSpace 去除空格
func TrimSpace(target interface{}) {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		return
	}
	t = t.Elem()
	v := reflect.ValueOf(target).Elem()
	for i := 0; i < t.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(strings.TrimSpace(v.Field(i).String()))
		}
	}
	return
}

// FileCompressZip 压缩文件
// files 文件数组，可以是不同dir下的文件或者文件夹
// dest 压缩文件存放地址
func FileCompressZip(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer func() { _ = d.Close() }()
	w := zip.NewWriter(d)
	defer func() { _ = w.Close() }()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

// FileCompressTar 压缩文件
// files 文件数组，可以是不同dir下的文件或者文件夹
// dest 压缩文件存放地址
func FileCompressTar(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer func() { _ = d.Close() }()
	w := tar.NewWriter(d)
	defer func() { _ = w.Close() }()
	for _, file := range files {
		err := compressTar(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	var (
		info   os.FileInfo
		header *zip.FileHeader
		writer io.Writer
		err    error
	)
	defer func() { _ = file.Close() }()
	info, err = file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			fil, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(fil, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		if header, err = zip.FileInfoHeader(info); nil != err {
			return err
		}
		header.Name = prefix + "/" + header.Name
		if writer, err = zw.CreateHeader(header); nil != err {
			return err
		}
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func compressTar(file *os.File, prefix string, tw *tar.Writer) error {
	var (
		info   os.FileInfo
		header *tar.Header
		err    error
	)
	defer func() { _ = file.Close() }()
	info, err = file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			fil, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compressTar(fil, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		if header, err = tar.FileInfoHeader(info, ""); nil != err {
			return err
		}
		header.Name = prefix + "/" + header.Name
		if err = tw.WriteHeader(header); nil != err {
			return err
		}
		_, err = io.Copy(tw, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// FileDeCompressTar 压缩文件
// 压缩文件路径
// 解压文件夹
//func FileDeCompressTar(tarFile, dest string) error {
//	srcFile, err := os.Open(tarFile)
//	if err != nil {
//		return err
//	}
//	defer func() { _ = srcFile.Close() }()
//	reader := tar.NewReader(srcFile)
//	return deCompressTar(reader, dest)
//}

// FileDeCompressZip 解压
func FileDeCompressZip(zipFile, dest string) error {
	var (
		reader *zip.ReadCloser
		err    error
	)
	if reader, err = zip.OpenReader(zipFile); nil != err {
		return err
	}
	return deCompressZip(reader, dest)
}

// deCompress 解压
func deCompressZip(reader *zip.ReadCloser, dest string) error {
	defer func() { _ = reader.Close() }()
	for _, innerFile := range reader.File {
		info := innerFile.FileInfo()
		if info.IsDir() {
			err := os.MkdirAll(innerFile.Name, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}
		srcFile, err := innerFile.Open()
		if err != nil {
			continue
		}
		err = os.MkdirAll(dest, 0755)
		if err != nil {
			return err
		}
		filePath := filepath.Join(dest, innerFile.Name)
		if exist := FilePathExists(filePath); !exist {
			lastIndex := strings.LastIndex(filePath, "/")
			parentPath := filePath[0:lastIndex]
			if err := os.MkdirAll(parentPath, os.ModePerm); nil != err {
				return err
			}
		}
		newFile, err := os.Create(filePath)
		if err != nil {
			_ = srcFile.Close()
			continue
		}
		if _, err = io.Copy(newFile, srcFile); nil != err {
			_ = newFile.Close()
			_ = srcFile.Close()
			return err
		}
		_ = newFile.Close()
		_ = srcFile.Close()
	}
	return nil
}

// deCompressTar 解压
//func deCompressTar(reader *tar.Reader, dest string) error {
//	for {
//		if header,err := reader.Next(); nil==err {
//			info := header.FileInfo()
//			if info.IsDir() {
//				err := os.MkdirAll(header.Name, os.ModePerm)
//				if err != nil {
//					return err
//				}
//				continue
//			}
//		}
//	}
//	defer func() { _ = reader.Close() }()
//	for _, innerFile := range reader.File {
//		info := innerFile.FileInfo()
//		if info.IsDir() {
//			err := os.MkdirAll(innerFile.Name, os.ModePerm)
//			if err != nil {
//				return err
//			}
//			continue
//		}
//		srcFile, err := innerFile.Open()
//		if err != nil {
//			continue
//		}
//		err = os.MkdirAll(dest, 0755)
//		if err != nil {
//			return err
//		}
//		filePath := filepath.Join(dest, innerFile.Name)
//		if exist := FilePathExists(filePath); !exist {
//			lastIndex := strings.LastIndex(filePath, "/")
//			parentPath := filePath[0:lastIndex]
//			if err := os.MkdirAll(parentPath, os.ModePerm); nil != err {
//				return err
//			}
//		}
//		newFile, err := os.Create(filePath)
//		if err != nil {
//			_ = srcFile.Close()
//			continue
//		}
//		if _, err = io.Copy(newFile, srcFile); nil != err {
//			_ = newFile.Close()
//			_ = srcFile.Close()
//			return err
//		}
//		_ = newFile.Close()
//		_ = srcFile.Close()
//	}
//	return nil
//}

// FileCopy 拷贝文件
//
// srcFilePath 源文件路径
//
// dstFilePath 目标文件路径
func FileCopy(srcFilePath, dstFilePath string) (written int64, err error) {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return -1, err
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.OpenFile(dstFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return -1, err
	}
	defer func() { _ = dstFile.Close() }()
	return io.Copy(dstFile, srcFile)
}

// 图片处理
// ImageMimeTypes refer net/http package
var ImageMimeTypes = map[string]string{
	"bmp": "image/bmp",
	"gif": "image/gif",
	"ief": "image/ief",
	"jpg": "image/jpeg",
	// "jpe":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"svg":  "image/svg+xml",
	"ico":  "image/x-icon",
	"webp": "image/webp",
}

// IsImageFile check file is image file.
func IsImageFile(path string) bool {
	mime := MimeType(path)
	if mime == "" {
		return false
	}

	for _, imgMime := range ImageMimeTypes {
		if imgMime == mime {
			return true
		}
	}
	return false
}

// IsZipFile check is zip file.
// from https://blog.csdn.net/wangshubo1989/article/details/71743374
func IsZipFile(filepath string) bool {
	f, err := os.Open(filepath)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4)
	if n, err := f.Read(buf); err != nil || n < 4 {
		return false
	}

	return bytes.Equal(buf, []byte("PK\x03\x04"))
}

// FileFilter for filter file path.
type FileFilter interface {
	FilterFile(filePath, filename string) bool
}

// FileFilterFunc for filter file path.
type FileFilterFunc func(filePath, filename string) bool

// FilterFile Filter for filter file path.
func (fn FileFilterFunc) FilterFile(filePath, filename string) bool {
	return fn(filePath, filename)
}

// DirFilter for filter dir path.
type DirFilter interface {
	FilterDir(dirPath, dirName string) bool
}

// DirFilterFunc for filter file path.
type DirFilterFunc func(dirPath, dirName string) bool

// FilterDir Filter for filter file path.
func (fn DirFilterFunc) FilterDir(dirPath, dirName string) bool {
	return fn(dirPath, dirName)
}

// // BodyFilter for filter file contents.
// type BodyFilter interface {
// 	FilterBody(contents, filePath string) bool
// }
//
// // BodyFilterFunc for filter file contents.
// type BodyFilterFunc func(contents, filePath string) bool
//
// // Filter for filter file path.
// func (fn BodyFilterFunc) FilterBody(contents, filePath string) bool {
// 	return fn(contents, filePath)
// }

// FilterFunc for filter file path.
type FilterFunc func(filePath, filename string) bool

// Filter for filter file path.
func (fn FilterFunc) Filter(filePath, filename string) bool {
	return fn(filePath, filename)
}

// FileMeta struct
type FileMeta struct {
	filePath string
	filename string
}

// FindResults struct
type FindResults struct {
	f *FileFilter

	// founded file paths.
	filePaths []string

	// filters
	dirFilters  []DirFilter  // filters for filter dir paths
	fileFilters []FileFilter // filters for filter file paths
	// bodyFilters []BodyFilter // filters for filter file contents
}

func (r *FindResults) append(filePath ...string) {
	r.filePaths = append(r.filePaths, filePath...)
}

// AddFilters Result get find paths
func (r *FindResults) AddFilters(filterFuncs ...FileFilter) *FindResults {
	return r
}

// Filter Result get find paths
func (r *FindResults) Filter() *FindResults {
	return r
}

// Each Result get find paths
func (r *FindResults) Each() *FindResults {
	return r
}

// Result get find paths
func (r *FindResults) Result() []string {
	return r.filePaths
}

// TODO use excludeDotFlag 1 file 2 dir 1|2 both
type exDotFlag uint8

const (
	ExDotFile exDotFlag = 1
	ExDotDir  exDotFlag = 2
)

// FileFinder struct
type FileFinder struct {
	// r *FindResults

	// mark has been run find()
	founded bool
	// dir paths for find file.
	dirPaths []string
	// file paths for filter.
	srcFiles []string

	// builtin include filters
	includeDirs []string // include dir names. eg: {"model"}
	includeExts []string // include ext names. eg: {".go", ".md"}

	// builtin exclude filters
	excludeDirs  []string // exclude dir names. eg: {"test"}
	excludeExts  []string // exclude ext names. eg: {".go", ".md"}
	excludeNames []string // exclude file names. eg: {"go.mod"}

	// builtin dot filters.
	// TODO use excludeDotFlag 1 file 2 dir 1|2 both
	// excludeDotFlag exDotFlag
	excludeDotDir  bool
	excludeDotFile bool

	// fileFlags int

	dirFilters  []DirFilter  // filters for filter dir paths
	fileFilters []FileFilter // filters for filter file paths

	// founded file paths.
	filePaths []string

	// the founded file instances
	osFiles map[string]*os.File
	osInfos map[string]os.FileInfo
}

// EmptyFinder new empty FileFinder instance
func EmptyFinder() *FileFinder {
	return &FileFinder{
		osInfos: make(map[string]os.FileInfo),
	}
}

// NewFinder new instance with source dir paths.
func NewFinder(dirPaths []string, filePaths ...string) *FileFinder {
	return &FileFinder{
		dirPaths:  dirPaths,
		filePaths: filePaths,
		osInfos:   make(map[string]os.FileInfo),
	}
}

// AddDirPath add source dir for find
func (f *FileFinder) AddDirPath(dirPaths ...string) *FileFinder {
	f.dirPaths = append(f.dirPaths, dirPaths...)
	return f
}

// AddDir add source dir for find. alias of AddDirPath()
func (f *FileFinder) AddDir(dirPaths ...string) *FileFinder {
	f.dirPaths = append(f.dirPaths, dirPaths...)
	return f
}

// ExcludeDotDir exclude dot dir names. eg: ".idea"
func (f *FileFinder) ExcludeDotDir(exclude ...bool) *FileFinder {
	if len(exclude) > 0 {
		f.excludeDotDir = exclude[0]
	} else {
		f.excludeDotDir = true
	}
	return f
}

// NoDotDir exclude dot dir names. alias of ExcludeDotDir().
func (f *FileFinder) NoDotDir(exclude ...bool) *FileFinder {
	return f.ExcludeDotDir(exclude...)
}

// ExcludeDotFile exclude dot dir names. eg: ".gitignore"
func (f *FileFinder) ExcludeDotFile(exclude ...bool) *FileFinder {
	if len(exclude) > 0 {
		f.excludeDotFile = exclude[0]
	} else {
		f.excludeDotFile = true
	}
	return f
}

// NoDotFile exclude dot dir names. alias of ExcludeDotFile().
func (f *FileFinder) NoDotFile(exclude ...bool) *FileFinder {
	return f.ExcludeDotFile(exclude...)
}

// ExcludeDir exclude dir names.
func (f *FileFinder) ExcludeDir(dirs ...string) *FileFinder {
	f.excludeDirs = append(f.excludeDirs, dirs...)
	return f
}

// ExcludeName exclude file names.
func (f *FileFinder) ExcludeName(files ...string) *FileFinder {
	f.excludeNames = append(f.excludeNames, files...)
	return f
}

// AddFilter for filter filepath or dirpath
func (f *FileFinder) AddFilter(filterFuncs ...interface{}) *FileFinder {
	return f.WithFilter(filterFuncs...)
}

// WithFilter add filter func for filtering filepath or dirpath
func (f *FileFinder) WithFilter(filterFuncs ...interface{}) *FileFinder {
	for _, filterFunc := range filterFuncs {
		if fileFilter, ok := filterFunc.(FileFilter); ok {
			f.fileFilters = append(f.fileFilters, fileFilter)
		} else if dirFilter, ok := filterFunc.(DirFilter); ok {
			f.dirFilters = append(f.dirFilters, dirFilter)
		}
	}
	return f
}

// AddFileFilter for filter filepath
func (f *FileFinder) AddFileFilter(filterFuncs ...FileFilter) *FileFinder {
	f.fileFilters = append(f.fileFilters, filterFuncs...)
	return f
}

// WithFileFilter for filter func for filtering filepath
func (f *FileFinder) WithFileFilter(filterFuncs ...FileFilter) *FileFinder {
	f.fileFilters = append(f.fileFilters, filterFuncs...)
	return f
}

// AddDirFilter for filter file contents
func (f *FileFinder) AddDirFilter(filterFuncs ...DirFilter) *FileFinder {
	f.dirFilters = append(f.dirFilters, filterFuncs...)
	return f
}

// WithDirFilter for filter func for filtering file contents
func (f *FileFinder) WithDirFilter(filterFuncs ...DirFilter) *FileFinder {
	f.dirFilters = append(f.dirFilters, filterFuncs...)
	return f
}

// // AddBodyFilter for filter file contents
// func (f *FileFinder) AddBodyFilter(filterFuncs ...BodyFilter) *FileFinder {
// 	f.bodyFilters = append(f.bodyFilters, filterFuncs...)
// 	return f
// }
//
// // WithBodyFilter for filter func for filtering file contents
// func (f *FileFinder) WithBodyFilter(filterFuncs ...BodyFilter) *FileFinder {
// 	f.bodyFilters = append(f.bodyFilters, filterFuncs...)
// 	return f
// }

// AddFilePaths set founded files
func (f *FileFinder) AddFilePaths(filePaths []string) {
	f.filePaths = append(f.filePaths, filePaths...)
}

// AddFilePath add source file
func (f *FileFinder) AddFilePath(filePaths ...string) *FileFinder {
	f.filePaths = append(f.filePaths, filePaths...)
	return f
}

// AddFile add source file. alias of AddFilePath()
func (f *FileFinder) AddFile(filePaths ...string) *FileFinder {
	f.filePaths = append(f.filePaths, filePaths...)
	return f
}

// FindAll find and return founded file paths.
func (f *FileFinder) FindAll() []string {
	f.find()

	return f.filePaths
}

// Find find file paths.
func (f *FileFinder) Find() *FileFinder {
	f.find()
	return f
}

// do finding
func (f *FileFinder) find() {
	// mark found
	if f.founded {
		return
	}
	f.founded = true

	for _, filePath := range f.filePaths {
		fi, err := os.Stat(filePath)
		if err != nil {
			continue // ignore I/O error
		}
		if fi.IsDir() {
			continue // ignore I/O error
		}

		// cache file info
		f.osInfos[filePath] = fi
	}

	// do finding
	for _, dirPath := range f.dirPaths {
		f.findInDir(dirPath)
	}
}

// code refer filepath.glob()
func (f *FileFinder) findInDir(dirPath string) {
	dfi, err := os.Stat(dirPath)
	if err != nil {
		return // ignore I/O error
	}
	if !dfi.IsDir() {
		return // ignore I/O error
	}

	// opening
	d, err := os.Open(dirPath)
	if err != nil {
		return // ignore I/O error
	}

	names, _ := d.Readdirnames(-1)
	sort.Strings(names)

	hasDirFilter := len(f.dirFilters) > 0
	hasFileFilter := len(f.fileFilters) > 0
	for _, name := range names {
		fullPath := filepath.Join(dirPath, name)
		fi, err := os.Stat(fullPath)
		if err != nil {
			continue // ignore I/O error
		}

		// --- dir
		if fi.IsDir() {
			if f.excludeDotDir && name[0] == '.' {
				continue
			}

			var ok bool
			if hasDirFilter {
				for _, df := range f.dirFilters {
					ok = df.FilterDir(fullPath, name)
					if true == ok { // 有一个满足即可
						break
					}
				}

				// find in sub dir.
				if ok {
					f.findInDir(fullPath)
				}
			} else {
				// find in sub dir.
				f.findInDir(fullPath)
			}

			continue
		}

		// --- file
		if f.excludeDotFile && name[0] == '.' {
			continue
		}

		// use custom filter functions
		var ok bool
		if hasFileFilter {
			for _, ff := range f.fileFilters {
				ok = ff.FilterFile(fullPath, name)
				if true == ok { // 有一个满足即可
					break
				}
			}
		} else {
			ok = true
		}

		// append
		if ok {
			f.filePaths = append(f.filePaths, fullPath)
			// cache file info
			f.osInfos[fullPath] = fi
		}
	}

	_ = d.Close()
}

// Each each file paths.
func (f *FileFinder) Each(fn func(filePath string)) {
	// ensure find is running
	f.find()

	for _, filePath := range f.filePaths {
		fn(filePath)
	}
}

// EachFile each file os.File
func (f *FileFinder) EachFile(fn func(file *os.File)) {
	// ensure find is running
	f.find()

	for _, filePath := range f.filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			continue
		}

		fn(file)
	}
}

// EachStat each file os.FileInfo
func (f *FileFinder) EachStat(fn func(fi os.FileInfo, filePath string)) {
	// ensure find is running
	f.find()

	for filePath, fi := range f.osInfos {
		fn(fi, filePath)
	}
}

// EachContents each file contents
func (f *FileFinder) EachContents(fn func(contents, filePath string)) {
	// ensure find is running
	f.find()

	for _, filePath := range f.filePaths {
		bts, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}

		fn(string(bts), filePath)
	}
}

// Reset data setting.
func (f *FileFinder) Reset() {
	f.founded = false

	f.filePaths = make([]string, 0)

	f.excludeNames = make([]string, 0)
	f.excludeExts = make([]string, 0)
	f.excludeDirs = make([]string, 0)
}

// String all file paths
func (f *FileFinder) String() string {
	return strings.Join(f.filePaths, "\n")
}

//
// ------------------ built in file path filters ------------------
//

// ExtFilterFunc filter filepath by given file ext.
//
// Usage:
//	f := EmptyFiler()
//	f.AddFilter(ExtFilterFunc([]string{".go", ".md"}, true))
//	f.AddFilter(ExtFilterFunc([]string{".log", ".tmp"}, false))
func ExtFilterFunc(exts []string, include bool) FileFilterFunc {
	return func(filePath, _ string) bool {
		fExt := path.Ext(filePath)

		for _, ext := range exts {
			if fExt == ext {
				return include
			}
		}
		return !include
	}
}

// SuffixFilterFunc filter filepath by given file ext.
//
// Usage:
//	f := EmptyFiler()
//	f.AddFilter(SuffixFilterFunc([]string{"util.go", ".md"}, true))
//	f.AddFilter(SuffixFilterFunc([]string{"_test.go", ".log"}, false))
func SuffixFilterFunc(suffixes []string, include bool) FileFilterFunc {
	return func(filePath, _ string) bool {
		for _, sfx := range suffixes {
			if strings.HasSuffix(filePath, sfx) {
				return include
			}
		}
		return !include
	}
}

// PathNameFilterFunc filter filepath by given path names.
func PathNameFilterFunc(names []string, include bool) FileFilterFunc {
	return func(filePath, _ string) bool {
		for _, name := range names {
			if strings.Contains(filePath, name) {
				return include
			}
		}
		return !include
	}
}

// DotFileFilterFunc filter dot filename. eg: ".gitignore"
func DotFileFilterFunc(include bool) FileFilterFunc {
	return func(filePath, filename string) bool {
		// filename := path.Base(filePath)
		if filename[0] == '.' {
			return include
		}

		return !include
	}
}

// ModTimeFilterFunc filter file by modify time.
func ModTimeFilterFunc(limitSec int, op rune, include bool) FileFilterFunc {
	return func(filePath, filename string) bool {
		fi, err := os.Stat(filePath)
		if err != nil {
			return !include
		}

		now := time.Now().Second()
		if op == '>' {
			if now-fi.ModTime().Second() > limitSec {
				return include
			}

			return !include
		}

		// '<'
		if now-fi.ModTime().Second() < limitSec {
			return include
		}

		return !include
	}
}

// GlobFilterFunc filter filepath by given patterns.
//
// Usage:
//	f := EmptyFiler()
//	f.AddFilter(GlobFilterFunc([]string{"*_test.go"}, true))
func GlobFilterFunc(patterns []string, include bool) FileFilterFunc {
	return func(_, filename string) bool {
		for _, pattern := range patterns {
			if ok, _ := path.Match(pattern, filename); ok {
				return include
			}
		}
		return !include
	}
}

// RegexFilterFunc filter filepath by given regex pattern
//
// Usage:
//	f := EmptyFiler()
//	f.AddFilter(RegexFilterFunc(`[A-Z]\w+`, true))
func RegexFilterFunc(pattern string, include bool) FileFilterFunc {
	reg := regexp.MustCompile(pattern)

	return func(_, filename string) bool {
		return reg.MatchString(filename)
	}
}

//
// ----------------- built in dir path filters -----------------
//

// DotDirFilterFunc filter dot dirname. eg: ".idea"
func DotDirFilterFunc(include bool) DirFilterFunc {
	return func(_, dirname string) bool {
		if dirname[0] == '.' {
			return include
		}

		return !include
	}
}

// DirNameFilterFunc filter filepath by given dir names.
func DirNameFilterFunc(names []string, include bool) DirFilterFunc {
	return func(_, dirName string) bool {
		for _, name := range names {
			if dirName == name {
				return include
			}
		}
		return !include
	}
}

const (
	// MimeSniffLen sniff Length, use for detect file mime type
	MimeSniffLen = 512
)

// OSTempFile create an temp file on os.TempDir()
// Usage:
// 	fsutil.OSTempFile("example.*.txt")
func OSTempFile(pattern string) (*os.File, error) {
	return ioutil.TempFile(os.TempDir(), pattern)
}

// TempFile is alias of ioutil.TempFile
// Usage:
// 	fsutil.TempFile("", "example.*.txt")
func TempFile(dir, pattern string) (*os.File, error) {
	return ioutil.TempFile(dir, pattern)
}

// OSTempDir creates a new temp dir on os.TempDir and return the temp dir path
// Usage:
// 	fsutil.OSTempDir("example.*.txt")
func OSTempDir(pattern string) (string, error) {
	return ioutil.TempDir(os.TempDir(), pattern)
}

// TempDir creates a new temp dir and return the temp dir path
// Usage:
// 	fsutil.TempDir("", "example.*.txt")
func TempDir(dir, pattern string) (string, error) {
	return ioutil.TempDir(dir, pattern)
}

// ExpandPath will parse `~` as user home dir path.
func ExpandPath(path string) string {
	path, _ = homedir.Expand(path)
	return path
}

// Realpath parse and get
func Realpath(path string) string {
	// TODO
	return path
}

// MimeType get File Mime Type name. eg "image/png"
func MimeType(path string) (mime string) {
	if path == "" {
		return
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}

	return ReaderMimeType(file)
}

// ReaderMimeType get the io.Reader mimeType
// Usage:
// 	file, err := os.Open(filepath)
// 	if err != nil {
// 		return
// 	}
//	mime := ReaderMimeType(file)
func ReaderMimeType(r io.Reader) (mime string) {
	var buf [MimeSniffLen]byte
	n, _ := io.ReadFull(r, buf[:])
	if n == 0 {
		return ""
	}

	return http.DetectContentType(buf[:n])
}

// Mkdir alias of os.MkdirAll()
func Mkdir(dirPath string, perm os.FileMode) error {
	return os.MkdirAll(dirPath, perm)
}

// MkParentDir quick create parent dir
func MkParentDir(fpath string) error {
	dirPath := filepath.Dir(fpath)
	if !IsDir(dirPath) {
		return os.MkdirAll(dirPath, 0775)
	}
	return nil
}

// MustReadFile read file contents, will panic on error
func MustReadFile(filePath string) []byte {
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return bs
}

// ReadExistFile read file contents if exist, will panic on error
func ReadExistFile(filePath string) []byte {
	if IsFile(filePath) {
		bs, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		return bs
	}
	return nil
}

// ************************************************************
//	open/create files
// ************************************************************

// OpenFile like os.OpenFile, but will auto create dir.
func OpenFile(filepath string, flag int, perm os.FileMode) (*os.File, error) {
	fileDir := path.Dir(filepath)

	// if err := os.Mkdir(dir, 0775); err != nil {
	if err := os.MkdirAll(fileDir, DefaultDirPerm); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filepath, flag, perm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// QuickOpenFile like os.OpenFile
func QuickOpenFile(filepath string) (*os.File, error) {
	return OpenFile(filepath, DefaultFileFlags, DefaultFilePerm)
}

/* TODO MustOpenFile() */

// CreateFile create file if not exists
// Usage:
// 	CreateFile("path/to/file.txt", 0664, 0666)
func CreateFile(fpath string, filePerm, dirPerm os.FileMode) (*os.File, error) {
	dirPath := path.Dir(fpath)
	if !IsDir(dirPath) {
		err := os.MkdirAll(dirPath, dirPerm)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
}

// MustCreateFile create file, will panic on error
func MustCreateFile(filePath string, filePerm, dirPerm os.FileMode) *os.File {
	file, err := CreateFile(filePath, filePerm, dirPerm)
	if err != nil {
		panic(err)
	}

	return file
}

// ************************************************************
//	copy files
// ************************************************************

// CopyFile copy file to another path.
func CopyFile(src string, dst string) error {
	return errors.New("TODO")
}

// MustCopyFile copy file to another path.
func MustCopyFile(src string, dst string) {
	panic("TODO")
}

// ************************************************************
//	remove files
// ************************************************************

// alias methods
var (
	MustRm  = MustRemove
	QuietRm = QuietRemove
)

// MustRemove removes the named file or (empty) directory.
// NOTICE: if error will panic
func MustRemove(fpath string) {
	if err := os.Remove(fpath); err != nil {
		panic(err)
	}
}

// QuietRemove removes the named file or (empty) directory.
// NOTICE: will ignore error
func QuietRemove(fpath string) {
	_ = os.Remove(fpath)
}

// DeleteIfExist removes the named file or (empty) directory on exists.
func DeleteIfExist(fpath string) error {
	if !FilePathExists(fpath) {
		return nil
	}

	return os.Remove(fpath)
}

// DeleteIfFileExist removes the named file on exists.
func DeleteIfFileExist(fpath string) error {
	if !IsFile(fpath) {
		return nil
	}

	return os.Remove(fpath)
}

// ************************************************************
//	other operates
// ************************************************************

// Unzip a zip archive
// from https://blog.csdn.net/wangshubo1989/article/details/71743374
func Unzip(archive, targetDir string) (err error) {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return
	}

	if err = os.MkdirAll(targetDir, DefaultDirPerm); err != nil {
		return
	}

	for _, file := range reader.File {
		fullPath := filepath.Join(targetDir, file.Name)
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(fullPath, file.Mode())
			if err != nil {
				return err
			}

			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		targetFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()
			return err
		}

		_, err = io.Copy(targetFile, fileReader)

		// close all
		fileReader.Close()
		targetFile.Close()

		if err != nil {
			return err
		}
	}

	return
}
