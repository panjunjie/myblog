package controller

import (
	"github.com/gorilla/csrf"
	"github.com/unrolled/render"
	"myblog/src/lib"
	"myblog/src/models"
	"net/http"
	"time"
)

var (
	LOGIN_ADMIN = "LOGIN_ADMIN"
)

func AccountLogin(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})

	var next = req.FormValue("next")
	ctx["next"] = next

	ctx[csrf.TemplateTag] = csrf.TemplateField(req)
	var account = models.Account{}
	if req.Method == "POST" {
		var username = req.FormValue("username")
		var password = req.FormValue("password")

		if username != "" && password != "" {
			err := dbMap.SelectOne(&account, "select * from account where username=$1 or email=$2 limit 1", username, username)

			if err != nil && account.Id > 0 {
				SetFlashMessages(req, w, "用户不存在！")
				http.Redirect(w, req, "/account/login", http.StatusFound)
				return
			}

			if account.Password == lib.MD5(password) {
				SetSession(req, w, LOGIN_ADMIN, account)
				dbMap.Exec("update account set lastlogin=$1 where $2", time.Now(), account.Id)

				if next != "" {
					http.Redirect(w, req, next, http.StatusFound)
					return
				} else {
					http.Redirect(w, req, "/", http.StatusFound)
					return
				}
			} else {
				SetFlashMessages(req, w, "账号和密码不匹配")
				http.Redirect(w, req, "/account/login", http.StatusFound)
				return
			}
		}
	} else {
		dbMap.SelectOne(&account, "select * from account limit 1")
		if account.Id <= 0 {
			account.UserName = "admin"
			account.Password = lib.MD5("admin")
			account.Status = true
			dbMap.Insert(&account)
			SetFlashMessages(req, w, "初始化账号和密码：admin / admin，请及时更改密码！")
		}

		ctx["flashes"] = GetFlashMessages(req, w)

		if currentUser := GetSession(req, LOGIN_ADMIN); currentUser != nil {
			if next == "" {
				next = "/"
			}
			http.Redirect(w, req, next, http.StatusFound)
			return
		}

		r.HTML(w, http.StatusOK, "account/login", ctx, render.HTMLOptions{Layout: ""})
		return
	}
}

func AccountLogout(w http.ResponseWriter, req *http.Request) {
	ClearSession(req, w, LOGIN_ADMIN)
	ClearAllSession(req, w)
	http.Redirect(w, req, "/", http.StatusFound)
}
