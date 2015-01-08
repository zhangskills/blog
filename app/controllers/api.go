package controllers

import (
	"blog/app/models"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"github.com/russross/blackfriday"
	log "github.com/xiocode/glog"
	"strings"
)

type Api struct {
	*revel.Controller
}

func (a *Api) Index() revel.Result {
	return a.RenderTemplate("app/index.html")
}

func (a *Api) SendJson(err error, val interface{}) revel.Result {
	return a.RenderJson(map[string]interface{}{
		"err": err,
		"val": val,
	})
}

func (a *Api) BlogList(page int) revel.Result {
	var blogs []*models.Blog
	start := getStart(page, pageSize)
	err := engine.Desc("id").Limit(pageSize, start).Find(&blogs)
	if err != nil && err != gorm.RecordNotFound {
		log.Errorln(err)
		return a.SendJson(err, "`")
	}
	count, err := engine.Count(&models.Blog{})
	if err != nil {
		log.Errorln(err)
		return a.SendJson(err, "1")
	}

	return a.SendJson(err, map[string]interface{}{
		"blogs": blogs,
		"count": count,
	})
}

func (a *Api) BlogShow(id int64) revel.Result {
	var blog models.Blog

	has, err := engine.Id(id).Get(&blog)
	if !has || err != nil {
		log.Error(err)
		return a.SendJson(err, "")
	}
	engine.Table("tag").Join("left", "blog_tag", "tag.id=blog_tag.tag_id").Where("blog_tag.blog_id=?", id).Find(&blog.Tags)

	content := string(blackfriday.MarkdownCommon([]byte(blog.Content)))
	blog.Content = strings.Replace(content, "<pre><code", "<pre class=\"prettyprint\"><code", -1)
	go addBlogView(id)

	return a.SendJson(err, blog)
}