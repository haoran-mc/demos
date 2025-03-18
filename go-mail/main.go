package main

import (
	"github.com/go-gomail/gomail"
)

func main() {
	from := "haoran.mc@qq.com"
	password := "****************"
	server_host := "smtp.qq.com"
	server_port := 587
	username := "haoran.mc"

	// from := "187********@163.com"
	// server_host := "smtp.163.com"
	// server_port := 25
	// username := "187********"
	// password := "****************"

	to := "haoran.mc@outlook.com"
	subject := "测试邮件"
	text := "<b>This is the body of the mail</b>"

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", text)
	// msg.Attach("/home/User/cat.jpg")

	n := gomail.NewDialer(server_host, server_port, username, password)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}
}
