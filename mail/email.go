package mail

import (
	"HuanJ/config"
	"HuanJ/logs"
	"HuanJ/utils"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// 发送邮箱验证码
func SendEmailCode(rdb *redis.Client, email string, emailCodeType int8) error {

	// 生成邮箱验证码
	emailCode := utils.RandomNum(100000, 999999)

	// 发送邮箱验证码
	notice := config.EmailMap[emailCodeType]
	err := SendToMail(email, notice.Title, fmt.Sprintf(notice.Body, emailCode))
	if err != nil {
		logs.Error("发送邮箱验证码失败: ", err.Error())
		return err
	}
	// 缓存邮箱验证码
	key := fmt.Sprintf("email:%d:%s", emailCodeType, email)
	err = rdb.Set(context.Background(), key, emailCode, 5*time.Minute).Err()
	if err != nil {
		logs.Error("缓存邮箱验证码失败: ", err.Error())
		return err
	}

	return nil
}

// Deprecated
func GetEmailCode(rdb *redis.Client, email string) (string, error) {
	key := "email:" + email
	emailCode, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		logs.Error("获取邮箱验证码失败: ", err.Error())
		return "", err
	}

	return emailCode, nil
}

// 校验邮箱验证码
func VerifyEmailCode(rdb *redis.Client, email, emailCode string, emailCodeType int8) (bool, error) {
	key := fmt.Sprintf("email:%d:%s", emailCodeType, email)
	code, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	return code == emailCode, nil
}
