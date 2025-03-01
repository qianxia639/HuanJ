package handler

// 人机校验
const Answer = "Weather"

// 申请状态
const (
	Pending  int8 = iota + 1 // 待处理
	Accepted                 // 已同意
	Rejected                 // 已拒绝
	Ignored                  // 已忽略
)

// 角色
const (
<<<<<<< HEAD
	GroupLeader int16 = iota + 1 // 群主
	Admin                        // 管理员
	Member                       // 普通成员
=======
	GroupOwner int16 = iota + 1 // 群主
	Admin                       // 管理员
	Member                      // 普通成员
)

// 接收者类型
const (
	User  int16 = iota + 1 // 用户
	Group                  // 群组
)

// 发送类型
const (
	PrivateChat int16 = iota + 1 // 私聊
	GroupChat                    // 群聊
	Heartbeat                    // 心跳
)

// 消息类型
const (
	Text  int16 = iota + 1 // 文本
	File                   // 文件
	Image                  // 图片
	Audio                  // 音频
	Video                  // 视频
>>>>>>> f4a28b92984017d1e60e2fdc387226202641208d
)
