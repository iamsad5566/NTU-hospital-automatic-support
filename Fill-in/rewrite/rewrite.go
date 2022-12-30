package rewrite

import (
	"fillin/config"
	"fillin/model"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func StartReWrite(setting config.Config) {
	data := read(setting)
	query := cache(data, setting)
	writeCSV(query, setting.OutputFile)
}

func read(setting config.Config) map[string][]string {
	f, _ := excelize.OpenFile(setting.FileName)
	defer f.Close()
	rows, _ := f.GetRows(setting.SourceSheet)
	targetsMap := model.GetDecimalSlice(setting.ReWriteColumns)

	output := make(map[string][]string)
	ch := make(chan []string)

	for _, row := range rows {
		go func(ch chan []string, row []string) {
			newSlice := make([]string, len(setting.ReWriteColumns)+1)

			pk := row[model.ToDecimal(setting.PrimaryKey)]
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

func cache(input map[string][]string, setting config.Config) []string {
	f, err := excelize.OpenFile(setting.FileName)

	cacheStr := make([]string, 0)

	if err != nil {
		log.Println(err)
		defer f.Close()
		return cacheStr
	}

	startPoint := 2
	rows, _ := f.GetRows(setting.TargetSheet)
	cacheStr = append(cacheStr, model.FilterTargetColumn(rows[0], setting.ReWriteColumns))
	endPoint := len(rows)
	ch := make(chan bool, 10)

	for i := startPoint; i <= endPoint; i++ {
		go func(ch chan bool, i int, cacheStr *[]string) {
			indexStr := strconv.Itoa(i)
			pk, _ := f.GetCellValue(setting.TargetSheet, setting.PrimaryKey+indexStr)
			for _, splitTarget := range setting.Split {
				pk = strings.ReplaceAll(pk, splitTarget, "-")
			}
			pk = strings.ToLower(pk)

			if len(input[pk]) > 0 {
				*cacheStr = append(*cacheStr, pk+","+model.SliceToStr(input[pk]))
			}
			ch <- true
		}(ch, i, &cacheStr)
	}

	for i := startPoint; i <= endPoint; i++ {
		<-ch
	}

	defer f.Close()
	return cacheStr
}

func writeCSV(query []string, fileName string) bool {
	file, _ := os.Create(fileName)
	for _, str := range query {
		file.WriteString(str + "\n")
	}
	file.Close()
	return true
}
