package controller

import (
	"github.com/go-gorp/gorp"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/pmylund/go-cache"
	"github.com/robfig/config"
	"github.com/unrolled/render"
	"html/template"
	"log"
	"myblog/src/conf"
	"myblog/src/dbutils"
	"myblog/src/lib"
	"net/http"
	"time"
)

const (
	SESSION_WEB           = "SESSION_WEB"
	SESSION_FLASH_MESSAGE = "SESSION_FLASH_MESSAGE"
	SESSION_FLASH_OBJECT  = "SESSION_FLASH_OBJECT"
	STORE_KEY             = "BLOGS"
	TIMEFORMAT            = "2006-01-02 15:04:05"
)

var (
	dbSqlx   *sqlx.DB
	dbMap    *gorp.DbMap
	r        *render.Render
	settings *config.Config

	MinTime int64 = 1136185445 //2006-01-02 15:04:05

	Store = sessions.NewCookieStore([]byte(STORE_KEY))

	SessionFlash *sessions.Session
	SessionWeb   *sessions.Session

	cacheObj = cache.New(15*time.Minute, 30*time.Second)
)

//控制器启动的时候 数据库连接初始化
func init() {
	settings = conf.Settings

	debug, err := settings.Bool("default", "debug")
	if err != nil {
		debug = false
	}

	r = render.New(render.Options{
		Directory:       "templates",                    // Specify what path to load the templates from.
		Layout:          "layout",                       // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions:      []string{".html"},              // Specify extensions to load for templates.
		Funcs:           []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
		Delims:          render.Delims{"{{", "}}"},      // Sets delimiters to the specified strings.
		Charset:         "UTF-8",                        // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON:      debug,                          // Output human readable JSON.
		IndentXML:       debug,                          // Output human readable XML.
		PrefixJSON:      []byte(""),                     // Prefixes JSON responses with the given bytes.
		PrefixXML:       []byte(""),                     // Prefixes XML responses with the given bytes.
		HTMLContentType: "text/html",                    // Output XHTML content type instead of default "text/html".
		IsDevelopment:   debug,                          // Render will now recompile the templates on every HTML response.
		UnEscapeHTML:    true,                           // Replace ensure '&<>' are output correctly (JSON only).
		StreamingJSON:   true,                           // Streams the JSON response via json.Encoder.
		RequirePartials: true,
	})

	dbSqlx = dbutils.Db
	dbMap = dbutils.DbMap
}

var AppHelpers = template.FuncMap{
	"html": func(x string) interface{} {
		return template.HTML(x)
	},
	"substring": func(str string, lens int) string {
		return lib.ShowSubstr(str, lens)
	},
	"add": func(a, b int) int {
		return a + b
	},
	"currentDate": func(format string) interface{} {
		return time.Now().Format(format)
	},
	"timeFormat": func(datetime time.Time, format string) interface{} {
		return datetime.Format(format)
	},
	"timeFormat2": func(t int64, format string) interface{} {
		if t > 0 {
			return time.Unix(t, 0).Format(format)
		}
		return TIMEFORMAT
	},
}

func SetFlashMessages(req *http.Request, w http.ResponseWriter, msg string) {
	SessionFlash, _ = Store.Get(req, SESSION_FLASH_MESSAGE)
	SessionFlash.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
	}
	SessionFlash.AddFlash(msg)
	SessionFlash.Save(req, w)
}

func GetFlashMessages(req *http.Request, w http.ResponseWriter) []interface{} {
	SessionFlash, _ := Store.Get(req, SESSION_FLASH_MESSAGE)
	var msg = SessionFlash.Flashes()
	SessionFlash.Save(req, w)
	return msg
}

func SetFlashObject(req *http.Request, w http.ResponseWriter, obj interface{}) {
	if obj != nil {
		SessionFlash, _ = Store.Get(req, SESSION_FLASH_OBJECT)
		SessionFlash.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   300,
			HttpOnly: true,
		}
		SessionFlash.AddFlash(obj)
		SessionFlash.Save(req, w)
	}
}

func GetFlashObject(req *http.Request, w http.ResponseWriter) (bool, interface{}) {
	SessionFlash, _ = Store.Get(req, SESSION_FLASH_OBJECT)
	var obj = SessionFlash.Flashes()
	SessionFlash.Save(req, w)
	if size := len(obj); size > 0 {
		return true, obj[size-1]
	}
	return false, nil
}

func SetSession(req *http.Request, w http.ResponseWriter, key string, value interface{}) {
	SessionWeb, _ = Store.Get(req, SESSION_WEB)
	SessionWeb.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 10 * 24,
		HttpOnly: true,
	}
	SessionWeb.Values[key] = value
	SessionWeb.Save(req, w)
}

func GetSession(req *http.Request, key string) interface{} {
	SessionWeb, _ = Store.Get(req, SESSION_WEB)
	return SessionWeb.Values[key]
}

//获取到key的value后，清除该session
func PopSession(req *http.Request, w http.ResponseWriter, key string) interface{} {
	SessionWeb, _ = Store.Get(req, SESSION_WEB)
	v := SessionWeb.Values[key]
	SessionWeb.Values[key] = nil
	SessionWeb.Save(req, w)
	return v
}

func ClearSession(req *http.Request, w http.ResponseWriter, key string) {
	SessionWeb, _ = Store.Get(req, SESSION_WEB)
	if nil != SessionWeb {
		SessionWeb.Values[key] = nil
		SessionWeb.Save(req, w)
	}
}

func ClearAllSession(req *http.Request, w http.ResponseWriter) {
	SessionWeb, _ = Store.Get(req, SESSION_WEB)
	if nil != SessionWeb {
		flashes := SessionWeb.Flashes()
		size := len(flashes)
		if size > 0 {
			for i := 0; i < size; i++ {
				flashes[i] = nil
			}
			SessionWeb.Save(req, w)
		}
	}
}

//获取当前页的url地址（包含参数）
func getCurrPathQuery(req *http.Request) string {
	var query = ""
	if req.URL.RawQuery != "" {
		query = "?" + req.URL.RawQuery
	}
	return req.URL.Path + query
}

func checkErr(err error, str string) {
	if err != nil {
		log.Fatalf("发生错误：%s\n%v", str, err)
	}
}

func showErr(err error, str string) {
	if err != nil {
		log.Printf("发生错误：%s\n%v", str, err)
	}
}

func showMsg(str string) {
	if str != "" {
		log.Println(str)
	}
}

func showInfo(obj interface{}) {
	debug, err := settings.Bool("default", "debug")
	if err != nil {
		debug = false
	}
	if obj != nil && debug {
		log.Println(obj)
	}
}
