package controllers

import (
	"blog/app/models"
	"blog/app/utils"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
	"github.com/revel/revel"
	"github.com/russross/blackfriday"
	log "github.com/xiocode/glog"
	"regexp"
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

func (a *Api) BlogList(tag string, page int) revel.Result {
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

func (a *Api) TagBlogList(tagName string, page int) revel.Result {
	var tag models.Tag
	has, err := engine.Where("name=?", tagName).Get(&tag)
	if !has {
		log.Error("标签不存在")
		return a.SendJson(errors.New("标签不存在"), "")
	} else if err != nil {
		log.Errorln(err)
		return a.SendJson(err, "")
	}
	count, err := engine.Count(&tag)

	var blogs []models.Blog
	start := getStart(page, pageSize)
	err = engine.Table("blog").Join("left", "blog_tag", "blog.id=blog_tag.blog_id").Where("tag_id=?", tag.Id).Limit(pageSize, start).Find(&blogs)
	if err != nil {
		log.Errorln(err)
		return a.SendJson(err, "")
	}

	return a.SendJson(err, map[string]interface{}{
		"blogs": blogs,
		"count": count,
	})
}

func (a *Api) TagNames() revel.Result {
	var tags []models.Tag
	err := engine.Find(&tags)
	if err != nil {
		log.Error(err)
		return a.SendJson(err, "")
	}
	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return a.SendJson(err, tagNames)
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

func (a *Api) HotTags() revel.Result {
	return a.SendJson(nil, hotTags)
}

func (a *Api) HotBlogs() revel.Result {
	return a.SendJson(nil, hotBlogs)
}

func init() {
	revel.OnAppStart(func() {
		conf.ACCESS_KEY = revel.Config.StringDefault("qiniu.access_key", "")
		conf.SECRET_KEY = revel.Config.StringDefault("qiniu.secret_key", "")
	})
}
func (a *Api) UploadToken(fileName string) revel.Result {
	r := regexp.MustCompile(`\.(jpe?g|png|bmp|gif)$`)
	suf := r.FindString(strings.ToLower(fileName))
	if suf == "" {
		return a.SendJson(errors.New("请选择图片格式"), "")
	}
	fileName = utils.NewFileName() + suf
	putPolicy := rs.PutPolicy{
		Scope: "ww-blog:" + fileName,
	}
	token := putPolicy.Token(nil)
	return a.SendJson(nil, map[string]interface{}{
		"token": token,
		"key":   fileName,
	})
}
