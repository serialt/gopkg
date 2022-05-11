/*
 * @Description   : IMAU of Serialt
 * @Author        : serialt
 * @Email         : serialt@qq.com
 * @Github        : https://github.com/serialt
 * @Created Time  : 2022-05-09 10:35:18
 * @Last modified : 2022-05-11 09:31:37
 * @FilePath      : /gopkg/excelize.go
 * @Other         :
 * @              :
 *
 *                 人和代码，有一个能跑就行
 *
 */

package gopkg

import (
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

/**
 * @description: 读取excel文件数据
 * @author: serialt
 * @param filepath {string}
 * @param fileHeader {[]string}
 * @param skipHeader {bool}
 * @return {error}
 */
func PasreExcel2List(filepath string, fileHeader []string, skipHeader bool) (data [][]string, err error) {

	file, err := excelize.OpenFile(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return
	}
	if skipHeader {
		if ForEqualStringSlice(rows[0], fileHeader) {
			rows = rows[1:]
		} else {
			err = errors.New("fileHeader not Equal")
			return
		}
	}
	data = rows
	return
}

/**
 * @description: 写入数据到excel文件
 * @author: serialt
 * @param filepath {string}
 * @param sheet {string}
 * @param data {[][]string}
 * @return {error}
 */
func PasreList2Excel(filepath, sheet string, fileHeader []string, data [][]string) (err error) {
	if IsFile(filepath) {
		err = errors.New("File is exists")
		return
	}
	excel := excelize.NewFile()
	excel.SetSheetRow(sheet, "A1", &fileHeader)
	for index, row := range data {
		axis := fmt.Sprintf("A%d", index+2)
		excel.SetSheetRow(sheet, axis, &row)

	}
	err = excel.SaveAs(filepath)
	return

}
