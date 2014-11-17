package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (a App) About() revel.Result {
	nav := 3
	return a.Render(nav)
}
