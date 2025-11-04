package main

import (
	"app/internal/util/excel"
	"testing"
)

func TestExcelRead(t *testing.T) {

	fileName := "Book1.xlsx"

	e := excel.NewExcel()

	err := e.OpenFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	err = e.Close()
	if err != nil {
		t.Fatal(err)
	}

	data, err := e.Read()
	if err != nil {
		t.Fatal(err)
	}

	_ = data

	t.Logf("Success data: %s", data)
}
