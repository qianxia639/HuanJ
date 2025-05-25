package mail

import (
	"bytes"
	"crypto/tls"
	"mime"
	nmail "net/mail"
	"net/smtp"
	"time"
)

var username = "example@qq.com"
var password = "password"
var smtpHost = "smtp.qq.com"

func SendToMail(to, subject, body string) error {

	msg := buildEmail(to, subject, body)

	// 使用SSL/TLS加密连接
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	addr := smtpHost + ":587"

	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	err = client.StartTLS(tlsconfig)
	if err != nil {
		return err
	}

	// SMTP认证
	auth := smtp.PlainAuth("", username, password, smtpHost)
	err = client.Auth(auth)
	if err != nil {
		return err
	}

	err = client.Mail(username)
	if err != nil {
		return err
	}

	err = client.Rcpt(to)
	if err != nil {
		return err
	}

	w, err := client.Data()
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

	_ = client.Quit()

	return nil
}

func buildEmail(to, subject, body string) []byte {

	// 编码主题(处理中文等特殊字符)
	encoodedSubject := mime.QEncoding.Encode("UTF-8", subject)

	buf := bytes.NewBuffer(nil)

	from := nmail.Address{Name: "幻晶", Address: username}

	headers := []string{
		"From: " + from.String(),
		"To: " + to,
		"Subject: " + encoodedSubject,
		"Date: " + time.Now().Format(time.RFC1123Z),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Transfer-Encoding: quoted-printable",
	}

	// 安全拼接

	for _, h := range headers {
		buf.WriteString(h + "\r\n")
	}
	buf.WriteString("\r\n") // 头部与正文分隔
	buf.WriteString(body)

	return buf.Bytes()
}
