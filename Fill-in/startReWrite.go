package main

import (
	"fillin/config"
	"fillin/rewrite"
)


func main() {
	setting := config.GetInstance()
	rewrite.StartReWrite(setting)
}