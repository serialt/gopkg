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
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// 当前项目根目录
var API_ROOT string

// 获取项目路径
func GetPath() string {

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

// 创建文件
func MkdirFile(path string) error {

	err := os.Mkdir(path, os.ModePerm) //在当前目录下生成md目录
	if err != nil {
		return err
	}
	return nil
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
