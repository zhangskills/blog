package controllers

import (
	"blog/app/models"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"github.com/russross/blackfriday"
	log "github.com/xiocode/glog"
	"strings"
)

type Blog struct {
	Base
}

func (b Blog) List(page int) revel.Result {
	var blogs []*models.Blog
	start := getStart(page, pageSize)
	err := engine.Desc("id").Limit(pageSize, start).Find(&blogs)
	if err != nil && err != gorm.RecordNotFound {
		log.Errorln(err)
		return b.RenderError(err)
	}
	count, err := engine.Count(&models.Blog{})
	if err != nil {
		log.Errorln(err)
		return b.RenderError(err)
	}

	pageName, nav := "全部", 1
	return b.Render(blogs, count, page, pageName, nav)
}

func (b Blog) Show(id int64) revel.Result {
	var blog models.Blog

	has, err := engine.Id(id).Get(&blog)
	if !has || err != nil {
		log.Error(err)
		return b.sendErrJson(err.Error())
	}
	engine.Table("tag").Join("left", "blog_tag", "tag.id=blog_tag.tag_id").Where("blog_tag.blog_id=?", id).Find(&blog.Tags)

	content := string(blackfriday.MarkdownCommon([]byte(blog.Content)))
	blog.Content = strings.Replace(content, "<pre><code", "<pre class=\"prettyprint\"><code", -1)
	go addBlogView(id)
	return b.Render(blog)
}
