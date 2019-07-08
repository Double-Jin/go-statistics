package createLog

import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)
const BaseUrl  = "http://localhost:8000/"
var UaList = []string{
	"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; AcooBrowser; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Acoo Browser; SLCC1; .NET CLR 2.0.50727; Media Center PC 5.0; .NET CLR 3.0.04506)",
	"Mozilla/4.0 (compatible; MSIE 7.0; AOL 9.5; AOLBuild 4337.35; Windows NT 5.1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
	"Mozilla/5.0 (Windows; U; MSIE 9.0; Windows NT 9.0; en-US)",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 2.0.50727; Media Center PC 6.0)",
	"Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 1.0.3705; .NET CLR 1.1.4322)",
	"Mozilla/4.0 (compatible; MSIE 7.0b; Windows NT 5.2; .NET CLR 1.1.4322; .NET CLR 2.0.50727; InfoPath.2; .NET CLR 3.0.04506.30)",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN) AppleWebKit/523.15 (KHTML, like Gecko, Safari/419.3) Arora/0.3 (Change: 287 c9dfb30)",
	"Mozilla/5.0 (X11; U; Linux; en-US) AppleWebKit/527+ (KHTML, like Gecko, Safari/419.3) Arora/0.6",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.2pre) Gecko/20070215 K-Ninja/2.1.1",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.9) Gecko/20080705 Firefox/3.0 Kapiko/3.0",
	"Mozilla/5.0 (X11; Linux i686; U;) Gecko/20070322 Kazehakase/0.4.5",
	"Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.8) Gecko Fedora/1.9.0.8-1.fc10 Kazehakase/0.5.6",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/535.20 (KHTML, like Gecko) Chrome/19.0.1036.7 Safari/535.20",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; fr) Presto/2.9.168 Version/11.52",
}
type resource struct {
	url string
	target string
	start int
	end int
}

func RuleResource() []resource  {
	var res []resource
	//首页
	r1 := resource{
		url:BaseUrl,
		target:"",
		start:0,
		end:0,

	}

	//列表页
	r2 := resource{
		url:BaseUrl+"list/{$id}.html",
		target:"{$id}",
		start:1,
		end:21,
	}

	//详情页
	r3 := resource{
		url:BaseUrl+"movie/{$id}.html",
		target:"{$id}",
		start:1,
		end:12924,
	}

	res = append(append(append(res,r1),r2),r3)
	return res

}

func BuildUrl(res []resource) []string {
	var list []string
	for _,resItem := range res {
		if len(resItem.target)==0 {
			list = append(list, resItem.url)
		}else {
			for i:=resItem.start;i<resItem.end;i++  {
				//返回将s中前n个不重叠old子串都替换为new的新字符串，如果n<0会替换所有old子串
				urlStr := strings.Replace(resItem.url,resItem.target,strconv.Itoa(i),-1)
				list = append(list,urlStr)

			}
		}
	}
	return list
}

func makeLog(current,refer,ua string) string {
	u:= url.Values{}
	u.Set("time",time.Now().String())
	u.Set("url",current)
	u.Set("refer",refer)
	u.Set("ua",ua)
	paramStr := u.Encode()
	logTemplate := "127.0.0.1 - - [19/Oct/2018:22:58:25 -0400] \"GET /dig?{$param_str} HTTP/1.1\" 200 43 \"-\" \"{$ua}\" \"-\""
	log:= strings.Replace(logTemplate,"{$param_str}",paramStr,-1)
	log = strings.Replace(log,"{$ua}",ua,-1)
	return log
}

func randInt( min, max int) int  {
	r := rand.New( rand.NewSource( time.Now().UnixNano()))

	if min > max {
		return max
	}

	return r.Intn(max - min) + min
}

func WriteLog(logStr,filePath string) bool{
	fd,err := os.OpenFile(filePath,os.O_RDWR|os.O_APPEND,0644)
	if err !=nil {
		return false
	}
	fd.Write([]byte(logStr))
	defer fd.Close()
	return true
}


func main()  {
	total := flag.Int("total", 100, "how many rows")
	filePath := flag.String("filePath","E:/phpStudy/PHPTutorial/nginx/logs/dig.log","log file path")
	flag.Parse()

	//需要构建出真实的网站URL集合
	res := RuleResource()
	list := BuildUrl(res)
	//按照要求， 生成$total 行日志内容，源自的这个集合

	//根据url生成日志内容
	//logStrChan := make(chan string,*total)
	logStr := ""
	for i := 0;i<=*total ;i++  {
		currentUrl := list[randInt(0,len(list)-1)]
		refer := list[randInt(0,len(list)-1)]
		ua := UaList[randInt(0,len(UaList)-1)]
		logStr += makeLog(currentUrl,refer,ua)+"\n"
		//logStrChan <- makeLog(currentUrl,refer,ua)+"\n"
	}

	//处理完批量写
	WriteLog(logStr,*filePath)

	fmt.Println("done.\n")
}

