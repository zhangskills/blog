package controllers

import (
	"blog/app/models"
	"errors"
	"github.com/revel/revel"
	log "github.com/xiocode/glog"
	"strings"
)

type AdminApi struct {
	*revel.Controller
}

func (a *AdminApi) SendJson(err error, val interface{}) revel.Result {
	res := map[string]interface{}{
		"val": val,
	}
	if err != nil {
		res["err"] = err.Error()
	}
	return a.RenderJson(res)
}

func (a *AdminApi) Index() revel.Result {
	return a.RenderTemplate("admin/index.html")
}

func (a *AdminApi) Login(name, password string) revel.Result {
	if name == revel.Config.StringDefault("user.name", "admin") && password == revel.Config.StringDefault("user.password", "admin123") {
		a.Session.Id()
		a.Session["login"] = "1"
		return a.SendJson(nil, "ok")
	}
	return a.SendJson(errors.New("用户名或密码不正确"), "err")
}

func (a *AdminApi) SaveBlog(blog models.Blog, tagNames string) revel.Result {
	a.Validation.Required(blog.Title)
	a.Validation.Required(blog.Content)
	if a.Validation.HasErrors() {
		log.Errorln(a.Validation.ErrorMap())
		return a.SendJson(errors.New("表单验证失败"), "")
	}
	var num int64
	var err error
	if blog.Id > 0 {
		num, err = engine.Id(blog.Id).Update(&blog)
	} else {
		num, err = engine.Insert(&blog)
	}
	if num > 0 && err == nil {
		//保存标签
		engine.Delete(&models.BlogTag{BlogId: blog.Id})
		for _, tagName := range strings.Split(tagNames, ",") {
			var tag models.Tag
			has, err := engine.Where("name=?", tagName).Get(&tag)
			if !has || err != nil {
				tag = models.Tag{Name: tagName}
				engine.Insert(&tag)
			}
			blog.Tags = append(blog.Tags, &tag)
			engine.Insert(&models.BlogTag{BlogId: blog.Id, TagId: tag.Id})
		}
	} else {
		log.Errorln(err)
		return a.SendJson(err, "")
	}
	autoRefreshHotTag.SetStatus(true)
	return a.SendJson(err, "ok")
}

func (a *AdminApi) DelBlog(id int64) revel.Result {
	num, err := engine.Delete(&models.Blog{Id: id})
	if num < 1 || err != nil {
		log.Error(err)
		return a.SendJson(err, "")
	}
	engine.Delete(&models.BlogTag{BlogId: id})
	autoRefreshHotTag.SetStatus(true)
	autoRefreshBlogView.SetStatus(true)

	return a.SendJson(err, "ok")
}
