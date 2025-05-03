package main

import (
	"github.com/ghulammuzz/misterblast/config"
	"github.com/ghulammuzz/misterblast/internal/app"
)

func init() {
	config.Init()
}

func main() {
	app.Start()
}
