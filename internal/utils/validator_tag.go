package utils

// 支持的性别
const (
	Male       = 1  // 男性
	Female     = 2  // 女性
	Non_Binary = -1 // 非二元性别
)

// 如果性别是支持的, 则返回true
func IsSupportedGender(gender int16) bool {
	switch gender {
	case Male, Female, Non_Binary:
		return true
	}
	return false
}
