package excel

import (
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
	return e.File.SaveAs(sheetName)
}

func (e *Excel) SetCols(cols []string, line int) error {
	colsIndex := e.MakeColumns(len(cols))
	if line != 0 {
		e.CurrentRow = line
	}

	currentCol := 0
	for _, colName := range cols {
		cell := colsIndex[currentCol] + strconv.Itoa(line)
		currentCol++
		err := e.Write(cell, colName)
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
	return nil
}

func (e *Excel) Write(cell string, value any) error {
	return e.File.SetCellValue(e.CurrentSheet, cell, value)
}

// GetColumn converts a 1-based column index to an Excel column name (A, B, ..., Z, AA, AB, ...).
func (e *Excel) GetColumn(index int) string {
	if index <= 0 {
		return ""
	}
	result := ""
	for index > 0 {
		index-- // Excel columns are 1-based
		char := rune('A' + (index % 26))
		result = string(char) + result
		index /= 26
	}
	return result
}

// MakeColumns returns a list of Excel column
// names from A to the column at the given 1-based index (e.g., A, B, ..., Z, AA, AB...).
func (e *Excel) MakeColumns(max int) []string {
	if max <= 0 {
		return nil
	}
	cols := make([]string, max)
	for i := 1; i <= max; i++ {
		cols[i-1] = e.GetColumn(i)
	}
	return cols
}
