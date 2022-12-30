package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	FileName       string
	OutputFile     string
	TargetSheet    string
	SourceSheet    string
	PrimaryKey     string
	Split          []string
	ReWriteColumns []string
}

func GetInstance() Config {
	godotenv.Load()
	newConfig := Config{}
	newConfig.FileName = os.Getenv("FILE_NAME")
	newConfig.OutputFile = os.Getenv("OUTPUT_FILE")
	newConfig.TargetSheet = os.Getenv("TARGET_SHEET")
	newConfig.SourceSheet = os.Getenv("SOURCE_SHEET")
	newConfig.PrimaryKey = os.Getenv("PRIMARY_KEY")
	newConfig.Split = strings.Split(os.Getenv("SPLIT"), ",")
	newConfig.ReWriteColumns = strings.Split(os.Getenv("RE_WRITE_COLUMNS"), ",")
	return newConfig
}
