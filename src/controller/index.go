/*
	博客的控制器
*/
package controller

import (
	"github.com/pmylund/go-cache"
	"github.com/unrolled/render"
	"myblog/src/models"
	"net/http"
)

//博客首页 获取最新20条文章
func Index(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})
	ctx["title"] = "首页"

	articleList := &[]models.ArticleExt{}
	obj, found := cacheObj.Get("getIndexArticle")
	if !found || obj == nil {
		dbMap.Select(articleList, "select a.*,(select count(1) from comment where articleid=a.id) commentcount from article a where a.status=true order by a.id desc LIMIT 20")
		cacheObj.Set("getIndexArticle", articleList, cache.DefaultExpiration)

	} else {
		articleList = obj.(*[]models.ArticleExt)
	}

	ctx["articleList"] = articleList

	// 把公共数据放入context ctx
	commonContextProcess(req, ctx)

	r.HTML(w, http.StatusOK, "index", ctx, render.HTMLOptions{Layout: "layout"})
}
