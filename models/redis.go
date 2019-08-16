package models

import(
	"github.com/astaxie/goredis"
)

const(
	URL_QUEUE="url_queue"  //队列
	URL_VISIT_SET="url_visit_set"//记录访问过的url
)
//创建客户端redis
var(
	client goredis.Client
)

func ConnectRedis(addr string){
	client.Addr=addr
}
//提取url放入队列 Lpush()
func PutinQueue(url string){
	client.Lpush(URL_QUEUE,[]byte(url))
}
//从队列取 Rpop()
func PopformQueue() string{
	res,err :=client.Rpop(URL_QUEUE)
	if err != nil{
		panic(err)
	}
	return string(res)
}
//获取队列长度  Llen()
func GetQueueLength() int{
	length,err :=client.Llen(URL_QUEUE)
	if err!=nil{
		return 0
	}
	return length
}
//访问过的放在一个集合  sadd()插入
func AddToSet(url string){
	client.Sadd(URL_VISIT_SET,[]byte(url))
}
//判断某个URL是否访问过  sismember()
func IsVisit(url string) bool{
	bIntVisit,err :=client.Sismember(URL_VISIT_SET,[]byte(url))
	if err !=nil{
		return false
	}
	return bIntVisit
}