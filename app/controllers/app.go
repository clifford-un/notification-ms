package controllers

import (
	"github.com/revel/revel"
)

// App
type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}
