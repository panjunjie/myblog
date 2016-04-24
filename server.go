package main

import (
	"github.com/codegangsta/negroni"
	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
	"myblog/src/conf"
	"myblog/src/controller"
	"myblog/src/middleware"

	"github.com/robfig/cron"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	//前台
	router.HandleFunc("/", controller.Index)
	router.HandleFunc("/index", controller.Index)
	router.HandleFunc("/article/{page:[0-9]+}", controller.ArticleList)

	router.HandleFunc("/category/{categoryId:[0-9]+}", controller.ArticleCategory)
	router.HandleFunc("/category/{categoryId:[0-9]+}/{page:[0-9]+}", controller.ArticleCategory)

	router.HandleFunc("/article/detail/{id:[0-9]+}", controller.ArticleDetail)
	router.HandleFunc("/article/{id:[0-9]+}/disable", controller.ArticleDisable)
	router.HandleFunc("/article/article_praise_ajax", controller.ArticlePraise)

	router.Handle("/captcha/{[a-z]+}", captcha.Server(240, 80))

	router.HandleFunc("/article/search", controller.SearchArticle)

	router.HandleFunc("/add_comment_ajax", controller.AddCommentAjax)
	router.HandleFunc("/article/rss", controller.ArticleRss)

	router.HandleFunc("/archives/{year}/{month}/", controller.ArticleArchives)

	router.HandleFunc("/friendlink", controller.FriendLinkList)
	router.HandleFunc("/friendlink/add", controller.AddFriendLinkAjax)
	router.HandleFunc("/friendlink/pass", controller.PassFriendLinkAjax)

	router.HandleFunc("/about", controller.AboutMe)

	router.HandleFunc("/account/login", controller.AccountLogin)
	router.HandleFunc("/account/logout", controller.AccountLogout)

	//后台
	adminRoutes := mux.NewRouter()

	adminRoutes.HandleFunc("/admin/category/list", controller.CategoryListAdmin)
	adminRoutes.HandleFunc("/admin/category/add", controller.CategoryAddAdmin)
	adminRoutes.HandleFunc("/admin/category/{id:[0-9]+}/edit", controller.CategoryUpdateAdmin)
	adminRoutes.HandleFunc("/admin/category/{id:[0-9]+}/del", controller.CategoryDelAdmin)

	adminRoutes.HandleFunc("/admin/article/list", controller.ArticleListAdmin)
	adminRoutes.HandleFunc("/admin/article/add", controller.ArticleAddAdmin)
	adminRoutes.HandleFunc("/admin/article/{id:[0-9]+}/edit", controller.ArticleUpdateAdmin)
	adminRoutes.HandleFunc("/admin/article/{id:[0-9]+}/del", controller.ArticleDelAdmin)

	adminRoutes.HandleFunc("/admin/friendlink/list", controller.FriendLinkListAdmin)
	adminRoutes.HandleFunc("/admin/friendlink/add", controller.FriendLinkAddAdmin)
	adminRoutes.HandleFunc("/admin/friendlink/{id:[0-9]+}/edit", controller.FriendLinkEditAdmin)
	adminRoutes.HandleFunc("/admin/friendlink/pass", controller.FriendLinkPassAdmin)
	adminRoutes.HandleFunc("/admin/friendlink/del", controller.FriendLinkDelAdmin)

	adminRoutes.HandleFunc("/admin/user/edit", controller.AboutMeEditAdmin)

	adminRoutes.HandleFunc("/admin/backup/myblog.xml", controller.BackUpBlog)
	adminRoutes.HandleFunc("/admin/settings/edit", controller.SettingsEditAdmin)
	adminRoutes.HandleFunc("/admin/advert/edit", controller.AdvertEditAdmin)

	adminRoutes.HandleFunc("/admin/flush/data", controller.FlushDataAdmin)

	adminRoutes.HandleFunc("/admin/kindeditor/upload", controller.KindEditorUploadJson)
	adminRoutes.HandleFunc("/admin/kindeditor/file/manager", controller.KindEditorFileManageJson)

	router.PathPrefix("/admin").Handler(negroni.New(
		middleware.NewCheckLogin(),
		negroni.Wrap(adminRoutes),
	))

	settings := conf.Settings
	debug, err := settings.Bool("default", "debug")
	if err != nil {
		debug = false
	}

	n := negroni.New()

	if debug {
		n.Use(negroni.NewLogger())
		n.Use(negroni.NewStatic(http.Dir(".")))
	}

	recovery := negroni.NewRecovery()
	recovery.PrintStack = debug
	n.Use(recovery)

	n.UseHandler(router)

	//作业计划
	c := cron.New()
	c.AddFunc("@hourly", controller.RunAllStaticstics)
	c.Start()

	port, err := settings.String("default", "port")
	if err != nil {
		port = "3003"
	}

	n.Run(":" + port)
}
