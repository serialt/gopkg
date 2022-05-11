package gopkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasreList2Excel(t *testing.T) {
	// rootPath := GetRootPath()
	rootPath := "/Users/serialt/Desktop/flkj/github/gopkg/tmp/"
	var data [][]string

	header := []string{"ID", "name", "email"}
	data = append(data, []string{"1001", "Tom", "tom@imau.cc"})
	data = append(data, []string{"1003", "Jerry", "jerry@imau.cc"})
	filepath := fmt.Sprintf("%s/excel.xlsx", rootPath)
	err := PasreList2Excel(filepath, "Sheet1", header, data)
	t.Logf("写入的数据: %v", data)
	assert.Empty(t, err)
}

func TestPasreExcel2List(t *testing.T) {
	filepath := fmt.Sprintf("%s/excel.xlsx", "/Users/serialt/Desktop/flkj/github/gopkg/tmp/")
	data, err := PasreExcel2List(filepath, []string{"ID", "name", "email"}, false)
	t.Logf("读取的数据: %v", data)
	assert.Empty(t, err)
	assert.True(t, ForEqualStringSlice(data[1], []string{"1001", "Tom", "tom@imau.cc"}))
	assert.True(t, ForEqualStringSlice(data[2], []string{"1003", "Jerry", "jerry@imau.cc"}))
}
