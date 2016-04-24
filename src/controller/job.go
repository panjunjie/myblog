package controller

//定时作业
import (
	"fmt"
	"myblog/src/models"
)

/*计算分类的博客数*/
func RunCategory() {
	index := 1
	if index <= 1 {
		//第一次启动项目的时候不执行,第二次的时候执行
		showInfo("统计分类博客数作业开始...")
		dbMap.Exec("update category a set itemcount=(select count(1) from article where status=true and categoryid=a.id) where a.status=true")
		showInfo("统计分类博客数作业结束^_^")
	}
	index = index + 1
}

//statistics category 1:统计博客数评论数访问数点赞数；2：归档博客；3：热门博客；4：最新点评博客；5：热评博客;6：点赞博客排名

func RunAllStaticstics() {
	RunArchives()
	RunHotBlog()
	RunHotPraise()
	RunNewComment()
	RunHotComment()
	RunTotalCount()
	RunCategory()
}

//统计总博客数，总评论数，总浏览数，总点赞数
func RunTotalCount() {
	index := 1
	if index <= 1 {
		//第一次启动项目的时候不执行,第二次的时候执行
		showInfo("统计总博客数，总评论数，总浏览数，总点赞数 作业开始...")
		dbMap.Exec("delete from statistics where category=$1", 1)

		blogTotal, _ := dbMap.SelectInt("select count(1) from article where status=true")
		commentTotal, _ := dbMap.SelectInt("select count(1) from comment where status=true")
		viewTotal, _ := dbMap.SelectInt("select sum(read) from article where status=true")
		PraiseTotal, _ := dbMap.SelectInt("select sum(upcount) from article where status=true")

		total := models.Statistics{}
		total.Category = 1
		total.Param1 = fmt.Sprintf("%d", blogTotal)
		total.Param2 = fmt.Sprintf("%d", commentTotal)
		total.Param3 = fmt.Sprintf("%d", viewTotal)
		total.Param4 = fmt.Sprintf("%d", PraiseTotal)

		err := dbMap.Insert(&total)
		showInfo(err)

	}
	index = index + 1
}

/**
**博客归档统计
**/
func RunArchives() {
	index := 1
	if index <= 1 {
		//第一次启动项目的时候不执行,第二次的时候执行
		showInfo("统计归档作业开始...")
		dbMap.Exec("delete from statistics where category=$1", 2)
		statistics_list := []models.Statistics{}
		_, err := dbMap.Select(&statistics_list, "select 2 category,to_char(created,'yyyy年mm月') param1,to_char(created,'yyyy') param2,to_char(created,'mm') param3,count(id) param4 from article group by to_char(created,'yyyy年mm月'),to_char(created,'yyyy'),to_char(created,'mm') order by to_char(created,'yyyy年mm月') desc")

		if err == nil {
			for _, model := range statistics_list {
				dbMap.Insert(&model)
			}
		}

	}
	index = index + 1
}

//季度浏览最多的N篇博客
func RunHotBlog() {
	index := 1
	if index <= 1 {
		//第一次启动项目的时候不执行,第二次的时候执行
		showInfo("季度浏览最多的N篇博客...")
		dbMap.Exec("delete from statistics where category=$1", 3)
		statistics_list := []models.Statistics{}
		_, err := dbMap.Select(&statistics_list, "select id param1,3 category,title title1,read param2 from article where created >= (select created from article order by created desc limit 1) - interval '3 month' order by read desc limit 10")
		if len(statistics_list) < 10 {
			statistics_list = []models.Statistics{}
			dbMap.Select(&statistics_list, "select id param1,3 category,title title1,read param2 from article order by read desc limit 10")
		}
		if err == nil {
			for _, model := range statistics_list {
				dbMap.Insert(&model)
			}
		}

	}
	index = index + 1
}

//季度最新点评的博客
func RunNewComment() {
	index := 1
	if index <= 1 {
		//第一次启动项目的时候不执行,第二次的时候执行
		showInfo("季度最新点评的N篇博客...")
		dbMap.Exec("delete from statistics where category=$1", 4)
		statistics_list := []models.Statistics{}
		_, err := dbMap.Select(&statistics_list, "select a.id param1,4 category,a.title title1,b.content title2,b.username param2,to_char(b.created,'yyyy-MM-dd hh24:MI:ss') param3,b.id param4 from article a join comment b on a.id=b.articleid where a.status=true and b.status=true order by b.id desc limit 10")

		if err == nil {
			for _, model := range statistics_list {
				dbMap.Insert(&model)
			}
		}

	}
	index = index + 1
}

//季度点评最多的博客
func RunHotComment() {
	index := 1
	if index <= 1 {
		//第一次启动项目的时候不执行,第二次的时候执行
		showInfo("季度点评最多的N篇博客...")
		dbMap.Exec("update article a set commentcount=(select count(1) from comment where articleid=a.id)")
		dbMap.Exec("delete from statistics where category=$1", 5)
		statistics_list := []models.Statistics{}
		_, err := dbMap.Select(&statistics_list, "select a.id param1,5 category,a.title title1,b.content title2,b.username param2,to_char(b.created,'yyyy-MM-dd hh24:MI:ss') param3,b.id param4 from article a join comment b on a.id=b.articleid where a.created >= (select created from article order by created desc limit 1) - interval '3 month' order by a.commentcount desc limit 10")

		if err == nil {
			for _, model := range statistics_list {
				dbMap.Insert(&model)
			}
		}

	}
	index = index + 1
}

//季度点赞最多的博客
func RunHotPraise() {
	index := 1
	if index <= 1 {
		//第一次启动项目的时候不执行,第二次的时候执行
		showInfo("季度点赞最多的N篇博客...")
		dbMap.Exec("delete from statistics where category=$1", 6)
		statistics_list := []models.Statistics{}
		_, err := dbMap.Select(&statistics_list, "select id param1,6 category,title title1,upcount param2 from article where created >= (select created from article order by created desc limit 1) - interval '3 month' order by upcount desc limit 10")

		if err == nil {
			for _, model := range statistics_list {
				dbMap.Insert(&model)
			}
		}

	}
	index = index + 1
}
