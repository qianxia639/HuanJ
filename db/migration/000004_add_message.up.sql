-- 消息表
CREATE TABLE IF NOT EXISTS "messages" (
    "id" SERIAL PRIMARY KEY,
    "session_id" VARCHAR(64) NOT NULL,
    "sender_id" INT NOT NULL,
    "receiver_id" INT NOT NULL,
    "send_type" SMALLINT NOT NULL DEFAULT 1,
    "receiver_type" SMALLINT NOT NULL DEFAULT 1,
    "message_type" SMALLINT NOT NULL DEFAULT 1,
    "content" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN "messages"."id" IS '消息ID';

COMMENT ON COLUMN "messages"."session_id" IS '会话ID, 用来做消息标识,格式:user/group:sender_id:receiver_id';

COMMENT ON COLUMN "messages"."sender_id" IS '发送者ID';

COMMENT ON COLUMN "messages"."receiver_id" IS '接收者ID, 用户或群组ID';

COMMENT ON COLUMN "messages"."receiver_type" IS '接收者类型, 1: 用户, 2: 群组';

COMMENT ON COLUMN "messages"."send_type" IS '发送类型, 1: 私聊, 2: 群聊, 3: 心跳';

COMMENT ON COLUMN "messages"."message_type" IS '消息类型, 1: 文字, 2: 文件, 3: 图片, 4: 语音, 5: 视频';

COMMENT ON COLUMN "messages"."content" IS '消息内容';

COMMENT ON COLUMN "messages"."created_at" IS '发送时间';