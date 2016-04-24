package controller

import (
	"fmt"
	"github.com/gorilla/feeds"
	"html/template"
	"myblog/src/models"
	"net/http"
	"time"
)

var site_index_url = ""

func ArticleRss(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, getArtileFeedString("rss"))
}

func ArticleAtom(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, getArtileFeedString("atom"))
}

func getArtileFeedString(myType string) string {
	site_index_url, _ = settings.String("default", "base-url")

	account := getAccount()
	userName := "MyBlog"
	userIntro := "MyBlog"
	userEmail := ""
	if account.Id != 0 {
		userName = account.UserName
		userIntro = account.Intro
		userEmail = account.Email
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       userName + "的博客",
		Link:        &feeds.Link{Href: site_index_url},
		Description: userIntro,
		Author:      &feeds.Author{userName, userEmail},
		Created:     now,
	}

	var articles []models.Article
	dbMap.Select(&articles, "select * from article where status=true order by created desc")
	if articles != nil {
		size := len(articles)
		var items = make([]*feeds.Item, size)

		for i := 0; i < size; i++ {
			the_time, _ := time.Parse(TIMEFORMAT, articles[i].Created.Format(TIMEFORMAT))

			item := &feeds.Item{
				Id:          fmt.Sprintf("%d", articles[i].Id),
				Title:       articles[i].Title,
				Link:        &feeds.Link{Href: fmt.Sprintf("%s/article/detail/%d", site_index_url, articles[i].Id)},
				Description: fmt.Sprintf("%s", template.HTML(articles[i].Content)),
				Author:      &feeds.Author{userName, userEmail},
				Created:     the_time,
			}
			items[i] = item

		}
		feed.Items = items
	}

	if myType == "rss" {
		rss, _ := feed.ToRss()
		return rss
	} else {
		atom, _ := feed.ToAtom()
		return atom
	}
}
