package controller

import (
	"github.com/dchest/captcha"
	"github.com/unrolled/render"
	"myblog/src/models"
	"net/http"
	"strings"
	"time"
)

func FriendLinkList(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})

	ctx["title"] = "友情交换链接"

	ctx = loginProcess(req, ctx)

	friendLinkList := []models.FriendLink{}
	if ctx["logined"].(bool) {
		dbMap.Select(&friendLinkList, "select * from friendlink order by status asc,id asc")
	} else {
		dbMap.Select(&friendLinkList, "select * from friendlink where status=true order by id asc")
	}

	ctx["friendLinkList"] = friendLinkList

	ctx["friendLinkCount"] = len(friendLinkList)

	ctx["captchaId"] = captcha.New()

	// 把必需公共数据放入context ctx
	ctx = requiredCtxProcess(req, ctx)

	r.HTML(w, http.StatusOK, "friendlink", ctx, render.HTMLOptions{Layout: "layout"})

}

// 添加友情链接
func AddFriendLinkAjax(w http.ResponseWriter, req *http.Request) {
	retMap := make(map[string]interface{})
	if strings.ToUpper(req.Method) == "POST" {
		name := req.FormValue("name")

		vcode := req.FormValue("vcode")
		vcodeSrc := req.FormValue("vcodeSrc")

		if !captcha.VerifyString(vcodeSrc, vcode) {
			retMap["code"] = "2"
			retMap["captchaId"] = captcha.New()
			r.JSON(w, http.StatusOK, retMap)
			return
		}

		count, _ := dbMap.SelectInt("select count(1) from friendlink where name=$1", name)

		if count > 0 {
			// 记录已经存在
			retMap["code"] = "3"
			retMap["captchaId"] = captcha.New()
			r.JSON(w, http.StatusOK, retMap)
			return
		}

		friendlink := models.FriendLink{
			Name:    name,
			Link:    req.FormValue("link"),
			Intro:   req.FormValue("intro"),
			Read:    1,
			Status:  false,
			Created: time.Now(),
			Updated: time.Now(),
		}

		err := dbMap.Insert(&friendlink)

		if err == nil {
			retMap["code"] = "1"
			retMap["captchaId"] = captcha.New()
			r.JSON(w, http.StatusOK, retMap)
			return
		}
	}

	retMap["code"] = "0"
	retMap["captchaId"] = captcha.New()
	r.JSON(w, http.StatusOK, retMap)
	return
}

func PassFriendLinkAjax(w http.ResponseWriter, req *http.Request) {
	retMap := make(map[string]interface{})
	if strings.ToUpper(req.Method) == "POST" {
		if GetSession(req, SESSION_WEB) != nil {
			id := req.FormValue("id")
			action := req.FormValue("action")

			status := false
			if action == "1" {
				status = true
			}

			dbMap.Exec("update friendlink set status=$1 where id=$2", status, id)
			retMap["code"] = "1"
			r.JSON(w, http.StatusOK, retMap)
			return
		}
	}
	retMap["code"] = "0"
	r.JSON(w, http.StatusOK, retMap)
}
