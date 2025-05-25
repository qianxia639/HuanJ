package mail

import (
	"testing"
)

func TestMail(t *testing.T) {

	// 跳过当前测试
	t.Skip()

	to := "2274000859@qq.com"
	subject := "测试邮件"
	body := "测试邮件发送"

	t.Log("send email")
	err := SendToMail(to, subject, body)
	if err != nil {
		t.Error("send mail error")
		t.Log(err)
		return
	} else {
		t.Log("send mail success!")
	}
}
