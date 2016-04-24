package controller

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"myblog/src/lib"
	"myblog/src/models"

	"net/http"
	"strconv"
	"strings"
	"time"
)

func ArticleList(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})
	ctx["title"] = "博客"

	param := mux.Vars(req)

	pageStr := param["page"]
	if pageStr == "" {
		pageStr = "1"
	}

	page, _ := strconv.Atoi(pageStr)

	recordCount, _ := dbMap.SelectInt("select count(*) count from article")

	pageSize := int64(20)
	paging := lib.NewPager(int64(page), pageSize, recordCount)
	offset := (paging.CurrentPage() - 1) * pageSize

	articleList := []models.ArticleExt{}
	dbMap.Select(&articleList, "select a.*,(select count(1) from comment where articleid=a.id) CommentCount from article a where a.status=true order by a.id desc LIMIT $1 OFFSET $2", pageSize, offset)

	ctx["articleList"] = articleList
	ctx["pager"] = paging

	// 把公共数据放入context ctx
	commonContextProcess(req, ctx)

	r.HTML(w, http.StatusOK, "article_list", ctx, render.HTMLOptions{Layout: "layout"})
}

func ArticleCategory(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})
	param := mux.Vars(req)

	pageStr := param["page"]
	if pageStr == "" {
		pageStr = "1"
	}

	categoryId := "0"
	if param["categoryId"] != "" {
		categoryId = param["categoryId"]
		ctx["categoryId"] = categoryId
	}

	page, _ := strconv.Atoi(pageStr)

	recordCount, _ := dbMap.SelectInt("select count(*) count from article where categoryid=$1", categoryId)

	pageSize := int64(20)
	paging := lib.NewPager(int64(page), pageSize, recordCount)
	offset := (paging.CurrentPage() - 1) * pageSize

	articleList := []models.ArticleExt{}
	dbMap.Select(&articleList, "select *,(select count(1) from comment where articleid=article.id) CommentCount from article where categoryid=$3 order by created desc LIMIT $1 OFFSET $2", pageSize, offset, categoryId)

	ctx["articleList"] = articleList
	ctx["pager"] = paging

	categoryName, _ := dbMap.SelectStr("select name from category where id=$1", categoryId)

	ctx["title"] = categoryName

	// 把公共数据放入context ctx
	commonContextProcess(req, ctx)

	r.HTML(w, http.StatusOK, "category", ctx, render.HTMLOptions{Layout: "layout"})
}

func ArticleDetail(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})
	param := mux.Vars(req)

	id := 0
	if param["id"] != "" {
		id, _ = strconv.Atoi(param["id"])
	}
	article := models.ArticleExt{}

	err := dbMap.SelectOne(&article, "select *,(select count(1) from comment where articleid=article.id) CommentCount from article where id=$1", id)
	if err == nil {
		// 增加阅读数
		c := make(chan int)
		go func() {
			dbMap.Exec("update article set read=read+1,updated=$2 where id=$1", id, time.Now())
			c <- 1
		}()
		<-c

		ctx["article"] = article
		ctx["title"] = article.Title
	} else {
		r.Text(w, http.StatusNotFound, "嗷.嗷.嗷...找不到页面了！")
		return
	}

	commentCount := 0
	commentList := []models.Comment{}
	_, err = dbMap.Select(&commentList, "select * from comment where status=true and articleid=$1 order by id asc", id)
	if err == nil {
		ctx["commentList"] = commentList
		commentCount = len(commentList)
	}

	ctx["commentCount"] = commentCount

	ctx["captchaId"] = captcha.New()

	//热门博客
	ctx["hotArticles"] = getHotArticleList()

	//最近浏览的博客
	ctx["recentView"] = getRecentView()

	// 把必需公共数据放入context ctx
	ctx = requiredCtxProcess(req, ctx)

	r.HTML(w, http.StatusOK, "detail", ctx, render.HTMLOptions{Layout: "layout"})
}

func ArticleDisable(w http.ResponseWriter, req *http.Request) {
	logined := _loginProcess(req)
	if logined {
		params := mux.Vars(req)
		id := params["id"]
		dbMap.Exec("update article set status=false where id=$1", id)
		ClearSession(req, w, "getIndexArticle")
		SetFlashMessages(req, w, "移除成功！")
	} else {
		SetFlashMessages(req, w, "无权限移除！")
	}
	http.Redirect(w, req, "/", http.StatusFound)
}

