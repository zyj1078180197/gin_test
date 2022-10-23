package test

import (
	"crypto/tls"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)

func TestEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Get <rain971124@163.com>"
	e.To = []string{"1078180197@qq.com"}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码是<b>123456</b><")
	//返回EOF时，关闭SSL重试
	er := e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "rain971124@163.com", "CVOWBHXPZLCUMOZV", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if er != nil {
		t.Error(er)
		return
	}

}
