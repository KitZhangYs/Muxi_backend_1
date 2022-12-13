package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Id       int         `json:"id"`
	Username string      `json:"username"`
	Password string      `json:"password"`
	Name     string      `json:"name"`
	Age      string      `json:"age"`
	Sex      string      `json:"sex"`
	SelfWord string      `json:"self_word"`
	Ucookie  http.Cookie `json:"ucookie"`
	State    string      `json:"state"`
}

var TimelyCookies = make(map[string]http.Cookie)

var (
	userName  string = "root"
	password  string = ""
	ipAddrees string = "localhost"
	port      int    = 3306
	dbName    string = "userinfo"
	charset   string = "utf8"
)

var Db *sqlx.DB

// 连接数据库
func connectMysql() *sqlx.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, ipAddrees, port, dbName, charset)
	Db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("mysql connect failed, detail is [%v]", err.Error())
	} else {
		fmt.Println("connected")
	}
	return Db
}

// 增加用户信息
func addRecord(Db *sqlx.DB, user *User) {
	for i := 0; i < 1; i++ {
		result, err := Db.Exec("insert into users  values(?,?,?,?,?,?,?,?)", nil, user.Username, user.Password, user.Name, user.Age, user.Sex, user.SelfWord, user.State)
		if err != nil {
			fmt.Printf("data insert faied, error:[%v]", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		fmt.Printf("insert success, last id:[%d]\n", id)
	}
}

// 在数据库中查询用户信息
func queryData(Db *sqlx.DB, table string, user *User) (bool, *User) {
	qry := fmt.Sprintf("select * from %s where username='%s'", table, user.Username)
	rows, err := Db.Query(qry)
	if err != nil {
		fmt.Printf("query faied, error:[%v]", err.Error())
		return false, nil
	}
	var s User //定义s储存返回数值
	//若数据库中没有查询到该用户，则不会进入循环
	for rows.Next() {
		err = rows.Scan(&s.Id, &s.Username, &s.Password, &s.Name, &s.Age, &s.Sex, &s.SelfWord, &s.State)
		if err != nil {
			fmt.Println("err", err.Error())
			return false, nil
		}
		fmt.Println("user was found")
		goto A
	}
	err = rows.Close()
	if err != nil {
		return false, nil
	}
	return false, nil
A:
	return true, &s
}

// FormConversion 解析请求头
func FormConversion(req *http.Request) User {
	contentLength := req.ContentLength
	request := make([]byte, contentLength)
	//_, err := req.Body.Read(request)
	request, err := io.ReadAll(req.Body)
	if err != nil {
		return User{}
	}
	var a User
	//fmt.Println(string(request))
	err = json.Unmarshal(request, &a)
	if err != nil {
		return User{}
	}
	return a
}

// 首页
func homepage(res http.ResponseWriter, req *http.Request) {
	var theUser User
	theUser = FormConversion(req)
	ok, FindUser := queryData(Db, "users", &theUser)
	cookie, Err := req.Cookie(theUser.Username)
	if Err == nil && cookie.Value == TimelyCookies[FindUser.Username].Value && FindUser.State == "online" && ok {
		_, err := fmt.Fprintf(res, "用户：%s,登陆成功，火速来当快递员，跟我一起激情满满滴寄件吧！", FindUser.Name)
		if err != nil {
			return
		}
	} else {
		_, err := res.Write([]byte("奶奶滴，登！为什么不登！\n不登？不登是吧！不登，也白想活着！"))
		if err != nil {
			return
		}
	}
}

// 登录
func login(res http.ResponseWriter, req *http.Request) {
	var LogUser User
	LogUser = FormConversion(req)
	ok, ThisUser := queryData(Db, "users", &LogUser)
	if ok {
		if ThisUser.Password == LogUser.Password {
			value := strconv.Itoa(int(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000)))
			cookie := &http.Cookie{
				Name:  LogUser.Username,
				Value: value,
			}
			http.SetCookie(res, cookie)
			LogUser.Ucookie = *cookie
			TimelyCookies[LogUser.Username] = LogUser.Ucookie
			LogUser.State = "online"
			qry := fmt.Sprintf("update users set state='online' where username='%s'", LogUser.Username)
			_, err := Db.Exec(qry)
			if err != nil {
				return
			}
			_, err = res.Write([]byte("密码正确，您已登录,/homepage有些好康的，比新游戏还刺激~"))
			if err != nil {
				return
			}
		} else {
			_, err := res.Write([]byte("密码错误！这是你吗？这特喵就不是你！"))
			if err != nil {
				return
			}
		}
	} else {
		_, err := res.Write([]byte("诶？用户捏？这人我不熟啊，快去/signup注册吧！"))
		if err != nil {
			return
		}
	}
}

// 注册
func signup(res http.ResponseWriter, req *http.Request) {
	var NewUser User
	NewUser = FormConversion(req)
	end := &NewUser
	end.State = "outline"
	ok, _ := queryData(Db, "users", end)
	if NewUser.Username != "" && NewUser.Password != "" && !ok {
		addRecord(Db, end)
		_, err := res.Write([]byte("注册成功，速来登录"))
		if err != nil {
			return
		}
	} else if ok {
		_, err := res.Write([]byte("注册失败，该用户名已被注册"))
		if err != nil {
			return
		}
	} else {
		_, err := res.Write([]byte("注册失败，请确定您输入了用户名与密码"))
		if err != nil {
			return
		}
	}
}

// ViewUserInformation 查询用户信息
func ViewUserInformation(res http.ResponseWriter, req *http.Request) {
	var theUser User
	theUser = FormConversion(req)
	ok, FindUser := queryData(Db, "users", &theUser)
	cookie, Err := req.Cookie(FindUser.Username)
	if Err == nil && cookie.Value == TimelyCookies[FindUser.Username].Value && FindUser.State == "online" && ok {
		_, err := fmt.Fprintf(res, "查询成功！\n用户名：%s\n昵称：%s\n年龄:%s\n性别：%s\n个性签名：%s\n", FindUser.Username, FindUser.Name, FindUser.Age, FindUser.Sex, FindUser.SelfWord)
		if err != nil {
			return
		}
	} else {
		_, err := res.Write([]byte("记得登录喵，记得登录谢谢喵"))
		if err != nil {
			return
		}
	}
}

// ChangeUserInformation 更改用户信息
func ChangeUserInformation(res http.ResponseWriter, req *http.Request) {
	var theUser User
	theUser = FormConversion(req)
	ok, FindUser := queryData(Db, "users", &theUser)
	cookie, Err := req.Cookie(FindUser.Username)
	if Err == nil && cookie.Value == TimelyCookies[FindUser.Username].Value && FindUser.State == "online" && ok {
		if theUser.Password != FindUser.Password {
			qry := fmt.Sprintf("update users set password='%s' where username='%s'", theUser.Password, theUser.Username)
			_, err := Db.Exec(qry)
			if err != nil {
				return
			}
			theUser.State = "outline"
			qry = fmt.Sprintf("update users set state='outline' where username='%s'", theUser.Username)
			_, err = Db.Exec(qry)
			if err != nil {
				return
			}
		}
		if theUser.Name != "" {
			qry := fmt.Sprintf("update users set name='%s' where username='%s'", theUser.Name, theUser.Username)
			_, err := Db.Exec(qry)
			if err != nil {
				return
			}
		}
		if theUser.Age != "" {
			qry := fmt.Sprintf("update users set age='%s' where username='%s'", theUser.Age, theUser.Username)
			_, err := Db.Exec(qry)
			if err != nil {
				return
			}
		}
		if theUser.SelfWord != "" {
			qry := fmt.Sprintf("update users set selfword='%s' where username='%s'", theUser.SelfWord, theUser.Username)
			_, err := Db.Exec(qry)
			if err != nil {
				return
			}
		}
		if theUser.Sex != "" {
			qry := fmt.Sprintf("update users set sex='%s' where username='%s'", theUser.Sex, theUser.Username)
			_, err := Db.Exec(qry)
			if err != nil {
				return
			}
		}
	} else {
		_, err := res.Write([]byte("请先登录喵，请先登录谢谢喵,温馨提示，用户名不可更改哦"))
		if err != nil {
			return
		}
	}
}

// 登出
func logout(res http.ResponseWriter, req *http.Request) {
	var theUser User
	theUser = FormConversion(req)
	ok, FindUser := queryData(Db, "users", &theUser)
	cookie, Err := req.Cookie(FindUser.Username)
	if Err == nil && cookie.Value == TimelyCookies[FindUser.Username].Value && FindUser.State == "online" && ok {
		theUser.State = "outline"
		qry := fmt.Sprintf("update users set state='outline' where username='%s'", theUser.Username)
		_, err := Db.Exec(qry)
		if err != nil {
			return
		}
	} else {
		_, err := res.Write([]byte("请先登录喵，请先登录谢谢喵"))
		if err != nil {
			return
		}
	}
}

func main() {
	mi := http.NewServeMux()
	ser := &http.Server{
		Addr:    ":57257",
		Handler: mi,
	}
	Db = connectMysql()
	mi.HandleFunc("/homepage", homepage)
	mi.HandleFunc("/login", login)
	mi.HandleFunc("/signup", signup)
	mi.HandleFunc("/detail", ViewUserInformation)
	mi.HandleFunc("/change", ChangeUserInformation)
	mi.HandleFunc("/logout", logout)
	err := ser.ListenAndServe()
	if err != nil {
		return
	}
}
