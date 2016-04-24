package controller

import (
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/unrolled/render"
	"myblog/src/lib"
	"myblog/src/models"

	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//后台方法
func CategoryListAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})

	ctx["title"] = "博客类型列表"

	categoryList := []models.Category{}
	dbMap.Select(&categoryList, "select * from category order by seq asc,id asc")

	ctx["categoryList"] = categoryList

	// 把必需公共数据放入context ctx
	ctx = requiredCtxProcess(req, ctx)

	ctx["flashes"] = GetFlashMessages(req, w)

	r.HTML(w, http.StatusOK, "admin/category_list", ctx, render.HTMLOptions{Layout: "layout"})
}

func CategoryAddAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})

	if strings.ToUpper(req.Method) == "POST" {
		category := models.Category{}
		errs := binding.Bind(req, &category)
		if errs.Handle(w) {
			return
		}

		//err := bind.Request(req).Field(&category, "category")

		err := dbMap.Insert(&category)
		if err == nil {
			SetFlashMessages(req, w, "新增成功！")
			http.Redirect(w, req, "/admin/category/list", http.StatusFound)
			return
		} else {
			SetFlashMessages(req, w, "保存数据出错！")
			http.Redirect(w, req, "/admin/category/add", http.StatusFound)
			return
		}
	} else {
		ctx["title"] = "新增博客类型"

		ctx["action"] = "add"
		ctx["category"] = models.Category{}

		// 把必需公共数据放入context ctx
		ctx = requiredCtxProcess(req, ctx)

		ctx["flashes"] = GetFlashMessages(req, w)

		r.HTML(w, http.StatusOK, "admin/category_add", ctx, render.HTMLOptions{Layout: "layout"})
		return
	}
}

func CategoryUpdateAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})
	params := mux.Vars(req)
	id := params["id"]
	if strings.ToUpper(req.Method) == "POST" {
		category := models.Category{}
		errs := binding.Bind(req, &category)
		if errs.Handle(w) {
			return
		}

		_, err := dbMap.Update(&category)
		if err == nil {
			SetFlashMessages(req, w, "编辑成功！")
			http.Redirect(w, req, "/admin/category/list", http.StatusFound)
			return
		} else {
			SetFlashMessages(req, w, "保存数据出错！")
			http.Redirect(w, req, "/admin/category/"+id+"/edit", http.StatusFound)
			return
		}
	} else {
		ctx["title"] = "修改博客类型"

		ctx["action"] = "update"

		category := models.Category{}
		dbMap.SelectOne(&category, "select * from category where id=$1", id)

		ctx["category"] = category

		// 把必需公共数据放入context ctx
		ctx = requiredCtxProcess(req, ctx)

		ctx["flashes"] = GetFlashMessages(req, w)

		r.HTML(w, http.StatusOK, "admin/category_add", ctx, render.HTMLOptions{Layout: "layout"})
		return
	}
}

func CategoryDelAdmin(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	count, _ := dbMap.SelectInt("select count(1) from article where categoryid=$1", id)

	if count > 0 {
		SetFlashMessages(req, w, "此分类有所属文章，需删除此类的文章后才能操作！")
	} else {
		dbMap.Exec("delete from category where id=$1", id)
		SetFlashMessages(req, w, "删除成功！")
	}

	http.Redirect(w, req, "/admin/category/list", http.StatusFound)
}

func ArticleListAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})

	ctx["title"] = "博客列表"
	pageStr := req.FormValue("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, _ := strconv.Atoi(pageStr)

	q := req.FormValue("q")
	ctx["q2"] = q

	recordCount := int64(0)
	if "" != q {
		q = "%" + q + "%"
		recordCount, _ = dbMap.SelectInt("select count(*) count from article where title like $1 or content like $2", q, q)
	} else {
		recordCount, _ = dbMap.SelectInt("select count(*) count from article")
	}

	pageSize := int64(20)
	paging := lib.NewPager(int64(page), pageSize, int64(recordCount))
	offset := (paging.CurrentPage() - 1) * pageSize

	articleList := []models.Article{}

	if "" != q {
		dbMap.Select(&articleList, "select * from article where title like $3 or content like $4 order by id desc LIMIT $1 OFFSET $2", pageSize, offset, q, q)
	} else {
		dbMap.Select(&articleList, "select * from article order by id desc LIMIT $1 OFFSET $2", pageSize, offset)
	}

	ctx["articleList"] = articleList
	ctx["pager"] = paging

	// 把必需公共数据放入context ctx
	ctx = requiredCtxProcess(req, ctx)

	ctx["flashes"] = GetFlashMessages(req, w)

	r.HTML(w, http.StatusOK, "admin/article_list", ctx, render.HTMLOptions{Layout: "layout"})
}

func ArticleAddAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})
	// account := GetSession(req, LOGIN_ADMIN).(*models.Account)

	if strings.ToUpper(req.Method) == "POST" {
		article := models.Article{}
		errs := binding.Bind(req, &article)
		if errs.Handle(w) {
			return
		}

		article.Tags = processTags(article.Tags)

		err := dbMap.Insert(&article)
		if err == nil {
			SetFlashMessages(req, w, "新增成功！")
			http.Redirect(w, req, "/admin/article/list", http.StatusFound)
			return
		} else {
			showErr(err, "新增博客出错")
			SetFlashMessages(req, w, "保存数据出错！")
			http.Redirect(w, req, "/admin/article/add", http.StatusFound)
			return
		}
	} else {
		ctx["title"] = "新增博客"
		ctx["action"] = "add"

		ctx["article"] = models.Article{}

		category := []models.Category{}
		dbMap.Select(&category, "select * from category where status=true order by seq asc,id asc")

		ctx["category"] = category

		// 把必需公共数据放入context ctx
		ctx = requiredCtxProcess(req, ctx)

		ctx["flashes"] = GetFlashMessages(req, w)

		r.HTML(w, http.StatusOK, "admin/article_add", ctx, render.HTMLOptions{Layout: "layout"})
		return
	}
}

func ArticleUpdateAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})

	if strings.ToUpper(req.Method) == "POST" {
		article := models.Article{}
		errs := binding.Bind(req, &article)
		if errs.Handle(w) {
			return
		}

		article.Tags = processTags(article.Tags)

		_, err := dbMap.Update(&article)
		if err == nil {
			SetFlashMessages(req, w, "编辑成功！")
			http.Redirect(w, req, "/admin/article/list", http.StatusFound)
			return
		} else {
			showErr(err, "编辑博客出错")
			SetFlashMessages(req, w, "编辑出错！")
			http.Redirect(w, req, "/admin/article/add", http.StatusFound)
			return
		}
	} else {
		params := mux.Vars(req)

		ctx["title"] = "新增博客"

		ctx["action"] = "update"

		article := models.Article{}
		dbMap.SelectOne(&article, "select * from article where id=$1", params["id"])

		ctx["article"] = article

		category := []models.Category{}
		dbMap.Select(&category, "select * from category where status=true order by seq asc,id asc")

		ctx["category"] = category

		// 把必需公共数据放入context ctx
		ctx = requiredCtxProcess(req, ctx)

		ctx["flashes"] = GetFlashMessages(req, w)

		r.HTML(w, http.StatusOK, "admin/article_add", ctx, render.HTMLOptions{Layout: "layout"})
		return
	}
}

func ArticleDelAdmin(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	dbMap.Exec("delete from comment where articleid=$1", id)
	dbMap.Exec("delete from article where id=$1", id)
	SetFlashMessages(req, w, "删除成功！")
	http.Redirect(w, req, "/admin/article/list", http.StatusFound)
}

func AboutMeEditAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})

	account := models.Account{}
	if strings.ToUpper(req.Method) == "POST" {
		errs := binding.Bind(req, &account)
		if errs.Handle(w) {
			return
		}

		if account.Password == "" {
			pwd, _ := dbMap.SelectStr("select password from account order by id asc limit 1")
			account.Password = pwd
		} else {
			account.Password = lib.MD5(account.Password)
		}

		//上传图片
		req.ParseMultipartForm(32 << 20)
		hatHeadFile, _, _ := req.FormFile("hathead")
		if hatHeadFile != nil {
			defer hatHeadFile.Close()
			account.Image = uploadHatHead(hatHeadFile)
		}

		count, _ := dbMap.Update(&account)
		if count > 0 {
			SetFlashMessages(req, w, "编辑成功！")
		} else {
			SetFlashMessages(req, w, "编辑失败！")
		}

		http.Redirect(w, req, "/admin/user/edit", http.StatusFound)
		return
	} else {
		ctx["title"] = "编辑用户信息"

		// 把必需公共数据放入context ctx
		ctx = requiredCtxProcess(req, ctx)

		dbMap.SelectOne(&account, "select * from account where status=true limit 1")

		ctx["account"] = account

		ctx["flashes"] = GetFlashMessages(req, w)

		r.HTML(w, http.StatusOK, "admin/user_add", ctx, render.HTMLOptions{Layout: "layout"})
	}
}

