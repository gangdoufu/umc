package sender

import (
	"context"
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"gopkg.in/gomail.v2"
)

var emailSender *Email

type Email struct {
	host     string
	port     int
	username string
	password string
	subject  string
}

func (e *Email) SendMessage(ctx context.Context, vo *InfoVo) error {
	m := e.getEmailMessage(vo)
	d := e.newDialer()
	return d.DialAndSend(m)
}

func (e Email) newDialer() *gomail.Dialer {
	return gomail.NewDialer(e.host, e.port, e.username, e.password)
}
func (e Email) getEmailMessage(vo *InfoVo) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", global.Config.Email.Username)
	m.SetHeader("To", vo.Receiver)
	m.SetHeader("Subject", subjectMap[vo.UseTo])
	m.SetBody("text/html", vo.Info)
	return m
}
