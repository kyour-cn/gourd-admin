package excel

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func NewExcel() *Excel {
	f := excelize.NewFile()

	sheetIndex := f.GetActiveSheetIndex()

	return &Excel{
		File:              f,
		CurrentRow:        1,
		CurrentSheetIndex: sheetIndex,
		CurrentSheet:      f.GetSheetName(sheetIndex),
	}
}

type Excel struct {
	CurrentRow        int    // 当前行
	CurrentSheetIndex int    // 当前sheet索引
	CurrentSheet      string // 当前sheet名称

	File *excelize.File
}

func (e *Excel) OpenFile(filePath string) error {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	e.File = file
	return nil
}

func (e *Excel) Read() ([][]string, error) {
	rows, err := e.File.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (e *Excel) Close() error {
	if e.File != nil {
		return e.File.Close()
	}
	return nil
}

// Save to file
func (e *Excel) Save(sheetName string) error {
	// 创建目录
	dir := filepath.Dir(sheetName)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return e.File.SaveAs(sheetName)
}

type Column struct {
	Name  string
	Width float64
}

func (e *Excel) SetCols(cols []Column, line int) error {
	colsIndex := e.MakeColumns(len(cols))
	if line != 0 {
		e.CurrentRow = line
	}

	// 冻结表头
	if err := e.File.SetPanes(e.CurrentSheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      e.CurrentRow,                       // 冻结
		TopLeftCell: "A" + strconv.Itoa(e.CurrentRow+1), // 活动单元格位置
		//ActivePane: "bottomLeft", // 激活底部窗格
	}); err != nil {
		return err
	}

	for i, col := range cols {
		cell := colsIndex[i] + strconv.Itoa(line)
		err := e.Write(cell, col.Name)
		if err != nil {
			return err
		}

		// 设置列宽
		err = e.File.SetColWidth(e.CurrentSheet, colsIndex[i], colsIndex[i], col.Width)
		if err != nil {
			return err
		}
	}
	e.CurrentRow++

	return nil
}

func (e *Excel) WriteLine(line int, vals []any) error {
	e.CurrentRow = line
	colsIndex := e.MakeColumns(len(vals))
	for i, val := range vals {
		cell := colsIndex[i] + strconv.Itoa(line)
		err := e.Write(cell, val)
		if err != nil {
			return err
		}
	}
	e.CurrentRow++
	return nil
}

func (e *Excel) Write(cell string, value any) error {
	return e.File.SetCellValue(e.CurrentSheet, cell, value)
}

// MakeColumns returns a list of Excel column
// names from A to the column at the given 1-based index (e.g., A, B, ..., Z, AA, AB...).
func (e *Excel) MakeColumns(max int) []string {
	cols := make([]string, max)
	for i := 1; i <= max; i++ {
		c, err := excelize.ColumnNumberToName(i)
		if err != nil {
			break
		}
		cols[i-1] = c
	}
	return cols
}