func BackUpBlog(w http.ResponseWriter, req *http.Request) {
	media_path, _ := settings.String("admin", "media_path")
	var dir = "backup"

	if lib.IsDirExists(filepath.Join(media_path, dir)) == false {
		os.MkdirAll(filepath.Join(media_path, dir), os.ModePerm)
	}

	fileName := filepath.Join(media_path, dir, "myblog.xml")

	content := `<?xml version="1.0" encoding="utf-8"?>` + "\n"
	content += `<blogs>` + "\n"

	articleList := []models.Article{}
	dbMap.Select(&articleList, "select title,content,origin,created from article order by id asc")

	size := len(articleList)
	for i := 0; i < size; i++ {
		content += "<blog>" + "\n"
		content += "<title>" + articleList[i].Title + "</title>" + "\n"
		content += "<content><![CDATA[" + articleList[i].Content + "]]></content>" + "\n"
		content += "<tags>" + articleList[i].Tags + "</tags>" + "\n"
		content += "<origin>" + articleList[i].Origin + "</origin>" + "\n"
		content += "<created>" + articleList[i].Created.Format(TIMEFORMAT) + "</created>" + "\n"
		content += "</blog>" + "\n"
	}
	content += `</blogs>` + "\n"

	lib.WriteStr2File(fileName, content)

	r.Data(w, http.StatusOK, []byte(content))
}

func SettingsEditAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})
	settings := models.Settings{}
	if strings.ToUpper(req.Method) == "POST" {
		errs := binding.Bind(req, &settings)
		if errs.Handle(w) {
			return
		}

		//上传图片
		req.ParseMultipartForm(32 << 20)
		faviconFile, faviconHeadFile, _ := req.FormFile("faviconFile")
		if faviconHeadFile != nil {
			defer faviconFile.Close()
			tmp, err := uploadFile(faviconHeadFile, "file")
			if err == nil {
				settings.Favicon = tmp
			}

		}

		donateFile, _, _ := req.FormFile("donateFile")
		if donateFile != nil {
			defer donateFile.Close()
			settings.Donate = uploadDonatePng(donateFile)
		}

		count, _ := dbMap.Update(&settings)
		if count > 0 {
			SetFlashMessages(req, w, "编辑成功！")
		} else {
			SetFlashMessages(req, w, "编辑失败！")
		}

		http.Redirect(w, req, "/admin/settings/edit", http.StatusFound)
		return
	} else {
		ctx["title"] = "博客设置"
		dbMap.SelectOne(&settings, "select * from settings order by id asc limit 1")
		if settings.Id <= 0 {
			settings.Favicon = "/static/images/favicon.ico"
			settings.BlogName = "我的博客"
			settings.BlogIntro = "myblog 是用 Go 语言开发的一款个人博客，界面简洁单纯，推崇极客思想，博客功能体系也扁平直接，应该是程序员自助撰写博客的首选方案！"
			settings.EmailSmtp = "smtp.163.com"
			settings.EmailSmtpPort = 465
			settings.EmailSender = "youremail@163.com"
			settings.EmailPwd = "123456"
			settings.Org = "myblog"
			settings.ICPNO = ""
			settings.Donate = "/static/images/donateforme.png"
			dbMap.Insert(&settings)
		}

		// 把必需公共数据放入context ctx
		ctx = requiredCtxProcess(req, ctx)

		ctx["settings"] = settings

		ctx["flashes"] = GetFlashMessages(req, w)

		r.HTML(w, http.StatusOK, "admin/settings", ctx, render.HTMLOptions{Layout: "layout"})
	}
}

func AdvertEditAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})

	advert := models.Advert{}
	if strings.ToUpper(req.Method) == "POST" {
		errs := binding.Bind(req, &advert)
		if errs.Handle(w) {
			return
		}

		count, _ := dbMap.Update(&advert)
		if count > 0 {
			SetFlashMessages(req, w, "编辑成功！")
		} else {
			SetFlashMessages(req, w, "编辑失败！")
		}

		http.Redirect(w, req, "/admin/advert/edit", http.StatusFound)
		return
	} else {
		ctx["title"] = "广告设置"

		dbMap.SelectOne(&advert, "select * from advert order by id asc limit 1")
		if advert.Id <= 0 {
			dbMap.Insert(&advert)
		}

		// 把必需公共数据放入context ctx
		ctx = requiredCtxProcess(req, ctx)

		ctx["advert"] = advert

		ctx["flashes"] = GetFlashMessages(req, w)

		r.HTML(w, http.StatusOK, "admin/advert_add", ctx, render.HTMLOptions{Layout: "layout"})
	}
}

func FriendLinkListAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})

	ctx["title"] = "友情链接列表"

	// 把必需公共数据放入context ctx
	ctx = requiredCtxProcess(req, ctx)

	friendLinkList := []models.FriendLink{}
	dbMap.Select(&friendLinkList, "select * from friendlink order by status asc,id asc")

	ctx["friendLinkList"] = friendLinkList

	ctx["flashes"] = GetFlashMessages(req, w)

	r.HTML(w, http.StatusOK, "admin/friendlink_list", ctx, render.HTMLOptions{Layout: "layout"})
}

func FriendLinkAddAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})

	if strings.ToUpper(req.Method) == "POST" {
		friendLink := models.FriendLink{}
		errs := binding.Bind(req, &friendLink)
		if errs.Handle(w) {
			return
		}

		err := dbMap.Insert(&friendLink)
		if err == nil {
			SetFlashMessages(req, w, "新增成功！")
			http.Redirect(w, req, "/admin/friendlink/list", http.StatusFound)
			return
		} else {
			SetFlashMessages(req, w, "保存数据出错！")
			http.Redirect(w, req, "/admin/friendlink/add", http.StatusFound)
			return
		}
	} else {
		// 把必需公共数据放入context ctx
		ctx = requiredCtxProcess(req, ctx)

		ctx["title"] = "添加友情链接"

		ctx["action"] = "add"
		ctx["friendlink"] = models.FriendLink{}

		ctx["flashes"] = GetFlashMessages(req, w)

		r.HTML(w, http.StatusOK, "admin/friendlink_add", ctx, render.HTMLOptions{Layout: "layout"})
		return
	}
}

func FriendLinkEditAdmin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})

	friendLink := models.FriendLink{}
	if strings.ToUpper(req.Method) == "POST" {
		errs := binding.Bind(req, &friendLink)
		if errs.Handle(w) {
			return
		}

		_, err := dbMap.Update(&friendLink)
		if err == nil {
			SetFlashMessages(req, w, "新增成功！")
			http.Redirect(w, req, "/admin/friendlink/list", http.StatusFound)
			return
		} else {
			SetFlashMessages(req, w, "保存数据出错！")
			http.Redirect(w, req, "/admin/friendlink/add", http.StatusFound)
			return
		}
	} else {
		params := mux.Vars(req)
		// 把必需公共数据放入context ctx
		ctx = requiredCtxProcess(req, ctx)

		ctx["title"] = "编辑友情链接"

		ctx["action"] = "update"

		dbMap.SelectOne(&friendLink, "select * from friendlink where id=$1", params["id"])

		ctx["friendlink"] = friendLink

		ctx["flashes"] = GetFlashMessages(req, w)

		r.HTML(w, http.StatusOK, "admin/friendlink_add", ctx, render.HTMLOptions{Layout: "layout"})
		return
	}
}

func FriendLinkPassAdmin(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	action := req.FormValue("action")

	if id != "" && action != "" {
		status := false
		if action == "1" {
			status = true
		}
		dbMap.Exec("update friendlink set status=$1 where id=$2", status, id)
		SetFlashMessages(req, w, "操作成功！")
		http.Redirect(w, req, "/admin/friendlink/list", http.StatusFound)
		return
	}

	SetFlashMessages(req, w, "操作失败！")
	http.Redirect(w, req, "/admin/friendlink/list", http.StatusFound)
}

func FriendLinkDelAdmin(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	if id != "" {
		dbMap.Exec("delete from friendlink where id=$1", id)
		SetFlashMessages(req, w, "操作成功！")
		http.Redirect(w, req, "/admin/friendlink/list", http.StatusFound)
		return
	}

	SetFlashMessages(req, w, "操作失败！")
	http.Redirect(w, req, "/admin/friendlink/list", http.StatusFound)
}

func FlushDataAdmin(w http.ResponseWriter, req *http.Request) {
	RunAllStaticstics()
	cacheObj.Flush()
	ClearAllSession(req, w)
	http.Redirect(w, req, "/", http.StatusFound)
}
