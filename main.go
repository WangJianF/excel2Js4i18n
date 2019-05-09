package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	f, err := excelize.OpenFile("./asset/i18n.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	result := make(map[string]map[string]map[string]string)
	colName := make([]string, 10)
	for _, name := range f.GetSheetMap() {
		rows, err := f.Rows(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		rowIndex := 0
		for rows.Next() {
			rowIndex++
			if rowIndex == 1 {
				continue
			}
			row, err := rows.Columns()
			if err != nil {
				fmt.Println(err)
				return
			}
			if rowIndex == 2 {
				for colIndex, colCell := range row {
					if colIndex == 0 {
						continue
					}
					if _, ok := result[colCell]; !ok {
						result[colCell] = make(map[string]map[string]string)
					}
					if _, ok := result[colCell][name]; !ok {
						result[colCell][name] = make(map[string]string)
					}
					colName[colIndex] = colCell
				}
				continue
			}
			for colIndex, colCell := range row {
				if colIndex < 1 {
					continue
				}
				result[colName[colIndex]][name][row[0]] = colCell
			}
		}
	}
	for key, value := range result {
		jsonStr, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err)
		}

		str := "if (!window.i18n) window.i18n = {};\nif (!window.i18n.languages) window.i18n.languages = {};\n"
		str += "window.i18n.languages." + key + "=" + string(jsonStr[:]) + ";"
		file := "./target/" + key + ".js"
		ioutil.WriteFile(file, []byte(str), 0644)
	}
}
