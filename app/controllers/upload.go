package controllers

import (
	"gBlog/app/utils"
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
	"github.com/revel/revel"
	"regexp"
	"strings"
)

func init() {
	revel.OnAppStart(func() {
		conf.ACCESS_KEY = revel.Config.StringDefault("qiniu.access_key", "")
		conf.SECRET_KEY = revel.Config.StringDefault("qiniu.secret_key", "")
	})
}

type Upload struct {
	Base
}

func (u *Upload) UploadToken(fileName string) revel.Result {
	r := regexp.MustCompile(`\.(jpe?g|png|bmp|gif)$`)
	suf := r.FindString(strings.ToLower(fileName))
	if suf == "" {
		return u.sendErrJson("请选择图片格式")
	}
	fileName = utils.NewFileName() + suf
	putPolicy := rs.PutPolicy{
		Scope: "ww-blog:" + fileName,
	}
	token := putPolicy.Token(nil)
	return u.sendOkJson(map[string]interface{}{
		"token": token,
		"key":   fileName,
	})
}
