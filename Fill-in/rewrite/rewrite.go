package rewrite

import (
	"fillin/config"
	"log"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func StartReWrite(setting config.Config) {
	data := read(setting)
	write(data, setting)
}

func read(setting config.Config) map[string][]string {
	f, _ := excelize.OpenFile(setting.FileName)
	rows, _ := f.GetRows(setting.SourceSheet)
	targetsMap := getDecimalSlice(setting.ReWriteColumns)

	output := make(map[string][]string)
	ch := make(chan []string)

	for _, row := range rows {
		go func(ch chan []string, row []string) {
			newSlice := make([]string, len(setting.ReWriteColumns)+1)

			pk := row[toDecimal(setting.PrimaryKey)]
			for _, splitTarget := range setting.Split {
				pk = strings.ReplaceAll(pk, splitTarget, "-")
			}

			newSlice[0] = strings.ToLower(pk)

			i := 1
			for col, val := range row {
				if targetsMap[col] {
					newSlice[i] = val
					i++
				}
			}
			ch <- newSlice
		}(ch, row)
	}

	for i := 0; i < len(rows); i++ {
		newSlice := <-ch
		key := newSlice[0]
		newSlice = newSlice[1:]
		output[key] = newSlice
	}
	return output
}

func write(input map[string][]string, setting config.Config) bool {
	f, err := excelize.OpenFile(setting.FileName)
	if err != nil {
		log.Println(err)
		return false
	}

	startPoint := 2
	rows, _ := f.GetRows(setting.TargetSheet)
	endPoint := len(rows)
	ch := make(chan int)

	for i := startPoint; i <= endPoint; i++ {
		go func(ch chan int, i int) {
			indexStr := strconv.Itoa(i)
			pk, _ := f.GetCellValue(setting.TargetSheet, setting.PrimaryKey+indexStr)
			for _, splitTarget := range setting.Split {
				pk = strings.ReplaceAll(pk, splitTarget, "-")
			}
			pk = strings.ToLower(pk)

			for j := 0; j < len(setting.ReWriteColumns); j++ {
				if len(input[pk]) > 0 {
					f.SetCellValue(setting.TargetSheet, setting.ReWriteColumns[j]+indexStr, input[pk][j])
				}
			}
			ch <- i
		}(ch, i)
	}

	for i := startPoint; i <= endPoint; i++ {
		<-ch
	}

	err = f.Save()
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func getDecimalSlice(strColumns []string) map[int]bool {
	output := make(map[int]bool)
	for _, str := range strColumns {
		output[toDecimal(str)] = true
	}
	return output
}

func toDecimal(str string) int {
	multiple := getMultiple(len(str) - 1)
	decimalIndex := 0
	for len(str) > 0 {
		num := int(str[0]-'A') + 1
		decimalIndex += num * multiple
		str = str[1:]
		multiple /= 26
	}
	return decimalIndex - 1
}

func getMultiple(digit int) int {
	output := 1
	for digit > 0 {
		output *= 26
		digit--
	}
	return output
}
