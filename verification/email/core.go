package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	log "github.com/sirupsen/logrus"
)

type EmailCore struct {
	// 邮件服务器地址
	SMTP_MAIL_HOST string
	// 端口
	SMTP_MAIL_PORT int64
	// 发送邮件用户账号
	SMTP_MAIL_USER string
	// 授权密码
	SMTP_MAIL_PWD string
	// 发送邮件昵称
	SMTP_MAIL_NICKNAME string
}

// set QQ-email as default email server
func NewEmailCore() *EmailCore {
	emailCore := EmailCore{}
	emailCore.SMTP_MAIL_HOST = ""
	emailCore.SMTP_MAIL_PORT = 544
	emailCore.SMTP_MAIL_USER = ""
	emailCore.SMTP_MAIL_PWD = ""
	emailCore.SMTP_MAIL_NICKNAME = ""
	return &emailCore
}

// send message of html template
func (emailCore *EmailCore) SendHtmlMessage(To, Subject, body string) error {
	header := make(map[string]string)
	header["From"] = emailCore.SMTP_MAIL_NICKNAME + "<" + emailCore.SMTP_MAIL_USER + ">"
	header["To"] = To
	header["Subject"] = Subject
	//html格式邮件
	header["Content-Type"] = "text/html; charset=UTF-8"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth(
		"",
		emailCore.SMTP_MAIL_USER,
		emailCore.SMTP_MAIL_PWD,
		emailCore.SMTP_MAIL_HOST,
	)
	err := SendMailWithTLS(
		fmt.Sprintf("%s:%d", emailCore.SMTP_MAIL_HOST, emailCore.SMTP_MAIL_PORT),
		auth,
		emailCore.SMTP_MAIL_USER,
		[]string{To},
		[]byte(message),
	)
	if err != nil {
		log.Errorln("Send email error:", err)
	} else {
		log.Errorln("Send mail success!")
	}
	return err
}

// 发送验证码
func (ec *EmailCore) SendCodeMessage(To, code string) {
	log.Debugln(code)
	ec.SendHtmlMessage(To, "OOPS Captcha", fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>OOPS验证服务</title>
	</head>
	<body>
		<h1>OOPS验证服务</h1>
		<p>你的验证码是 %s</p>
		<p>请妥善保管你的验证码，不要将验证码告诉他人</p>
	</body>
	</html>`, code))
}

// 注册邮箱验证的连接
func (ec *EmailCore) SendRegisterLinkMessage(To, Subject, body string) {
	ec.SendHtmlMessage(To, Subject, ``)
}

// Dial return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("tls.Dial Error:", err)
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// SendMailWithTLS send email with tls
func SendMailWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Errorln("Create smtp client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Errorln("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
