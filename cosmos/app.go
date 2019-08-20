package cosmos

import (
	"github.com/cosmos/cosmos-sdk/types/module"
)

type App struct {
	modules []module.AppModule
}

func New() *App {
	return &App{
		modules: []module.AppModule{},
	}
}

func (a *App) RegisterModule(module module.AppModule) {
	a.modules = append(a.modules, module)
}

func (a *App) Start() {
}
