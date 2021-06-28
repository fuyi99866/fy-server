package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"go_server/global"
	"net/smtp"
)

func EmailTest(subject string, body string) error {
	to := []string{global.GVA_CONFIG.Email.From}
	return send(to, subject, body)
}


func send(to []string, subject string, body string) error {
	/*from := global.GVA_CONFIG.Email.From
	nickname := global.GVA_CONFIG.Email.Nickname
	secret := global.GVA_CONFIG.Email.Secret
	host := global.GVA_CONFIG.Email.Host
	port := global.GVA_CONFIG.Email.Port
	isSSL := global.GVA_CONFIG.Email.IsSSL*/

	from := "xxx@163.com"
	nickname := "test"
	secret := "xxx"
	host := "smtp.163.com"
	port := 465
	isSSL := true

	auth := smtp.PlainAuth("", from, secret, host)
	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		e.From = from
	}
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)
	var err error
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if isSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	return err
}
