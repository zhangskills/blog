package controllers

import (
	"github.com/revel/revel"
)

type Base struct {
	*revel.Controller
}

func (b Base) sendOkJson(msg interface{}) revel.Result {
	return b.RenderJson(map[string]interface{}{
		"msg": msg,
		"err": false,
	})
}

func (b Base) sendErrJson(msg string) revel.Result {
	return b.RenderJson(map[string]interface{}{
		"msg":    msg,
		"err":    true,
		"errMap": b.Validation.ErrorMap(),
	})
}
