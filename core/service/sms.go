package service

import (
	"net/mail"

	"github.com/elvis88/baas/common/sms"
)

const (
	mailURL      = "smtp.qq.com:465"
	mailFrom     = "342529499@qq.com"
	mailFromName = "BaasSupport"
	mailFromPWD  = "ctcqcccjmzmdbhjf"

	mailBody = `<html>
    <head>
        <meta charset="UTF-8" />
        <title>fractal monitor</title>
    </head>
    <body>
       您好，您在Baas的%s验证码为：%s, 该验证码在5分钟内有效。 
    </body>
</html>`
)

var emailClient *sms.SMTPClient

func init() {
	email, err := sms.NewSMTPClient(mailURL, mail.Address{Name: mailFromName, Address: mailFrom}, mailFromPWD)
	if err != nil {
		panic("Failed to send email: " + err.Error())
	}
	emailClient = email
}
