package main

import (
	"fmt"
	"net/smtp"

	"github.com/astaxie/beego/config"
)

// Send an email
func Send(body, title, to string) error {

	conf, err := config.NewConfig("ini", "smtp.conf")
	if err != nil {
		return err
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	from := conf.String("smtp::email")

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n%s%s",
		from, to, title, mime, body)

	port, err := conf.Int("smtp::port")
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%d",
		conf.String("smtp::server"), port)

	auth := smtp.PlainAuth("",
		from, conf.String("smtp::password"),
		conf.String("smtp::server"))

	err = smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
