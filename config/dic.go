package config

var EmailMap = map[int8]Notice{
	1: {Title: "绑定或修改邮箱", Body: "您正在尝试绑定或修改邮箱! 您的验证码为 %v ! 如非本人操作,请勿泄露此验证码! 验证码5分钟内有效。"},
	2: {},
}
