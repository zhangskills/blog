package controllers

import (
	"github.com/revel/revel"
)

type User struct {
	*revel.Controller
}

func (u User) ToLogin() revel.Result {
	return u.Render()
}

func (u *User) Login(name, password string) revel.Result {
	if name == revel.Config.StringDefault("user.name", "admin") && password == revel.Config.StringDefault("user.password", "admin123") {
		u.Session.Id()
		u.Session["login"] = "1"
		return u.Redirect("/admin/blogs/1")
	}
	u.Flash.Error("用户名或密码不正确")
	u.FlashParams()
	return u.Redirect(User.ToLogin)
}
