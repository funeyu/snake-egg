package util

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CsvData interface {
	Headers() []string
	Bodys() [][] string
}

func IsExsit(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func WriteCsv(path string, data CsvData) {
	var file *os.File
	var isExsit bool
	if IsExsit(path) {
		isExsit = true
		f, err := os.OpenFile(path,  os.O_WRONLY | os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("open file failed:", err)
		}
		file = f
	} else {
		isExsit = false
		f, err := os.Create(path)
		file.WriteString("\xEF\xBB\xBF") //防止中文乱码
		if err != nil {
			fmt.Println("create file failed:", err)
		}
		file = f
	}

	defer file.Close()
	w := csv.NewWriter(file)
	if !isExsit {
		headers := data.Headers()
		w.Write(headers)
	}

	for _, d := range data.Bodys() {
		w.Write(d)
	}
	w.Flush()
}

func ReadCsv(path string) [][]string {
	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}