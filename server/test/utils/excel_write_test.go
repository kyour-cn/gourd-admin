package main

import (
	"app/internal/util/excel"
	"testing"
)

func TestExcelWrite(t *testing.T) {

	fileName := "Book1.xlsx"

	e := excel.NewExcel()

	err := e.SetCols([]excel.Column{
		{Name: "ID", Width: 10},
		{Name: "Name", Width: 20},
		{Name: "Age", Width: 10},
	}, 1)
	if err != nil {
		t.Error(err)
	}

	data := [][]any{
		{"1", "张三", 18},
		{"2", "李四", 19},
		{"3", "王五", 20},
		{"4", "赵六", 21},
	}
	for _, row := range data {
		err := e.WriteLine(e.CurrentRow, row)
		if err != nil {
			t.Error(err)
		}
	}
	err = e.Write("A7", "测试内容")
	if err != nil {
		t.Error(err)
	}

	err = e.Save(fileName)
	if err != nil {
		t.Error(err)
	}

	err = e.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Success file: %s", fileName)
}
