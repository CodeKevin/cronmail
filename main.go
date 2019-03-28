package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"time"

	top "cronmail/controller"
	"cronmail/model"

	"github.com/robfig/cron"
	gomail "gopkg.in/gomail.v2"
)

func main() {
	//cronStart()
	server()
}
func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body := top.GetData("beijing")
		io.WriteString(w, body)
	}) //设置访问的路由
	err := http.ListenAndServe(":8080", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func sendMail(users []model.User) {
	for _, user := range users {
		body := top.GetData(user.City) //这里可以优化！！！
		go send(body, user)
	}

}
func send(body string, user model.User) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "xxxx@xx.com", user.Name+"的助理") // 发件人
	m.SetHeader("To",                                          // 收件人
		m.FormatAddress(user.Address, user.Name),
	)
	date := time.Now().Format("2006-01-02")
	m.SetHeader("Subject", "Hi,"+user.Name+". 今天是: "+date+" !") // 主题
	m.SetBody("text/html", body)                                // 正文

	d := gomail.NewPlainDialer("smtp.xxx.com", 465, "xxxx@xx.com", "授权码") // 发送邮件服务器、端口、发件人账号、发件人密码
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
func cronStart() {
	c := cron.New()
	spec := "0 0 8 * * ?"
	var list []model.User
	kevin := model.User{Address: "xxxx@xx.com", Name: "Kevin", City: "shanghai"}
	list = append(list, kevin)
	c.AddFunc(spec, func() {
		sendMail(list)
	})
	c.Start()
	select {}
}
