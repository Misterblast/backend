package main

import (
	"github.com/ghulammuzz/misterblast/config"
	"github.com/ghulammuzz/misterblast/internal/app"
	metrics "github.com/ghulammuzz/misterblast/pkg/prom"
)

func init() {
	config.Init()
	metrics.Init()
}

func main() {
	app.Start()
}
