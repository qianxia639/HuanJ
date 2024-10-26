package handler

// 性别
var Gender = map[int8]struct{}{
	1: {}, // 男
	2: {}, // 女
	3: {}, // 未知
}

// 申请状态
const (
	PENDING  = iota + 1 // 待确认
	ACCEPTED            // 已确认
)
