package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"regexp"
	"strings"
)

var (
	db orm.Ormer
)

type MovieIn struct{
	Id int64
	Movie_id int64
	Movie_name string
	Movie_pic string
	Movie_director string
	Movie_writer string
	Movie_country string
	Movie_language string
	Movie_main_character string
	Movie_type string
	Movie_on_time string
	Movie_span string
	Movie_grade string
}
func init(){
	orm.Debug=true
	orm.RegisterDataBase("default", "mysql", "root:root1128@tcp(127.0.0.1:3306)/test?charset=utf8", 30)
	orm.RegisterModel(new(MovieIn)) //生成数据表表单映射
	db = orm.NewOrm()
}
//添加电影
func AddMovie(movie_in *MovieIn)(int64 ,error){
	movie_in.Id=0
	id,err :=db.Insert(movie_in)
	return  id ,err
}

//----------------------------------------------------------------------------------------------------------------------
func GetModeName(movieHtml string) string{
	if movieHtml == ""{
		return ""
	}
	//regexp是正则的包 写正则的规则
	reg :=regexp.MustCompile(`<span\s*property="v:itemreviewed">(.*?)</span>`)
	//然后进行匹配 -1 表示全部返回 如果写一个1 他就返回匹配到的第一个  然后返回是一个[][]string
	result :=reg.FindAllStringSubmatch(movieHtml,-1)
	//没匹配内容返回为空
	if len(result)==0{
		return ""
	}
	return string(result[0][1])
}

//图片地址
func GetMoviePic(movieHtml string) string{
	if movieHtml == ""{
		return ""
	}
	reg :=regexp.MustCompile(`<img\s*src="(.*?)"\s*title="点击看更多海报"\s*alt=".*"\s*rel="v:image"\s*/>`)
	result :=reg.FindAllStringSubmatch(movieHtml,-1)
	if len(result) == 0{
		return ""
	}
	return string(result[0][1])
}

func GetMovieDirector(movieHtml string) string{
	if movieHtml == ""{
		return ""
	}
	reg:=regexp.MustCompile(`<a.*?rel="v:directedBy">(.*?)</a>`)
	result :=reg.FindAllStringSubmatch(movieHtml,-1)
	if len(result)== 0 {
		return ""
	}
	return string(result[0][1])
}

func GetMovieBianju(movieHtml string) string{
	if movieHtml ==""{
		return ""
	}
	reg :=regexp.MustCompile(`<a\s*href="/celebrity/\d+/">(.*?)</a>`)
	result :=reg.FindAllStringSubmatch(movieHtml,-1)
	if len(result)==0{
		return ""
	}
	bianju :=""
	for _,v :=range result{
		bianju += v[1]+"/"
	}
	return strings.Trim(bianju,"/")
}

func GetMovieMainCharacters(movieHtml string) string{
	reg :=regexp.MustCompile(`<a.*?rel="v:starring">(.*?)</a>`)
	result :=reg.FindAllStringSubmatch(movieHtml,-1)
	if len(result)==0{
		return ""
	}
	mainCharacters:=""
	for _,v:=range result{
		mainCharacters +=v[1]+"/"
	}
	return  strings.Trim(mainCharacters,"/")
}

func GetMovieGrade(movieHtml string)string{
	reg := regexp.MustCompile(`<strong.*?property="v:average">(.*?)</strong>`)
	result :=reg.FindAllStringSubmatch(movieHtml,-1)
	if len(result) == 0 {
		return ""
	}
	return string(result[0][1])
}

func Getage(movieHtml string) string{
		if movieHtml==""{
			return ""
		}
	    reg := regexp.MustCompile(`<span\s*class="pl">制片国家/地区:</span>\s*(.*?)\s*<br/>`)
		result :=reg.FindAllStringSubmatch(movieHtml,-1)
		if len(result)==0 {
			return  ""
		}
		return string(result[0][1])
}

func GetMovieLanguage(movieHtml string) string{
	if movieHtml == ""{
		return ""
	}
	reg:=regexp.MustCompile(`<span\s*class="pl">语言:</span>\s*(.*?)<br/>`)
	result :=reg.FindAllStringSubmatch(movieHtml,-1)
	if len(result)==0{
		return ""
	}
	return string(result[0][1])
}

func GetMovieGenre(movieHtml string) string{
	reg := regexp.MustCompile(`<span.*?property="v:genre">(.*?)</span>`)
	result :=reg.FindAllStringSubmatch(movieHtml,-1)
	 if len(result) == 0{
		return ""
	 }
	movieGenre :=""
	for _,v :=range result{
		movieGenre += v[1]+"/"
	}
	return strings.Trim(movieGenre,"/")
}

func GetMovieOnTime(movieHtml string) string{
	reg :=regexp.MustCompile(`<span.*?property="v:initialReleaseDate".*?>(.*?)</span>`)
	result :=reg.FindAllStringSubmatch(movieHtml,-1)
	if len(result) == 0{
		return ""
	}
	return string(result[0][1])
}

func GetMovieRunningTime(movieHtml string) string{
	reg :=regexp.MustCompile(`<span.*?property="v:runtime".*?>(.*?)</span>`)
	result :=reg.FindAllStringSubmatch(movieHtml,-1)
	if len(result) ==0{
		return ""
	}
	return string(result[0][1])
}

//获取当前豆瓣下面页下对的所有相关电影url
func GetMovieUrls(movieHtml string) []string {
	reg := regexp.MustCompile(`<a.*?href="(https://movie.douban.com/.*?)"`)
	//
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	var movieSets []string
	for _,v :=range result{
		movieSets =append(movieSets,v[1])
	}
	return movieSets
}
