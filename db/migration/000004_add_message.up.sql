-- 消息表
CREATE TABLE IF NOT EXISTS "messages" (
    "id" BIGSERIAL PRIMARY KEY,
    "sender_id" INT NOT NULL,
    "receiver_id" INT NOT NULL,
    "message_type" SMALLINT NOT NULL DEFAULT 1,
    "content" JSONB NOT NULL,
    "content_type" SMALLINT NOT NULL DEFAULT 1,
    "message_status" SMALLINT NOT NULL DEFAULT 1,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN "messages"."id" IS '消息ID';

COMMENT ON COLUMN "messages"."sender_id" IS '发送者ID';

COMMENT ON COLUMN "messages"."receiver_id" IS '接收者ID, 用户或群组ID';

COMMENT ON COLUMN "messages"."message_type" IS '消息类型, 1: 私聊, 2: 群聊, 3: 系统消息';

COMMENT ON COLUMN "messages"."content" IS '消息内容';

COMMENT ON COLUMN "messages"."content_type" IS '消息内容类型, 1: 文字, 2: 文件, 3: 图片, 4: 语音, 5: 视频';

COMMENT ON COLUMN "messages"."message_status" IS '消息状态, 1: 已发送, 2: 已读, 3: 删除, 4: 撤回';

COMMENT ON COLUMN "messages"."created_at" IS '发送时间';

COMMENT ON COLUMN "messages"."updated_at" IS '更新时间';