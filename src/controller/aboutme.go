package controller

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/unrolled/render"
	"myblog/src/models"
	"net/http"
	"strings"
	"time"
)

func AboutMe(w http.ResponseWriter, req *http.Request) {
	if strings.ToUpper(req.Method) == "POST" {

		replyId := req.FormValue("replyid")
		refContent := ""
		tmpMyMsg := models.MyMsg{}
		if replyId != "" {
			err := dbMap.SelectOne(&tmpMyMsg, "select * from mymsg where id=$1", replyId)
			if err == nil && tmpMyMsg.Id > 0 {
				refContent = fmt.Sprintf(`<div class="ref">%s<h4>引用"%s"的留言</h4>%s</div>`, tmpMyMsg.RefContent, tmpMyMsg.UserName, tmpMyMsg.Content)
			}
		}

		mymsg := models.MyMsg{
			UserName:   req.FormValue("username"),
			Email:      req.FormValue("email"),
			Content:    req.FormValue("content"),
			RefContent: refContent,
			Status:     true,
			Created:    time.Now(),
			Updated:    time.Now(),
		}

		if !captcha.VerifyString(req.FormValue("vcodeSrc"), req.FormValue("vildCode")) {
			SetFlashMessages(req, w, "验证码有误！")
		} else {
			err := dbMap.Insert(&mymsg)
			if err == nil {

				domain, _ := settings.String("default", "domain")
				account := getAccount()

				//给网友发 email
				if replyId != "" && tmpMyMsg.Id > 0 {
					if tmpMyMsg.Email != "" && mymsg.Email != "" && tmpMyMsg.Email != mymsg.Email {
						html := fmt.Sprintf(`<p>%s:</p><p>%s 在 <a href="%s/about">博客</a> 中，回复您说：<br/>"%s"</p>`,
							tmpMyMsg.UserName,
							mymsg.UserName,
							domain,
							mymsg.Content)
						sendEmail(mymsg.Email, tmpMyMsg.Email, mymsg.UserName+"在博客中回复您了", html)
					}
				}

				//给博主发 Email
				if mymsg.Email != "" && account.Email != mymsg.Email {
					html := fmt.Sprintf(`<p>%s:</p><p>%s 在您 <a href="%s/about">博客</a> 中，给您留言说：<br/>"%s"</p>`,
						"博主",
						mymsg.UserName,
						domain,
						mymsg.Content)
					sendEmail(mymsg.Email, account.Email, mymsg.UserName+"在您博客中留言了", html)
				}

				SetFlashMessages(req, w, "留言成功！")
			} else {
				SetFlashMessages(req, w, "留言失败啦！")
			}
		}
		http.Redirect(w, req, "/about", http.StatusFound)
	} else {
		ctx := make(map[string]interface{})
		ctx["title"] = "关于我"

		//留言
		msgList := []models.MyMsg{}
		dbMap.Select(&msgList, "select * from mymsg where status=true order by id asc")

		ctx["msgList"] = msgList
		ctx["msgCount"] = len(msgList)

		ctx["captchaId"] = captcha.New()

		// 把必需公共数据放入context ctx
		ctx = requiredCtxProcess(req, ctx)

		ctx["flashes"] = GetFlashMessages(req, w)

		r.HTML(w, http.StatusOK, "about_me", ctx, render.HTMLOptions{Layout: "layout"})
	}
}

// func TestReq(w http.ResponseWriter, req *http.Request) {
// 	c := new(http.Client)
// 	reqs := request.NewRequest(c)
// 	resp, _ := reqs.Get("http://service.lalawitkey.com/account/login")
// 	//j, _ := resp.Json()
// 	txt, _ := resp.Text()
// 	r.JSON(w, http.StatusOK, txt)
// }
