# myblog
学习 Go 语言，练手的个人博客。

自己学习 Go语言，只看一些语法，不实践，感觉有点缺点什么，所以干脆开发一个自己的博客，很普通很通用的那种，开源出来做个纪念，但有点苦恼，一时想不到好的名字，暂时叫它 myblog 吧，放到 github 里。

myblog 这么一个简单的 web 项目，竟然花了好几天的时间，想想做程序的就是有点苦逼，东西改来改去的，看似小不点，花费的时间就这么没了。

我个人比较喜欢 大猩猩的东西，这个组织的东西很低调，代码质量很高，而且扩展性好，最关键是他们几个人口碑好，虽产品一直不是很热门，但很多名人光顾他们的代码，时不时推崇他们的做法。所以：
1. 路由器采用 mux  https://github.com/gorilla/mux  
2. 会话 session 采用 https://github.com/gorilla/sessions 
3. 订阅 rss 采用 https://github.com/gorilla/feeds

另外 codegangsta 也不错，Martini 名声在外，干脆使用他的 web 中间件 negroni
4. web 中间件 github.com/codegangsta/negroni

还有一些很有用的 中间件，都是看 codegangsta  的东西 偶尔发现的
5. html json xml 组件采用 https://github.com/unrolled/render
6. 表单绑定 中间件采用 https://github.com/mholt/binding

数据库是重要的成员
7. 数据库采用 PostgreSQL 9.4.5 http://get.enterprisedb.com/postgresql/postgresql-9.4.5-3-linux-x64.run
8. 数据库操作中间件 采用 gorp  https://github.com/go-gorp/gorp

不得不说 gorp 很好，我曾经尝试过 sqlx  https://github.com/jmoiron/sqlx 后来由于误解，又重回到 gorp 了，sqlx 其实也不错，特别发现原生 database/sql 不够用了，sqlx 是不二之选，作者是 大猩猩的成员之一。

最后发现用的东西都是 非常冷门的组件，没有一个是知名的，哎，不得不说：受 Go 语言影响可不浅，热衷的都是小众的东西。

myblog 功能点：
1：发布博客
2：发表、回复评论，支持盖楼式评论。
3：支持富文本上传单张或多张图片。
