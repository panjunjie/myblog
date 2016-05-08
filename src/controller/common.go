package controller

import (
	"github.com/go-gomail/gomail"
	"github.com/pmylund/go-cache"
	"myblog/src/lib"
	"myblog/src/models"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//把必需公共数据放入ctx中
func requiredCtxProcess(req *http.Request, ctx map[string]interface{}) map[string]interface{} {
	if ctx == nil {
		ctx = make(map[string]interface{})
	}

	ctx["logined"] = _loginProcess(req)

	// 博客总数和评论总数
	ctx["articleTotal"] = getArticleTotal()
	ctx["commentTotal"] = getCommentTotal()
	ctx["readArticleTotal"] = getArticleReadTotal()

	//博客设置
	ctx["settings"] = getSettings()

	// 用户信息
	ctx["account"] = getAccount()

	//广告记录
	ctx["advert"] = getAdvert()

	ctx["timefmt"] = TIMEFORMAT

	return ctx
}

//把全部公共数据放入ctx中
func commonContextProcess(req *http.Request, ctx map[string]interface{}) map[string]interface{} {
	if ctx == nil {
		ctx = make(map[string]interface{})
	}

	ctx["logined"] = _loginProcess(req)

	// 博客总数和评论总数
	ctx["articleTotal"] = getArticleTotal()
	ctx["commentTotal"] = getCommentTotal()
	ctx["readArticleTotal"] = getArticleReadTotal()

	//博客设置
	ctx["settings"] = getSettings()

	// 用户信息
	ctx["account"] = getAccount()

	// 分类
	ctx["categoryList"] = getCategoryList()

	//归档博客
	ctx["archiveList"] = getArchiveList()

	// 最新评论
	ctx["newComment"] = getNewCommentList()

	// 热门博客
	ctx["hotArticles"] = getHotArticleList()

	//最近浏览的博客
	ctx["recentView"] = getRecentView()

	// 友情链接
	ctx["friendLinkList"] = getFriendLinkList()

	//广告记录
	ctx["advert"] = getAdvert()

	//构建语言
	ctx["langVersion"] = strings.Replace(runtime.Version(), "go", "Go ", -1)

	return ctx
}

func loginProcess(req *http.Request, ctx map[string]interface{}) map[string]interface{} {
	if ctx == nil {
		ctx = make(map[string]interface{})
	}
	ctx["logined"] = _loginProcess(req)
	return ctx
}

func _loginProcess(req *http.Request) bool {
	logined := false
	if GetSession(req, LOGIN_ADMIN) != nil {
		logined = true
	}
	return logined
}

// 获取用户信息
func getAccount() *models.Account {
	obj, found := cacheObj.Get("getAccount")
	if !found || obj == nil {
		account := models.Account{}
		dbMap.SelectOne(&account, "select * from account where status=true order by id asc limit 1")
		cacheObj.Set("getAccount", &account, cache.DefaultExpiration)
		return &account
	}
	return obj.(*models.Account)
}

func getSettings() *models.Settings {
	obj, found := cacheObj.Get("getSettings")
	if !found || obj == nil {
		settings := models.Settings{}
		dbMap.SelectOne(&settings, "select * from settings order by id asc limit 1")
		cacheObj.Set("getSettings", &settings, cache.DefaultExpiration)
		return &settings
	}
	return obj.(*models.Settings)
}

func getAdvert() *models.Advert {
	obj, found := cacheObj.Get("getAdvert")
	if !found || obj == nil {
		advert := models.Advert{}
		dbMap.SelectOne(&advert, "select * from advert order by id asc limit 1")
		cacheObj.Set("getAdvert", &advert, cache.DefaultExpiration)
		return &advert
	}
	return obj.(*models.Advert)
}

//statistics category 1:统计博客数评论数访问数点赞数；2：归档博客；3：热门博客；4：最新点评博客；5：热评博客;6：点赞博客排名

func getTotalCount() [4]int64 {
	obj, found := cacheObj.Get("getTotalCount")
	if !found || obj == nil {
		statistics := models.Statistics{}
		err := dbMap.SelectOne(&statistics, "select id,param1,param2,param3,param4 from statistics where category=$1 limit 1", 1)
		ret := [4]int64{0, 0, 0, 0}
		if err == nil && statistics.Id > 0 {
			ret[0], _ = strconv.ParseInt(statistics.Param1, 10, 64)
			ret[1], _ = strconv.ParseInt(statistics.Param2, 10, 64)
			ret[2], _ = strconv.ParseInt(statistics.Param3, 10, 64)
			ret[3], _ = strconv.ParseInt(statistics.Param4, 10, 64)
		}
		cacheObj.Set("getTotalCount", ret, cache.DefaultExpiration)
		return ret
	}
	return obj.([4]int64)
}

// 获取文章总数
func getArticleTotal() int64 {
	return getTotalCount()[0]
}

// 获取总评论数
func getCommentTotal() int64 {
	return getTotalCount()[1]
}

func getArticleReadTotal() int64 {
	return getTotalCount()[2]
}

//获取分类
func getCategoryList() *[]models.Category {
	obj, found := cacheObj.Get("getCategoryList")
	if !found || obj == nil {
		categoryList := []models.Category{}
		dbMap.Select(&categoryList, "select id,name,intro,itemcount from category where status=true order by seq asc")
		cacheObj.Set("getCategoryList", &categoryList, cache.DefaultExpiration)
		return &categoryList
	}
	return obj.(*[]models.Category)
}

//博客归档
func getArchiveList() *[]models.Statistics {
	obj, found := cacheObj.Get("getArchiveList")
	if !found || obj == nil {
		list := []models.Statistics{}
		dbMap.Select(&list, "select id,param1,param2,param3,param4 from statistics where category=$1", 2)
		cacheObj.Set("getArchiveList", &list, cache.DefaultExpiration)
		return &list
	}
	return obj.(*[]models.Statistics)
}

// 获取阅读量最多的博客
func getHotArticleList() *[]models.Statistics {
	obj, found := cacheObj.Get("getHotArticleList")
	if !found || obj == nil {
		list := []models.Statistics{}
		dbMap.Select(&list, "select id,title1,param1,param2 from statistics where category=$1", 3)
		cacheObj.Set("getHotArticleList", &list, cache.DefaultExpiration)
		return &list
	}
	return obj.(*[]models.Statistics)
}

func getRecentView() *[]models.Article {
	obj, found := cacheObj.Get("getRecentView")
	if !found || obj == nil {
		list := []models.Article{}
		dbMap.Select(&list, "select id,title from article order by updated desc limit 10")
		cacheObj.Set("getRecentView", &list, cache.DefaultExpiration)
		return &list
	}
	return obj.(*[]models.Article)
}

// 获取最新的评论
func getNewCommentList() *[]models.Statistics {
	obj, found := cacheObj.Get("getNewCommentList")
	if !found || obj == nil {
		list := []models.Statistics{}
		dbMap.Select(&list, "select id,title1,title2,param1,param2,param3,param4 from statistics where category=$1", 4)
		cacheObj.Set("getNewCommentList", &list, cache.DefaultExpiration)
		return &list
	}
	return obj.(*[]models.Statistics)
}

// 获取顶最多的博客
func getUpArticleList() *[]models.Statistics {
	obj, found := cacheObj.Get("getUpArticleList")
	if !found || obj == nil {
		list := []models.Statistics{}
		dbMap.Select(&list, "select id,param1,param2,param3,param4 from statistics where category=$1", 3)
		cacheObj.Set("getUpArticleList", &list, cache.DefaultExpiration)
		return &list
	}
	return obj.(*[]models.Statistics)
}

// 友情链接
func getFriendLinkList() *[]models.FriendLink {
	obj, found := cacheObj.Get("getFriendLinkList")
	if !found || obj == nil {
		friendLinkList := []models.FriendLink{}
		dbMap.Select(&friendLinkList, "select name,link,intro from friendlink a where a.status=true order by created desc")
		cacheObj.Set("getFriendLinkList", &friendLinkList, cache.DefaultExpiration)
		return &friendLinkList
	}
	return obj.(*[]models.FriendLink)
}

func processTags(tags string) string {
	if tags != "" {
		tags = strings.Replace(tags, "，", ",", -1)
		tags = strings.Replace(tags, "、", ",", -1)
		tags = strings.Replace(tags, "；", ",", -1)
		tags = strings.Replace(tags, ";", ",", -1)
		tags = strings.Replace(tags, "|", ",", -1)
		tags = strings.Replace(tags, "&", ",", -1)
	}
	return tags
}

func sendEmail(from, to, subject, body string) {

	settings := getSettings()

	email_server := settings.EmailSmtp
	email_port := settings.EmailSmtpPort
	sender := settings.EmailSender
	sender_pwd := settings.EmailPwd

	if email_server == "" || sender == "" || sender_pwd == "" || email_port <= 0 {
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewPlainDialer(email_server, email_port, sender, sender_pwd)

	if err := d.DialAndSend(m); err != nil {
		showErr(err, "发送邮件出错:"+email_server+"|"+sender+"|"+sender_pwd)
	}
}

func addArticleRead(articleId int64) {
	_, err := dbMap.Exec("update article set read=read+1,updated=$2 where id=$1", articleId, time.Now())
	lib.ShowErr(err, "增加博客阅读数出错：")
}
