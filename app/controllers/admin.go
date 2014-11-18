package controllers

import (
	"blog/app/models"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	log "github.com/xiocode/glog"
	"strings"
)

type Admin struct {
	Base
}

func (a Admin) Blogs(page int) revel.Result {
	start := getStart(page, pageSize)

	var blogs []*models.Blog
	err := engine.Desc("id").Limit(pageSize, start).Find(&blogs)
	if err != nil && err != gorm.RecordNotFound {
		log.Errorln(err)
		return a.RenderError(err)
	}
	count, err := engine.Count(&models.Blog{})
	if err != nil {
		log.Errorln(err)
		return a.RenderError(err)
	}
	return a.Render(blogs, count, page)
}

func (a *Admin) GetBlog(id int64) revel.Result {
	var blog models.Blog

	has, err := engine.Id(id).Get(&blog)
	if !has || err != nil {
		log.Error(err)
		return a.sendErrJson(err.Error())
	}
	engine.Table("tag").Join("left", "blog_tag", "tag.id=blog_tag.tag_id").Where("blog_tag.blog_id=?", id).Find(&blog.Tags)
	return a.sendOkJson(blog)
}

func (a Admin) DelBlog(id int64) revel.Result {
	num, err := engine.Delete(&models.Blog{Id: id})
	if num < 1 || err != nil {
		log.Error(err)
		return a.sendErrJson(err.Error())
	}
	engine.Delete(&models.BlogTag{BlogId: id})
	autoRefreshHotTag.SetStatus(true)
	autoRefreshBlogView.SetStatus(true)
	return a.sendOkJson("ok")
}

func (a Admin) SaveBlog(blog models.Blog, tagNames string) revel.Result {
	a.Validation.Required(blog.Title)
	a.Validation.Required(blog.Content)
	if a.Validation.HasErrors() {
		log.Errorln(a.Validation.ErrorMap())
		a.Validation.Keep()
		a.FlashParams()
		return a.sendErrJson("")
	}
	var num int64
	var err error
	if blog.Id > 0 {
		num, err = engine.Update(&blog)
	} else {
		num, err = engine.Insert(&blog)
	}
	if num > 0 && err == nil {
		//保存标签
		for _, tagName := range strings.Split(tagNames, ",") {
			var tag models.Tag
			has, err := engine.Where("name=?", tagName).Get(&tag)
			if !has || err != nil {
				tag = models.Tag{Name: tagName}
				engine.Insert(&tag)
			}
			blog.Tags = append(blog.Tags, &tag)
			engine.Delete(&models.BlogTag{BlogId: blog.Id, TagId: tag.Id})
			engine.Insert(&models.BlogTag{BlogId: blog.Id, TagId: tag.Id})
		}
	} else {
		log.Errorln(err)
		return a.sendErrJson(err.Error())
	}
	autoRefreshHotTag.SetStatus(true)
	return a.sendOkJson("")
}

func (a Admin) CheckUser() revel.Result {
	a.Session.Id()
	if a.Session["login"] != "1" {
		return a.Redirect(User.ToLogin)
	}
	return nil
}

func init() {
	revel.InterceptMethod(Admin.CheckUser, revel.BEFORE)
}
