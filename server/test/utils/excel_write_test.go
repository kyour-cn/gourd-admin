package main

import (
	"app/internal/util/excel"
	"testing"
)

func TestExcelWrite(t *testing.T) {

	fileName := "Book1.xlsx"

	e := excel.NewExcel()

	err := e.SetCols([]string{
		"ID",
		"Name",
		"Age",
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
	for i, row := range data {
		err := e.WriteLine(i+2, row)
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
