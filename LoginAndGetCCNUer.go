package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// 定义学生结构体
type student struct {
	Id   string `json:"id"`
	Pid  string `json:"Pid"`
	Name string `json:"name"`
	//Label       string `json:"label"`
	//SzLogonName string `json:"szLogonname"`
	//SzHandPhone string `json:"szHandPhone"`
	//SzTel       string `json:"szTel"`
	//SzEmail     string `json:"szEmail"`
}

// 获取登录界面html
func fech(url string) string {
	client := &http.Client{}
	//发送Get请求
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.62")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

// 解析html
func parse(html string) (string, string, string) {
	//替换空格
	html = strings.Replace(html, "\n", "", -1)
	//定义body正则
	ReBody := regexp.MustCompile(`<body id="cas">(.*?)</body>`)
	//找到body
	Body := ReBody.FindString(html)
	//定义jsessionid正则
	FindJsessionid := regexp.MustCompile(`jsessionid=(.*?)"`)
	//在body中找到jsessionid
	Jsessionid := FindJsessionid.FindString(Body)
	//除去jsessionid末尾的"
	Jsessionid = Jsessionid[:(len(Jsessionid) - 1)]
	//拼接登录所需的url
	url1 := "https://account.ccnu.edu.cn/cas/login;" + Jsessionid
	//进一步解析jsessionid
	ReJsessionid := regexp.MustCompile(`=(.*?)\?`)
	Jsessionid = ReJsessionid.FindString(Jsessionid)
	//去除jsessionid首部的=与尾部的？
	Jsessionid = Jsessionid[1:(len(Jsessionid) - 1)]
	//定义lt正则
	FindLt := regexp.MustCompile(`name="lt" value="(.*?)"`)
	//在body中找到lt值
	Lt := FindLt.FindString(Body)
	//进一步解析lt
	ReLt := regexp.MustCompile(`value="(.*?)"`)
	Lt = ReLt.FindString(Lt)
	Lt = Lt[7:(len(Lt) - 1)]
	//返回url及jsessionid及lt
	return url1, Jsessionid, Lt
}

func LoginAndSearch(url1, Jsessionid, Lt string) {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	client := http.Client{
		Jar: jar, //初始化cookie容器
	}
	url1 = fmt.Sprintf(url1)
	var LogUser, LogPwd string
	fmt.Println("输入账号及密码，分两行输入")
	_, err := fmt.Scanln(&LogUser)
	if err != nil {
		return
	}
	_, err = fmt.Scanln(&LogPwd)
	if err != nil {
		return
	}
	//创建负载表单
	FormData := fmt.Sprintf("username=%v&password=%v&lt=%v&execution=e1s1&_eventId=submit&submit=%E7%99%BB%E5%BD%95", LogUser, LogPwd, Lt) //登录的body
	ReqBody := strings.NewReader(FormData)
	//新建POST请求
	LogReq, err := http.NewRequest(http.MethodPost, url1, ReqBody)
	if err != nil {
		return
	}
	//添加请求头
	LogReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	LogReq.Header.Set("Accept-Encoding", "gzip, deflate, br")
	LogReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	LogReq.Header.Set("Cache-Control", "max-age=0")
	LogReq.Header.Set("Connection", "keep-alive")
	LogReq.Header.Set("Content-Length", "166")
	LogReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//更改Jsessionid格式
	Jsessionid = "JSESSIONID=" + Jsessionid
	LogReq.Header.Set("Cookie", Jsessionid)
	LogReq.Header.Set("Host", "account.ccnu.edu.cn")
	LogReq.Header.Set("Origin", "https://account.ccnu.edu.cn")
	LogReq.Header.Set("Referer", "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page=")
	LogReq.Header.Set("sec-ch-ua", "\"Microsoft Edge\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"")
	LogReq.Header.Set("sec-ch-ua-mobile", "?0")
	LogReq.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	LogReq.Header.Set("Sec-Fetch-Dest", "document")
	LogReq.Header.Set("Sec-Fetch-Mode", "navigate")
	LogReq.Header.Set("Sec-Fetch-Site", "same-origin")
	LogReq.Header.Set("Sec-Fetch-User", "?1")
	LogReq.Header.Set("Upgrade-Insecure-LogRequests", "1")
	LogReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.62")
	_, err = client.Do(LogReq)
	if err != nil {
		return
	}
	fmt.Println("输入查询范围")
	var start, end int
	_, err = fmt.Scan(&start, &end)
	if err != nil {
		return
	}
	for i := start; i <= end; i++ {
		URL := fmt.Sprintf("http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/data/searchAccount.aspx?type=logonname&ReservaApply=ReservaApply&term=%v&_=1669882638320", i)
		req, err := http.NewRequest(http.MethodGet, URL, nil)
		if err != nil {
			log.Println("err")
		}
		req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Host", "kjyy.ccnu.edu.cn")
		req.Header.Set("Referer", "http://kjyy.ccnu.edu.cn/clientweb/xcus/ic2/Default.aspx")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.56")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		resp, err := client.Do(req)
		if err != nil {
			log.Println("err")
		}
		var body []byte
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Println("err")
		}
		body = append(body[:0], body[1:]...)
		length := len(body)
		body = append(body[:(length - 1)])
		var ThisStudent student
		err = json.Unmarshal(body, &ThisStudent)
		if err != nil {
			return
		}
		fmt.Println(ThisStudent)
		title := strconv.Itoa(i)
		//保存每一份数据到单独的md文件
		save(title, body)
	}
}

// 保存数据
func save(title string, content []byte) {
	//保存到当前目录下worm目录中的students目录
	err := os.WriteFile("./worm/students/"+title+".md", content, 0644)
	if err != nil {
		return
	}
}

func main() {
	url := "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
	s := fech(url)
	url1, Jsessionid, Lt := parse(s)
	LoginAndSearch(url1, Jsessionid, Lt)
}
