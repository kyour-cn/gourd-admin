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

	File         *excelize.File
	StreamWriter *excelize.StreamWriter
}

// OpenFile opens an Excel file.
func (e *Excel) OpenFile(filePath string) error {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	e.File = file
	return nil
}

func (e *Excel) SetSheet(sheetName string) error {
	index, err := e.File.GetSheetIndex(sheetName)
	if err != nil {
		return err
	}
	e.File.SetActiveSheet(index)
	e.CurrentSheetIndex = index
	e.CurrentSheet = sheetName
	return nil
}

func (e *Excel) Read() ([][]string, error) {
	rows, err := e.File.GetRows(e.CurrentSheet)
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

// StreamWriteBegin begins a stream write operation for the current sheet.
func (e *Excel) StreamWriteBegin(sheetName string) error {
	err := e.SetSheet(sheetName)
	if err != nil {
		return err
	}
	sw, err := e.File.NewStreamWriter(e.CurrentSheet)
	if err != nil {
		return err
	}

	e.StreamWriter = sw
	return nil
}

// StreamWriteEnd ends a stream write operation for the current sheet.
func (e *Excel) StreamWriteEnd() error {
	if e.StreamWriter != nil {
		err := e.StreamWriter.Flush()
		if err != nil {
			return err
		}
		e.StreamWriter = nil
	}
	return nil
}

// Save to file
func (e *Excel) Save(sheetName string) error {

	// 确保流式写入结束
	if e.StreamWriter == nil {
		err := e.StreamWriteEnd()
		if err != nil {
			return err
		}
	}

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

// SetCols sets the columns of the current sheet.
func (e *Excel) SetCols(cols []Column, line int, freeze bool) error {
	if e.StreamWriter == nil {
		// 冻结表头
		if freeze {
			err := e.File.SetPanes(e.CurrentSheet, &excelize.Panes{
				Freeze:      true,
				YSplit:      e.CurrentRow,                       // 冻结
				TopLeftCell: "A" + strconv.Itoa(e.CurrentRow+1), // 活动单元格位置
			})
			if err != nil {
				return err
			}
		}

		colsIndex := e.MakeColumns(len(cols))
		if line != 0 {
			e.CurrentRow = line
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
	} else {
		// 冻结表头
		if freeze {
			err := e.StreamWriter.SetPanes(&excelize.Panes{
				Freeze:      true,
				YSplit:      e.CurrentRow,                       // 冻结
				TopLeftCell: "A" + strconv.Itoa(e.CurrentRow+1), // 活动单元格位置
			})
			if err != nil {
				return err
			}
		}

		row := make([]any, len(cols))
		for i, col := range cols {
			row[i] = col.Name
			// 设置列宽
			err := e.StreamWriter.SetColWidth(i+1, i+1, col.Width)
			if err != nil {
				return err
			}
		}

		// 写入表头文字
		cell, err := excelize.CoordinatesToCellName(1, e.CurrentRow)
		if err != nil {
			return err
		}
		if err := e.StreamWriter.SetRow(cell, row); err != nil {
			return err
		}
	}

	e.CurrentRow++

	return nil
}

// WriteLine writes a line of data to the current sheet.
func (e *Excel) WriteLine(line int, vals []any) error {
	if e.StreamWriter == nil {
		e.CurrentRow = line
		colsIndex := e.MakeColumns(len(vals))
		for i, val := range vals {
			cell := colsIndex[i] + strconv.Itoa(line)
			err := e.Write(cell, val)
			if err != nil {
				return err
			}
		}
	} else {
		cell, err := excelize.CoordinatesToCellName(1, line)
		if err != nil {
			return err
		}
		if err := e.StreamWriter.SetRow(cell, vals); err != nil {
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
