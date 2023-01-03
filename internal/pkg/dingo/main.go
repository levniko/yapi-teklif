package dingo

import (
	"log"

	"github.com/sarulabs/di/v2"
	"github.com/yapi-teklif/internal/pkg/dingo/container/dic"
)

var Application *App

type App struct {
	Container *dic.Container
}

func New() {
	Application = &App{}
	container, err := dic.NewContainer(di.App)
	if err != nil {
		log.Fatal("Error dic.NewContainer")
	}
	Application.Container = container
}
