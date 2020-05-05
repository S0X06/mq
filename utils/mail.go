package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"gopkg.in/gomail.v2"
)

var (
	//服务器
	Host = "smtp.163.com"
	Port = 465
	//发送
	FromUser = "735594423@qq.com"
	PassWord = "123"
)

func SendEmail(mailTo string, subject string, body string) error {

	m := gomail.NewMessage()

	//设置发件人
	m.SetHeader("From", FromUser)

	//设置发送给多个用户
	mailArrTo := strings.Split(mailTo, ",")
	m.SetHeader("To", mailArrTo...)

	//设置邮件主题
	m.SetHeader("Subject", subject)

	//设置邮件正文
	m.SetBody("text/html", body)

	d := gomail.NewDialer(Host, Port, FromUser, PassWord)

	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func SendSmtEmail() {
	// 连接到远程 SMTP 服务器。
	client, err := smtp.Dial("mail.example.com:25")
	if err != nil {
		log.Fatal(err)
	}
	// 设置寄件人和收件人
	client.Mail("sender@example.org")
	client.Rcpt("recipient@example.net")
	// 发送邮件主体。
	wc, err := client.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString("This is the email body.")
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
}

func SendSmtEmailPlainAuth() {
	// 设置认证信息。
	auth := smtp.PlainAuth(
		"",
		"user@example.com",
		"password",
		"mail.example.com",
	)
	// 连接到服务器, 认证, 设置发件人、收件人、发送的内容,
	// 然后发送邮件。
	err := smtp.SendMail(
		"mail.example.com:25",
		auth,
		"sender@example.org",
		[]string{"recipient@example.net"},
		[]byte("This is the email body."),
	)
	if err != nil {
		log.Fatal(err)
	}
}
