package sugar

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileCommon_PathExists(t *testing.T) {
	path := "/etc"
	exist := FilePathExists(path)
	t.Log(path, "exist =", exist)

	path = "/etc/hello"
	exist = FilePathExists(path)
	t.Log(path, "exist =", exist)
}

func TestFileCommon_ReadFirstLine(t *testing.T) {
	profile, err := FileReadFirstLine("/etc/profile")
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadFirstLine_Fail(t *testing.T) {
	_, err := FileReadFirstLine("/etc/hello")
	t.Skip(err)
}

func TestFileCommon_ReadPointLine(t *testing.T) {
	profile, err := FileReadPointLine("/etc/profile", 1)
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadPointLine_KeyPoint(t *testing.T) {
	_, _ = FileAppend("./tmp/log/yes/go/point.txt", []byte("haha"), false)
	profile, err := FileReadPointLine("./tmp/log/yes/go/point.txt", 1)
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadPointLine_Fail_IndexOut(t *testing.T) {
	_, err := FileReadPointLine("/etc/profile", 300)
	t.Skip(err)
}

func TestFileCommon_ReadPointLine_Fail_NotExist(t *testing.T) {
	_, err := FileReadPointLine("/etc/hello", 1)
	t.Skip(err)
}

func TestFileCommon_ReadLines(t *testing.T) {
	profile, err := FileReadLines("/etc/profile")
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadLines_Fail(t *testing.T) {
	_, err := FileReadLines("/etc/hello")
	t.Skip(err)
}

func TestFileCommon_ParentPath(t *testing.T) {
	t.Log(FileParentPath("/etc/yes/go/test.txt"))
}

func TestFileCommon_Append(t *testing.T) {
	if _, err := FileAppend("./tmp/log/yes/go/test.txt", []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Append_Force(t *testing.T) {
	if _, err := FileAppend("./tmp/log/yes/go/test.txt", []byte("haha"), true); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Append_UnForce(t *testing.T) {
	if _, err := FileAppend("./tmp/log/yes/go/test.txt", []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Append_Fail_PermissionFileForce(t *testing.T) {
	_, err := FileAppend("/etc/www.json", []byte("haha"), true)
	t.Skip(err)
}

func TestFileCommon_Append_Fail_PermissionFileUnForce(t *testing.T) {
	_, err := FileAppend("/etc/www.json", []byte("haha"), false)
	t.Skip(err)
}

func TestFileCommon_Modify(t *testing.T) {
	if _, err := FileModify("./tmp/log/yes/go/test.txt", 1, []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Modify_Force(t *testing.T) {
	if _, err := FileModify("./tmp/log/yes/go/test.txt", 1, []byte("haha"), true); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Modify_UnForce(t *testing.T) {
	if _, err := FileModify("./tmp/log/yes/go/test.txt", 1, []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Modify_Fail_PermissionFileForce(t *testing.T) {
	_, err := FileModify("/etc/www.json", 1, []byte("haha"), true)
	t.Skip(err)
}

func TestFileCommon_Modify_Fail_PermissionFileUnForce(t *testing.T) {
	_, err := FileModify("/etc/www.json", 1, []byte("haha"), false)
	t.Skip(err)
}

func TestFileCommon_LoopDirs(t *testing.T) {
	if arr, err := FileLoopDirs("./tmp/log"); nil != err {
		t.Skip(err)
	} else {
		t.Log(arr)
	}
}

func TestFileCommon_LoopDirs_Fail(t *testing.T) {
	_, err := FileLoopDirs("./tmp/logger")
	t.Skip(err)
}

func TestFileCommon_LoopFiles(t *testing.T) {
	if arr, err := FileLoopFiles("./tmp/log"); nil != err {
		t.Skip(err)
	} else {
		t.Log(arr)
	}
}

func TestFileCommon_LoopFiles_Fail(t *testing.T) {
	_, err := FileLoopFiles("./tmp/logger")
	t.Skip(err)
}

func TestFileCommon_LoopOneDirs(t *testing.T) {
	array, err := FileLoopOneDirs("./tmp")
	t.Skip(array, err)
}

func TestFileCommon_Copy(t *testing.T) {
	if _, err := FileAppend("./tmp/copy/1.txt", []byte("hello"), true); nil != err {
		t.Skip(err)
	}
	if _, err := FileCopy("./tmp/copy/1.txt", "./tmp/copy/2.txt"); nil != err {
		t.Skip(err)
	}
}

func TestFileCompress(t *testing.T) {
	f, err := os.Open("./tmp/copy")
	assert.Nil(t, err)
	err = FileCompressZip([]*os.File{f}, "./tmp/copy.zip")
	assert.Nil(t, err)
}

func TestFileCompressTar(t *testing.T) {
	f, err := os.Open("./tmp/copy")
	assert.Nil(t, err)
	err = FileCompressTar([]*os.File{f}, "./tmp/copy.tar")
	assert.Nil(t, err)
	//err = FileDeCompressTar("./example/cas.tar", "./example/castar")
	//assert.Nil(t, err)
}

func TestFileDeCompressZip(t *testing.T) {
	err := FileDeCompressZip("./tmp/copy.zip", "./tmp/log/")
	assert.Nil(t, err)
	// err = FileDeCompressZip("./example/sql.zip", "./example/sql_de")
	// assert.Nil(t, err)
}

func TestCommon(t *testing.T) {
	assert.Equal(t, "", FileExt("testdata/testjpg"))
	assert.Equal(t, ".jpg", Suffix("testdata/test.jpg"))

	// IsZipFile
	assert.False(t, IsZipFile("testdata/test.jpg"))

	assert.Equal(t, "test.jpg", Name("path/to/test.jpg"))

	if runtime.GOOS == "windows" {
		assert.Equal(t, "path\\to", PathDir("path/to/test.jpg"))
	} else {
		assert.Equal(t, "path/to", PathDir("path/to/test.jpg"))
	}
}

func TestPathExists(t *testing.T) {
	assert.False(t, IsDir("/not-exist"))
	assert.False(t, IsDirExists("/not-exist"))
	assert.True(t, IsFile("testdata/test.jpg"))
	assert.True(t, IsFile("testdata/test.jpg"))
}

func TestIsFile(t *testing.T) {
	assert.False(t, FileExists(""))
	assert.False(t, IsFile(""))
	assert.False(t, IsFile("/not-exist"))
	assert.False(t, FileExists("/not-exist"))
	assert.True(t, IsFile("testdata/test.jpg"))
	assert.True(t, FileExist("testdata/test.jpg"))
}

func TestIsDir(t *testing.T) {
	assert.False(t, IsDir(""))
	assert.False(t, DirExist(""))
	assert.False(t, IsDir("/not-exist"))
	assert.True(t, IsDir("testdata"))
	assert.True(t, DirExist("testdata"))
}

func TestIsAbsPath(t *testing.T) {
	assert.True(t, IsAbsPath("/data/some.txt"))

	assert.NoError(t, DeleteIfFileExist("/not-exist"))
}

func TestMkdir(t *testing.T) {
	// TODO windows will error
	if IsWin() {
		return
	}

	err := os.Chmod("./testdata", os.ModePerm)

	if assert.NoError(t, err) {
		assert.NoError(t, Mkdir("./testdata/sub/sub21", os.ModePerm))
		assert.NoError(t, Mkdir("./testdata/sub/sub22", 0666))
		assert.NoError(t, Mkdir("./testdata/sub/sub23/sub31", 0777)) // 066X will error

		assert.NoError(t, os.RemoveAll("./testdata/sub"))
	}
}

func TestCreateFile(t *testing.T) {
	// TODO windows will error
	// if envutil.IsWin() {
	// 	return
	// }

	file, err := CreateFile("./testdata/test.txt", 0664, 0666)
	if assert.NoError(t, err) {
		assert.Equal(t, "./testdata/test.txt", file.Name())
		assert.NoError(t, file.Close())
		assert.NoError(t, os.Remove(file.Name()))
	}

	file, err = CreateFile("./testdata/sub/test.txt", 0664, 0777)
	if assert.NoError(t, err) {
		assert.Equal(t, "./testdata/sub/test.txt", file.Name())
		assert.NoError(t, file.Close())
		assert.NoError(t, os.RemoveAll("./testdata/sub"))
	}

	file, err = CreateFile("./testdata/sub/sub2/test.txt", 0664, 0777)
	if assert.NoError(t, err) {
		assert.Equal(t, "./testdata/sub/sub2/test.txt", file.Name())
		assert.NoError(t, file.Close())
		assert.NoError(t, os.RemoveAll("./testdata/sub"))
	}
}

func TestQuickOpenFile(t *testing.T) {
	fname := "./testdata/quick-open-file.txt"
	file, err := QuickOpenFile(fname)
	if assert.NoError(t, err) {
		assert.Equal(t, fname, file.Name())
		assert.NoError(t, file.Close())
		assert.NoError(t, os.Remove(file.Name()))
	}
}