// ajax 方式提价评论
func AddCommentAjax(w http.ResponseWriter, req *http.Request) {
	retMap := make(map[string]interface{})
	if strings.ToUpper(req.Method) == "POST" {
		userId, _ := strconv.ParseInt(req.FormValue("userId"), 10, 64)
		articleId, _ := strconv.ParseInt(req.FormValue("articleId"), 10, 64)
		userName := req.FormValue("userName")
		email := req.FormValue("email")
		replyer := req.FormValue("replyer")
		content := req.FormValue("content")
		blogTile := req.FormValue("blogTile")
		commentId := req.FormValue("commentId")

		vcode := req.FormValue("vcode")
		vcodeSrc := req.FormValue("vcodeSrc")

		if !captcha.VerifyString(vcodeSrc, vcode) {
			retMap["code"] = "2"
			retMap["captchaId"] = captcha.New()
			r.JSON(w, http.StatusOK, retMap)
			return
		}

		refContent := ""
		tmpComment := models.Comment{}
		if commentId != "" {
			err := dbMap.SelectOne(&tmpComment, "select * from comment where id=$1", commentId)
			if err == nil && tmpComment.Id > 0 {
				refContent = fmt.Sprintf(`<div class="ref">%s<h4>引用"%s"的评论</h4>%s</div>`, tmpComment.RefContent, tmpComment.UserName, tmpComment.Content)
			}

		}

		comment := models.Comment{
			UserId:     userId,
			UserName:   userName,
			Email:      email,
			Replyer:    replyer,
			ArticleId:  articleId,
			Content:    content,
			RefContent: refContent,
			Status:     true,
		}

		err := dbMap.Insert(&comment)

		if err == nil {
			retMap["code"] = "1"
			retMap["captchaId"] = captcha.New()
			retMap["commentId"] = comment.Id
			retMap["content"] = comment.RefContent + comment.Content

			domain, _ := settings.String("default", "domain")
			html := ""
			if email != "" && tmpComment.Email != "" {
				html = fmt.Sprintf(`<p>%s:</p><p>%s 在<a href='%s/article/detail/%d#items%d'>《%s》</a>博客中，回复您说：<br/>"%s"</p>`,
					replyer,
					userName,
					domain,
					articleId,
					comment.Id,
					blogTile,
					content)
				sendEmail(email, tmpComment.Email, userName+"在博客中回复了您", html)
			}
			html = ""
			account := getAccount()
			//不是博主评论的时候，才发 email 给自己
			if userId <= 0 && email != "" && account != nil && account.Email != "" {
				html = fmt.Sprintf(`<p>博主:</p><p>%s 在<a href='%s/article/detail/%d#items%d'>《%s》</a>博客中评论说：<br/>"%s"</p>`,
					userName,
					domain,
					articleId,
					comment.Id,
					blogTile,
					content)
				sendEmail(email, account.Email, userName+"在您博客中发表了评论", html)
			}

			r.JSON(w, http.StatusOK, retMap)
			return
		}
	}

	retMap["code"] = "0"
	retMap["captchaId"] = captcha.New()
	r.JSON(w, http.StatusOK, retMap)
}

func SearchArticle(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	ctx := make(map[string]interface{})

	pageStr := req.FormValue("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, _ := strconv.Atoi(pageStr)

	q := req.FormValue("q")
	ctx["title"] = "搜索：" + q
	key := q
	ctx["q"] = q

	recordCount := int64(0)
	if "" != q {
		q = "%" + q + "%"
		recordCount, _ = dbMap.SelectInt("select count(*) count from article where title like $1 or content like $2", q, q)
	} else {
		recordCount, _ = dbMap.SelectInt("select count(*) count from article")
	}

	pageSize := int64(12)
	paging := lib.NewPager(int64(page), pageSize, int64(recordCount))
	offset := (paging.CurrentPage() - 1) * pageSize

	articleList := []models.ArticleExt{}

	if "" != q {
		dbMap.Select(&articleList, "select *,(select count(1) from comment where articleid=article.id) CommentCount from article where title like $3 or content like $4 order by id desc LIMIT $1 OFFSET $2", pageSize, offset, q, q)
	} else {
		dbMap.Select(&articleList, "select *,(select count(1) from comment where articleid=article.id) CommentCount from article order by id desc LIMIT $1 OFFSET $2", pageSize, offset)
	}

	ctx["articleList"] = articleList
	ctx["pager"] = paging

	// 把公共数据放入context ctx
	commonContextProcess(req, ctx)

	ctx["tips"] = fmt.Sprintf("站内搜索%s，找到相关结果约%d个，耗时约%v", key, recordCount, time.Since(start))

	r.HTML(w, http.StatusOK, "search", ctx, render.HTMLOptions{Layout: "layout"})
}

func ArticlePraise(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("articleId")
	if id != "" {
		dbMap.Exec("update article set upcount=upcount+1 where id=$1", id)
		SetSession(req, w, "praiseTime", time.Now())
		r.Text(w, http.StatusOK, "1")
		return
	}
	r.Text(w, http.StatusOK, "0")
}

func ArticleArchives(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})
	ctx["title"] = "博客归档"

	param := mux.Vars(req)

	pageStr := param["page"]
	if pageStr == "" {
		pageStr = "1"
	}

	page, _ := strconv.Atoi(pageStr)

	year := param["year"]
	month := param["month"]
	if year == "" || month == "" {
		r.Text(w, http.StatusNotFound, "嗷.嗷.嗷...找不到页面了！")
		return
	}

	recordCount, _ := dbMap.SelectInt("select count(*) count from article where status=true and extract(year from created)=$1 and extract(month from created)=$2", year, month)

	pageSize := int64(20)
	paging := lib.NewPager(int64(page), pageSize, recordCount)
	offset := (paging.CurrentPage() - 1) * pageSize

	articleList := []models.ArticleExt{}
	dbMap.Select(&articleList, "select a.*,(select count(1) from comment where articleid=a.id) CommentCount from article a where a.status=true and extract(year from created)=$3 and extract(month from created)=$4 order by a.id desc LIMIT $1 OFFSET $2", pageSize, offset, year, month)

	ctx["articleList"] = articleList
	ctx["pager"] = paging

	// 把公共数据放入context ctx
	commonContextProcess(req, ctx)

	r.HTML(w, http.StatusOK, "article_list", ctx, render.HTMLOptions{Layout: "layout"})
}
