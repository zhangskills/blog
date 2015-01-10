package controllers

import (
	"blog/app/models"
	"blog/app/utils"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/gosexy/to"
	"github.com/revel/revel"
	log "github.com/xiocode/glog"
	"sync"
	"time"
)

var (
	engine *xorm.Engine

	hotTags  []*models.KeyCount
	hotBlogs []*models.Blog

	blogViewMap  = make(map[int64]int64, 0)
	blogViewSync sync.Mutex

	autoRefreshHotTag   = models.NewAutoRefresh(5*time.Second, refreshHotTags)
	autoRefreshBlogView = models.NewAutoRefresh(5*time.Second, func() {
		var oldBlogViewMap map[int64]int64
		oldBlogViewMap, blogViewMap = blogViewMap, make(map[int64]int64, 0)
		for blogId, viewNum := range oldBlogViewMap {
			_, err := engine.Id(blogId).Incr("view_num", viewNum).Update(&models.Blog{})
			if err != nil {
				log.Errorln(err)
			}
		}
		refreshHotBlogs()
	})

	pageSize = 10
)

func refreshHotTags() {
	m, err := engine.Query("select a.name,count(1) from tag a,blog_tag b where a.id=b.tag_id group by tag_id order by count(1) desc limit 50")
	if err != nil {
		log.Errorln(err)
	} else {
		hotTags = hotTags[:0]
		for _, f := range m {
			hotTags = append(hotTags, &models.KeyCount{
				Key:   to.String(f["name"]),
				Count: to.Int64(f["count(1)"]),
			})
		}
	}
}

func refreshHotBlogs() {
	hotBlogs = hotBlogs[:0]
	err := engine.Desc("view_num").Limit(10).Find(&hotBlogs)
	if err != nil {
		log.Errorln(err)
	}
}

func addBlogView(blogId int64) {
	blogViewSync.Lock()
	defer blogViewSync.Unlock()

	if viewNum, ok := blogViewMap[blogId]; ok {
		blogViewMap[blogId] = viewNum + 1
	} else {
		blogViewMap[blogId] = 1
	}
	autoRefreshBlogView.SetStatus(true)
}

func getStart(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * pageSize
}

func init() {
	flag.Set("stderrthreshold", "INFO")
	flag.Set("logtostderr", "true")

	revel.TemplateFuncs["showDate"] = func(t time.Time) string {
		return t.Format("2006-01-02")
	}
	revel.TemplateFuncs["showDateTime"] = func(timeUnix int64) string {
		return time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	}

	revel.TemplateFuncs["pages"] = func(count int64) (pages []int) {
		pageCount := int(count-1)/pageSize + 2
		for i := 1; i < pageCount; i++ {
			pages = append(pages, i)
		}
		return
	}
	revel.TemplateFuncs["substr"] = func(s string, num int) string {
		str := utils.SubstrByByte(s, num*3)
		if len(str) < len(s) {
			return str + "..."
		}
		return str
	}

	revel.OnAppStart(func() {
		//初始化数据库
		driverName := revel.Config.StringDefault("db.driverName", "mysql")
		dataSourceName := revel.Config.StringDefault("db.dataSourceName", "")

		var err error
		engine, err = xorm.NewEngine(driverName, dataSourceName)
		if err != nil {
			log.Error(err)
			return
		}

		engine.ShowErr = revel.Config.BoolDefault("xorm.showErr", true)
		engine.ShowSQL = revel.Config.BoolDefault("xorm.showSQL", true)

		engine.CreateTables(&models.Blog{}, &models.Tag{}, &models.BlogTag{})

		cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
		engine.SetDefaultCacher(cacher)

		refreshHotTags()
		refreshHotBlogs()
	})

}
