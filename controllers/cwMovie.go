package controllers

import(
	"beego_movie/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"time"
)

type CwMovieController struct {
	beego.Controller
}
func (c *CwMovieController) CwMovie(){
	var movieIn models.MovieIn  //models里面movie_in数据结构
	models.ConnectRedis("127.0.0.1:6379")
	//爬虫url
	Yurl :="https://movie.douban.com/subject/26794435/"
	models.PutinQueue(Yurl)
	for{
		length :=models.GetQueueLength()
		if length ==0 {
			break //如果url队列为空，就终止继续循环
		}
		Yurl =models.PopformQueue()
		//判断Yurl是否访问过
		if models.IsVisit(Yurl){
			continue
		}
		//然后http请求GET方法
		ysp :=httplib.Get(Yurl)
		//设置User-agent以及cookie是为了防止403
		ysp.Header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:50.0) Gecko/20100101 Firefox/50.0")
		ysp.Header("Cookie",`bid=tv0lb0GApLA; __utmz=30149280.1565105160.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utmz=223695111.1565105160.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); ll="118172"; __yadk_uid=rPXU2pV08tosySZITVKc6WqiDxg1yM34; trc_cookie_storage=taboola%2520global%253Auser-id%3D26be11bb-c5c3-4089-ac13-5c97c018848e-tuct443da79; _vwo_uuid_v2=D0BE3A7C8D6FDBAAFB2F24B1968F1641E|7a9598107c12f0e75124be05627ec0e4; _pk_id.100001.4cf6=bf61d8d98743d47a.1565105160.4.1565721031.1565690492.; _pk_ses.100001.4cf6=*; __utma=30149280.1917044307.1565105160.1565690451.1565721031.4; __utmb=30149280.0.10.1565721031; __utmc=30149280; __utma=223695111.520476801.1565105160.1565690451.1565721031.4; __utmb=223695111.0.10.1565721031; __utmc=223695111; __gads=ID=d0b1ccfd8077bfba:T=1565721033:S=ALNI_MbtqdWrTe5WLVeoldAi-uBeQJLO8Q`)
		//获取url数据
		yMovieHtml, err :=ysp.String()
		if err!=nil{
			panic(err)
		}
		//获取电影名称
		movieIn.Movie_name   			=models.GetModeName(yMovieHtml)
		//如果为空，不是电影，如果不为空，是电影
		if movieIn.Movie_name !=""{
			movieIn.Movie_pic 					=models.GetMoviePic(yMovieHtml)//图片
			movieIn.Movie_director      		=models.GetMovieDirector(yMovieHtml)//导演1
			movieIn.Movie_writer				=models.GetMovieBianju(yMovieHtml) //编剧
			movieIn.Movie_country               =models.Getage(yMovieHtml)  //产地
			movieIn.Movie_language              =models.GetMovieLanguage(yMovieHtml) //语言
			movieIn.Movie_main_character 		=models.GetMovieMainCharacters(yMovieHtml) //主演1
			movieIn.Movie_type                  =models.GetMovieGenre(yMovieHtml)  //类型1
			movieIn.Movie_on_time               =models.GetMovieOnTime(yMovieHtml) //上映时间1
			movieIn.Movie_span                  =models.GetMovieRunningTime(yMovieHtml) //片长1
			movieIn.Movie_grade					=models.GetMovieGrade(yMovieHtml) //评分1
			//存入数据库
			models.AddMovie(&movieIn)
		}

		//提取该页面的所有连接
		urls :=models.GetMovieUrls(yMovieHtml)
		//redis
		//遍历url
		//为了把url写入队列
		//同样需要开启一个协程，这个协程专门负责从队列中取，负责get，set，
		//第一判断这个url是不是一个电影，是的话加入到数据库，
		//第二是提取这个电影有关的url
		//第三把url放入set(集合)里，表明这个url已经访问过
		for _,url :=range urls{
			models.PutinQueue(url)
			c.Ctx.WriteString("<br>"+url+"<br/>")
		}
		//Yurl要记录到set集合里，表明这个url访问过
		   models.AddToSet(Yurl)
		   time.Sleep(time.Second) //适当休息
	}
	       c.Ctx.WriteString("爬虫结束了")
}
