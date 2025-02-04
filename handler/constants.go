package handler

// 人机校验
const ANSWER = "Ice"

// 性别
var Gender = map[int8]struct{}{
	1: {}, // 男
	2: {}, // 女
	3: {}, // 未知
}

// 申请状态
const (
	PENDING  int8 = iota + 1 // 待处理
	ACCEPTED                 // 已同意
	REJECTED                 // 已拒绝
	IGNORED                  // 已忽略
)

// 邀请码
const (
	UNUSED  = -1 // 未使用
	USED    = 1  // 已使用
	EXPIRED = -2 // 已过期
)
