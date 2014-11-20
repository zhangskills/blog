package controllers

import (
	"blog/app/models"
	"github.com/revel/revel"
	log "github.com/xiocode/glog"
)

type Tag struct {
	*revel.Controller
}

func (t Tag) Cloud() revel.Result {
	nav := 2
	return t.Render(nav)
}

func (t Tag) BlogList(tagName string, page int) revel.Result {

	var tag models.Tag
	has, err := engine.Where("name=?", tagName).Get(&tag)
	if !has {
		log.Error("标签不存在")
		return t.RenderError(err)
	} else if err != nil {
		log.Errorln(err)
		return t.RenderError(err)
	}
	count, err := engine.Count(&tag)

	var blogs []models.Blog
	start := getStart(page, pageSize)
	err = engine.Table("blog").Join("left", "blog_tag", "blog.id=blog_tag.blog_id").Where("tag_id=?", tag.Id).Limit(pageSize, start).Find(&blogs)
	if err != nil {
		log.Errorln(err)
		return t.RenderError(err)
	}
	pageName := "标签：" + tagName

	t.Render(blogs, count, page, pageName)
	return t.RenderTemplate("blog/list.html")
}
